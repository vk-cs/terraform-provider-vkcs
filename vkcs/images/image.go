package images

import (
	"archive/tar"
	"compress/bzip2"
	"compress/gzip"
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"

	"github.com/gofrs/flock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/ulikunitz/xz"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var imagesDefaultStoreEndpointMasks = []string{"*.devmail.ru$", "^ams.*"}

func resourceImagesImageMemberStatusFromString(v string) images.ImageMemberStatus {
	switch v {
	case string(images.ImageMemberStatusAccepted):
		return images.ImageMemberStatusAccepted
	case string(images.ImageMemberStatusPending):
		return images.ImageMemberStatusPending
	case string(images.ImageMemberStatusRejected):
		return images.ImageMemberStatusRejected
	case string(images.ImageMemberStatusAll):
		return images.ImageMemberStatusAll
	}

	return ""
}

func resourceImagesImageVisibilityFromString(v string) images.ImageVisibility {
	switch v {
	case string(images.ImageVisibilityPublic):
		return images.ImageVisibilityPublic
	case string(images.ImageVisibilityPrivate):
		return images.ImageVisibilityPrivate
	case string(images.ImageVisibilityShared):
		return images.ImageVisibilityShared
	case string(images.ImageVisibilityCommunity):
		return images.ImageVisibilityCommunity
	}

	return ""
}

func fileMD5Checksum(f *os.File) (string, error) {
	hash := md5.New()
	if _, err := io.Copy(hash, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func resourceImagesImageFileProps(filename string) (int64, string, error) {
	var filesize int64
	var filechecksum string

	file, err := os.Open(filename)
	if err != nil {
		return -1, "", fmt.Errorf("error opening file for Image: %s", err)
	}
	defer file.Close()

	fstat, err := file.Stat()
	if err != nil {
		return -1, "", fmt.Errorf("error reading image file %q: %s", file.Name(), err)
	}

	filesize = fstat.Size()
	filechecksum, err = fileMD5Checksum(file)
	if err != nil {
		return -1, "", fmt.Errorf("error computing image file %q checksum: %s", file.Name(), err)
	}

	return filesize, filechecksum, nil
}

func resourceImagesImageFile(client *gophercloud.ServiceClient, d *schema.ResourceData) (string, error) {
	if filename := d.Get("local_file_path").(string); filename != "" {
		return filename, nil
	}

	furl := d.Get("image_source_url").(string)
	if furl == "" {
		return "", fmt.Errorf("error in config. no file specified")
	}

	dir := d.Get("image_cache_path").(string)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", fmt.Errorf("unable to create dir %s: %s", dir, err)
	}

	filename := filepath.Join(dir, fmt.Sprintf("%x.img", md5.Sum([]byte(furl))))
	delFile := func() {
		if err := os.Remove(filename); err != nil {
			log.Printf("[DEBUG] Failed to cleanup the %q file: %s", filename, err)
		}
	}
	lockFilename := filename + ".lock"

	lock := flock.New(lockFilename)
	err := lock.Lock()
	if err != nil {
		return "", fmt.Errorf("unable to create file lock on file %s: %s", lockFilename, err)
	}
	defer func() {
		err := lock.Unlock()
		if err != nil {
			log.Printf("[WARN] There was an error unlocking filelock: %s", err)
		}
	}()

	info, err := os.Stat(filename)
	if err != nil && !os.IsNotExist(err) {
		return "", fmt.Errorf("error while trying to access file %q: %s", filename, err)
	}

	// check if the file size is zero
	// it could be a leftover from older provider versions
	if info != nil {
		if info.Size() != 0 {
			log.Printf("[DEBUG] File exists %s", filename)
			return filename, nil
		}
		// delete the zero size file
		delFile()
	}

	log.Printf("[DEBUG] File doesn't exists %s. will download from %s", filename, furl)
	file, err := os.Create(filename)
	if err != nil {
		return "", fmt.Errorf("error creating file %q: %s", filename, err)
	}
	defer file.Close()

	httpClient := &client.ProviderClient.HTTPClient
	request, err := http.NewRequest("GET", furl, nil)
	if err != nil {
		delFile()
		return "", fmt.Errorf("error creating a new request: %s", err)
	}

	username := d.Get("image_source_username").(string)
	password := d.Get("image_source_password").(string)
	if username != "" && password != "" {
		request.SetBasicAuth(username, password)
	}

	resp, err := httpClient.Do(request)
	if err != nil {
		delFile()
		return "", fmt.Errorf("error downloading image from %q: %s", furl, err)
	}

	// check for credential error among other errors
	if resp.StatusCode != http.StatusOK {
		delFile()
		return "", fmt.Errorf("error downloading image from %q, status code is %d", furl, resp.StatusCode)
	}

	defer resp.Body.Close()
	reader := resp.Body

	compressionFormat := d.Get("compression_format").(string)
	if compressionFormat == compressionFormatAuto {
		// If we're here "Content-Encoding" in not filled, we'll read
		// "Content-Type" to select format
		compressionFormat, err = getCompressionFormatFromContentType(resp.Header.Get("Content-Type"))
		if err != nil {
			delFile()
			return "", fmt.Errorf("error decompressing image %q: %s", furl, err)
		}
	}
	if compressionFormat != "" {
		decompressReader, err := selectDecompressReader(reader, compressionFormat)
		if err != nil {
			delFile()
			return "", fmt.Errorf("error decompressing image %q: %s", furl, err)
		}
		defer decompressReader.Close()
		reader = decompressReader
	}

	archivingFormat := d.Get("archiving_format").(string)
	if archivingFormat != "" {
		unzipReader, err := selectUnzipReader(reader, archivingFormat)
		if err != nil {
			delFile()
			return "", fmt.Errorf("error unzipping image %q: %s", furl, err)
		}
		defer unzipReader.Close()
		reader = unzipReader
	}

	if _, err = io.Copy(file, reader); err != nil {
		delFile()
		return "", fmt.Errorf("error downloading image %q to file %q: %s", furl, filename, err)
	}

	return filename, nil
}

func selectDecompressReader(src io.Reader, format string) (io.ReadCloser, error) {
	switch format {
	case compressionFormatBZIP2:
		return io.NopCloser(bzip2.NewReader(src)), nil
	case compressionFormatGZIP:
		return gzip.NewReader(src)
	case compressionFormatXZ:
		xzReader, err := xz.NewReader(src)
		if err != nil {
			return nil, err
		}
		return io.NopCloser(xzReader), nil
	}

	return nil, fmt.Errorf("format %s is not supported", format)
}

func getCompressionFormatFromContentType(contentType string) (string, error) {
	switch contentType {
	case "gzip", "application/gzip", "application/x-gzip":
		return compressionFormatGZIP, nil
	case "bzip2", "application/bzip2", "application/x-bzip2":
		return compressionFormatBZIP2, nil
	case "xz", "application/xz", "application/x-xz":
		return compressionFormatXZ, nil
	}
	return "", fmt.Errorf("content-type %s is not supported", contentType)
}

func selectUnzipReader(src io.Reader, format string) (io.ReadCloser, error) {
	if format == archivingFormatTAR {
		reader := tar.NewReader(src)
		for {
			header, err := reader.Next()
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				return nil, err
			}

			switch header.Typeflag {
			case tar.TypeReg:
				return io.NopCloser(reader), nil
			default:
				return nil, fmt.Errorf("got unexpected type: %s in %s", string(header.Typeflag), header.Name)
			}
		}
	}
	return nil, fmt.Errorf("format %s is not supported", format)
}

func resourceImagesImageRefreshFunc(client *gophercloud.ServiceClient, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		img, err := images.Get(client, id).Extract()
		if err != nil {
			return nil, "", err
		}
		log.Printf("[DEBUG] VKCS image status is: %s", img.Status)

		return img, string(img.Status), nil
	}
}

func resourceImagesImageBuildTags(v []interface{}) []string {
	tags := make([]string, len(v))
	for i, tag := range v {
		tags[i] = tag.(string)
	}

	return tags
}

func resourceImagesImageExpandProperties(v map[string]interface{}) map[string]string {
	properties := map[string]string{}
	for key, value := range v {
		if v, ok := value.(string); ok {
			properties[key] = v
		}
	}

	return properties
}

func resourceImagesImageNeedsDefaultStore(endpoint string) bool {
	endpointURL, _ := url.Parse(endpoint)
	hostname := endpointURL.Hostname()
	for _, mask := range imagesDefaultStoreEndpointMasks {
		matches, _ := regexp.MatchString(mask, hostname)
		if matches {
			return true
		}
	}
	return false
}

func resourceImagesImageUpdateComputedAttributes(_ context.Context, diff *schema.ResourceDiff, meta interface{}) error {
	if diff.HasChange("properties") {
		// Only check if the image has been created.
		if diff.Id() != "" {
			// Try to reconcile the properties set by the server
			// with the properties set by the user.
			//
			// old = user properties + server properties
			// new = user properties only
			o, n := diff.GetChange("properties")

			newProperties := resourceImagesImageExpandProperties(n.(map[string]interface{}))

			for oldKey, oldValue := range o.(map[string]interface{}) {
				if oldKey == "store" {
					if v, ok := oldValue.(string); ok {
						newProperties[oldKey] = v
					}
				}

				// direct_url is provided by some storage drivers.
				if oldKey == "direct_url" {
					if v, ok := oldValue.(string); ok {
						newProperties[oldKey] = v
					}
				}
			}

			// Set the diff to the newProperties
			//
			// If the user has changed properties, they will be caught at this
			// point, too.
			if err := diff.SetNew("properties", newProperties); err != nil {
				log.Printf("[DEBUG] unable set diff for properties key: %s", err)
			}
		}
	}

	return nil
}

// v - slice of images to filter
// p - field "properties" of schema.Resource from dataSourceImagesImageIDs
// or dataSourceImagesImage. If p is empty no filtering applies and the
// function returns the v.
func imagesFilterByProperties(v []images.Image, p map[string]string) []images.Image {
	var result []images.Image

	if len(p) > 0 {
		for _, image := range v {
			if len(image.Properties) > 0 {
				match := true
				for searchKey, searchValue := range p {
					imageValue, ok := image.Properties[searchKey]
					if !ok {
						match = false
						break
					}

					if searchValue != imageValue {
						match = false
						break
					}
				}

				if match {
					result = append(result, image)
				}
			}
		}
	} else {
		result = v
	}

	return result
}
