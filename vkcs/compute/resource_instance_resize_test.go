package compute

import (
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
)

func TestInstanceResizeState(t *testing.T) {
	const targetFlavorID = "target-flavor"

	tests := []struct {
		name    string
		server  *servers.Server
		want    string
		wantErr string
	}{
		{
			name:   "resize in progress",
			server: &servers.Server{Status: "RESIZE", Flavor: map[string]interface{}{"id": "current-flavor"}},
			want:   instanceResizeStatePending,
		},
		{
			name:   "resize in progress with target flavor",
			server: &servers.Server{Status: "RESIZE", Flavor: map[string]interface{}{"id": targetFlavorID}},
			want:   instanceResizeStatePending,
		},
		{
			name:   "verify resize requires confirmation",
			server: &servers.Server{Status: "VERIFY_RESIZE", Flavor: map[string]interface{}{"id": targetFlavorID}},
			want:   instanceResizeStatePending,
		},
		{
			name:   "active with old flavor",
			server: &servers.Server{Status: "ACTIVE", Flavor: map[string]interface{}{"id": "current-flavor"}},
			want:   instanceResizeStatePending,
		},
		{
			name:   "active with target flavor",
			server: &servers.Server{Status: "ACTIVE", Flavor: map[string]interface{}{"id": targetFlavorID}},
			want:   instanceResizeStateDone,
		},
		{
			name:   "shutoff with target flavor",
			server: &servers.Server{Status: "SHUTOFF", Flavor: map[string]interface{}{"id": targetFlavorID}},
			want:   instanceResizeStateDone,
		},
		{
			name:   "shelved offloaded with target flavor",
			server: &servers.Server{Status: "SHELVED_OFFLOADED", Flavor: map[string]interface{}{"id": targetFlavorID}},
			want:   instanceResizeStateDone,
		},
		{
			name:   "shelved offloaded with old flavor",
			server: &servers.Server{Status: "SHELVED_OFFLOADED", Flavor: map[string]interface{}{"id": "current-flavor"}},
			want:   instanceResizeStatePending,
		},
		{
			name:    "error with fault message",
			server:  &servers.Server{Status: "ERROR", Fault: servers.Fault{Message: "resize failed"}},
			wantErr: "resize failed",
		},
		{
			name:    "error without fault message",
			server:  &servers.Server{Status: "ERROR"},
			wantErr: "instance entered ERROR state during resize",
		},
		{
			name:    "deleted",
			server:  &servers.Server{Status: "DELETED"},
			wantErr: "instance was deleted during resize",
		},
		{
			name:    "soft deleted",
			server:  &servers.Server{Status: "SOFT_DELETED"},
			wantErr: "instance was deleted during resize",
		},
		{
			name:    "revert resize",
			server:  &servers.Server{Status: "REVERT_RESIZE"},
			wantErr: "resize was reverted",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := instanceResizeState(tt.server, targetFlavorID)
			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error %q, got nil", tt.wantErr)
				}
				if err.Error() != tt.wantErr {
					t.Fatalf("expected error %q, got %q", tt.wantErr, err.Error())
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if got != tt.want {
				t.Fatalf("expected state %q, got %q", tt.want, got)
			}
		})
	}
}

func TestIsResizeAlreadyConfirmed(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "nil",
			err:  nil,
			want: false,
		},
		{
			name: "conflict",
			err: gophercloud.ErrDefault409{
				ErrUnexpectedResponseCode: gophercloud.ErrUnexpectedResponseCode{
					Actual: 409,
				},
			},
			want: true,
		},
		{
			name: "server error",
			err: gophercloud.ErrDefault500{
				ErrUnexpectedResponseCode: gophercloud.ErrUnexpectedResponseCode{
					Actual: 500,
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isResizeAlreadyConfirmed(tt.err); got != tt.want {
				t.Fatalf("expected %v, got %v", tt.want, got)
			}
		})
	}
}
