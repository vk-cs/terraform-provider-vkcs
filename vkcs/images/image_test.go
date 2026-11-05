package images

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResourceImagesImageFile(t *testing.T) {
	image := ResourceImagesImage().TestResourceData()
	imgCachePath := fmt.Sprintf("%s/.terraform/image_cache", os.Getenv("HOME"))
	image.Set("image_source_url", "http://download.cirros-cloud.net/0.3.5/cirros-0.3.5-x86_64-disk.img")
	image.Set("image_cache_path", imgCachePath)

	client := new(gophercloud.ServiceClient)
	client.ProviderClient = new(gophercloud.ProviderClient)

	filename, err := resourceImagesImageFile(client, image)

	assert.Equal(t, nil, err)
	assert.Equal(t, filepath.Join(imgCachePath, "e18f6225aa529b597b260081c6ecb1da.img"), filename)
}

func gzipBytes(t *testing.T, data []byte) []byte {
	t.Helper()
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	_, err := w.Write(data)
	require.NoError(t, err)
	require.NoError(t, w.Close())
	return buf.Bytes()
}

func TestResourceImagesImageFile_decompressOptionsChange(t *testing.T) {
	raw := []byte("raw image content")
	compressed := gzipBytes(t, raw)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/octet-stream")
		_, _ = w.Write(compressed)
	}))
	defer srv.Close()

	client := new(gophercloud.ServiceClient)
	client.ProviderClient = new(gophercloud.ProviderClient)
	cacheDir := t.TempDir()

	download := func(compressionFormat string) []byte {
		image := ResourceImagesImage().TestResourceData()
		require.NoError(t, image.Set("image_source_url", srv.URL+"/disk.raw.gz"))
		require.NoError(t, image.Set("image_cache_path", cacheDir))
		require.NoError(t, image.Set("compression_format", compressionFormat))

		filename, err := resourceImagesImageFile(client, image)
		require.NoError(t, err)
		content, err := os.ReadFile(filename)
		require.NoError(t, err)
		return content
	}

	assert.Equal(t, compressed, download(""), "without compression_format the file is cached as is")
	assert.Equal(t, raw, download(compressionFormatGZIP), "compression_format must not reuse the compressed file from cache")
	assert.Equal(t, compressed, download(""), "removing compression_format must not reuse the decompressed file from cache")
}

func TestImageCacheFilename(t *testing.T) {
	const (
		dir = "/cache"
		url = "http://example.com/disk.img.tar.gz"
	)
	legacy := filepath.Join(dir, fmt.Sprintf("%x.img", md5.Sum([]byte(url))))

	assert.Equal(t, legacy, imageCacheFilename(dir, url, "", ""), "without decompress options the name must not change")

	names := map[string]string{"": legacy}
	for _, opts := range [][2]string{
		{compressionFormatAuto, ""},
		{compressionFormatGZIP, ""},
		{compressionFormatBZIP2, ""},
		{compressionFormatXZ, ""},
		{"", archivingFormatTAR},
		{compressionFormatGZIP, archivingFormatTAR},
	} {
		name := imageCacheFilename(dir, url, opts[0], opts[1])
		assert.Equal(t, filepath.Clean(dir), filepath.Dir(name))
		for prev, prevName := range names {
			assert.NotEqual(t, prevName, name, "options %v collide with %q", opts, prev)
		}
		names[opts[0]+"/"+opts[1]] = name
	}
}

func TestResourceImagesImageFile_atomicWrite(t *testing.T) {
	raw := []byte("raw image content")
	cacheDir := t.TempDir()
	var srvURL string
	var existedDuringDownload bool

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/octet-stream")
		_, _ = w.Write(raw[:5])
		w.(http.Flusher).Flush()
		time.Sleep(100 * time.Millisecond)
		_, err := os.Stat(imageCacheFilename(cacheDir, srvURL+r.URL.Path, "", ""))
		existedDuringDownload = err == nil
		_, _ = w.Write(raw[5:])
	}))
	defer srv.Close()
	srvURL = srv.URL

	client := new(gophercloud.ServiceClient)
	client.ProviderClient = new(gophercloud.ProviderClient)
	image := ResourceImagesImage().TestResourceData()
	require.NoError(t, image.Set("image_source_url", srv.URL+"/disk.raw"))
	require.NoError(t, image.Set("image_cache_path", cacheDir))

	filename, err := resourceImagesImageFile(client, image)
	require.NoError(t, err)

	assert.False(t, existedDuringDownload, "cache file must appear only after the download is complete")
	content, err := os.ReadFile(filename)
	require.NoError(t, err)
	assert.Equal(t, raw, content)
	assert.NoFileExists(t, filename+".tmp")
}

func TestResourceImagesImageFile_noCacheFileOnError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/octet-stream")
		_, _ = w.Write([]byte("not a gzip stream"))
	}))
	defer srv.Close()

	client := new(gophercloud.ServiceClient)
	client.ProviderClient = new(gophercloud.ProviderClient)
	cacheDir := t.TempDir()
	image := ResourceImagesImage().TestResourceData()
	require.NoError(t, image.Set("image_source_url", srv.URL+"/disk.raw.gz"))
	require.NoError(t, image.Set("image_cache_path", cacheDir))
	require.NoError(t, image.Set("compression_format", compressionFormatGZIP))

	_, err := resourceImagesImageFile(client, image)
	require.Error(t, err)

	filename := imageCacheFilename(cacheDir, srv.URL+"/disk.raw.gz", compressionFormatGZIP, "")
	assert.NoFileExists(t, filename)
	assert.NoFileExists(t, filename+".tmp")
}
