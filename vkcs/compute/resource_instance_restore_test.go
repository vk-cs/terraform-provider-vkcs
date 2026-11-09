package compute

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/backup/v1/restores"
)

func TestRestoredComputeInstanceID(t *testing.T) {
	tests := []struct {
		name          string
		restore       *restores.Restore
		wantID        string
		wantReady     bool
		wantErrSubstr string
	}{
		{
			name: "restore operation is still pending",
			restore: &restores.Restore{
				ID:              "restore-1",
				ResourcesStatus: map[string]string{"OS::Nova::Server#server-1": "restoring"},
			},
		},
		{
			name: "returns the created Nova server ID",
			restore: &restores.Restore{
				ID:              "restore-2",
				ResourcesStatus: map[string]string{"OS::Nova::Server#server-2": "available"},
			},
			wantID:    "server-2",
			wantReady: true,
		},
		{
			name: "reports a failed restore reason",
			restore: &restores.Restore{
				ID:              "restore-3",
				ResourcesStatus: map[string]string{"OS::Nova::Server#server-3": "error"},
				ResourcesReason: map[string]string{"OS::Nova::Server#server-3": "quota exceeded"},
			},
			wantErrSubstr: "quota exceeded",
		},
		{
			name: "reports an error even when another resource is still pending",
			restore: &restores.Restore{
				ID: "restore-8",
				ResourcesStatus: map[string]string{
					"OS::Cinder::Volume#volume-8": "restoring",
					"OS::Nova::Server#server-8":   "error",
				},
				ResourcesReason: map[string]string{"OS::Nova::Server#server-8": "server restore failed"},
			},
			wantErrSubstr: "server restore failed",
		},
		{
			name: "rejects a restore containing multiple VMs",
			restore: &restores.Restore{
				ID: "restore-4",
				ResourcesStatus: map[string]string{
					"OS::Nova::Server#server-4": "available",
					"OS::Nova::Server#server-5": "available",
				},
			},
			wantErrSubstr: "multiple Nova servers",
		},
		{
			name: "rejects a completed restore without a Nova server",
			restore: &restores.Restore{
				ID:              "restore-6",
				ResourcesStatus: map[string]string{"OS::Cinder::Volume#volume-6": "available"},
			},
			wantErrSubstr: "without a restored Nova server",
		},
		{
			name:    "waits when restore statuses are not populated yet",
			restore: &restores.Restore{ID: "restore-7"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotID, gotReady, err := restoredComputeInstanceID(test.restore)
			if test.wantErrSubstr != "" {
				if err == nil || !strings.Contains(err.Error(), test.wantErrSubstr) {
					t.Fatalf("restoredComputeInstanceID() error = %v, want substring %q", err, test.wantErrSubstr)
				}
				return
			}
			if err != nil {
				t.Fatalf("restoredComputeInstanceID() error = %v", err)
			}
			if gotID != test.wantID {
				t.Errorf("restoredComputeInstanceID() ID = %q, want %q", gotID, test.wantID)
			}
			if gotReady != test.wantReady {
				t.Errorf("restoredComputeInstanceID() ready = %t, want %t", gotReady, test.wantReady)
			}
		})
	}
}

func TestComputeInstanceRestoreIDSchema(t *testing.T) {
	restoreID := ResourceComputeInstance().Schema["restore_id"]
	if restoreID == nil {
		t.Fatal("restore_id schema attribute is missing")
	}
	if !restoreID.Optional || !restoreID.ForceNew {
		t.Fatalf("restore_id must be optional and force replacement, got Optional=%t ForceNew=%t", restoreID.Optional, restoreID.ForceNew)
	}

	conflicts := make(map[string]bool, len(restoreID.ConflictsWith))
	for _, attribute := range restoreID.ConflictsWith {
		conflicts[attribute] = true
	}
	for _, attribute := range []string{
		"image_id", "image_name", "flavor_id", "flavor_name", "availability_zone", "network", "network_mode",
		"block_device", "key_pair", "security_groups", "security_group_ids", "user_data",
	} {
		if !conflicts[attribute] {
			t.Errorf("restore_id must conflict with %q", attribute)
		}
	}
}

func TestRestoredComputeInstanceWaitState(t *testing.T) {
	tests := []struct {
		name          string
		state         string
		instanceFound bool
		wantState     string
		wantErrSubstr string
	}{
		{name: "active VM is ready", state: "ACTIVE", instanceFound: true, wantState: "ACTIVE"},
		{name: "build waits", state: "BUILD", instanceFound: true, wantState: "BUILD"},
		{name: "provisioning waits", state: "PROVISIONING", instanceFound: true, wantState: "BUILD"},
		{name: "shutoff waits for active state", state: "SHUTOFF", instanceFound: true, wantState: "BUILD"},
		{name: "not-yet-visible VM waits", state: "DELETED", wantState: "BUILD"},
		{name: "deleted VM fails", state: "DELETED", instanceFound: true, wantErrSubstr: "deleted before becoming active"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotState, err := restoredComputeInstanceWaitState(test.state, test.instanceFound)
			if test.wantErrSubstr != "" {
				if err == nil || !strings.Contains(err.Error(), test.wantErrSubstr) {
					t.Fatalf("restoredComputeInstanceWaitState() error = %v, want substring %q", err, test.wantErrSubstr)
				}
				return
			}
			if err != nil {
				t.Fatalf("restoredComputeInstanceWaitState() error = %v", err)
			}
			if gotState != test.wantState {
				t.Errorf("restoredComputeInstanceWaitState(%q, %t) = %q, want %q", test.state, test.instanceFound, gotState, test.wantState)
			}
		})
	}
}

func TestRestoredComputeInstanceStateRefreshFunc(t *testing.T) {
	t.Run("waits while the restored VM is not visible yet", func(t *testing.T) {
		var missingServer *servers.Server
		refresh := restoredComputeInstanceStateRefreshFunc(func() (interface{}, string, error) {
			return missingServer, "DELETED", nil
		}, "server-1")

		_, state, err := refresh()
		if err != nil {
			t.Fatalf("restored refresh returned error: %v", err)
		}
		if state != "BUILD" {
			t.Fatalf("restored refresh state = %q, want BUILD while Nova is not yet visible", state)
		}
	})

	t.Run("fails when a visible restored VM has been deleted", func(t *testing.T) {
		refresh := restoredComputeInstanceStateRefreshFunc(func() (interface{}, string, error) {
			return &servers.Server{ID: "server-2"}, "DELETED", nil
		}, "server-2")

		_, _, err := refresh()
		if err == nil || !strings.Contains(err.Error(), "deleted before becoming active") {
			t.Fatalf("restored refresh error = %v, want deleted VM error", err)
		}
	})
}

func TestWaitForRestoredComputeInstanceID(t *testing.T) {
	t.Run("polls until a VM is available", func(t *testing.T) {
		responses := []*restores.Restore{
			{
				ID:              "restore-1",
				ResourcesStatus: map[string]string{"OS::Nova::Server#server-1": "restoring"},
			},
			{
				ID:              "restore-1",
				ResourcesStatus: map[string]string{"OS::Nova::Server#server-1": "available"},
			},
		}
		calls := 0
		getter := func(_ context.Context, restoreID string) (*restores.Restore, error) {
			if restoreID != "restore-1" {
				t.Fatalf("restore ID = %q, want %q", restoreID, "restore-1")
			}
			responseIndex := calls
			calls++
			if responseIndex >= len(responses) {
				responseIndex = len(responses) - 1
			}
			return responses[responseIndex], nil
		}

		gotID, err := waitForRestoredComputeInstanceID(context.Background(), "restore-1", time.Second, time.Millisecond, getter)
		if err != nil {
			t.Fatalf("waitForRestoredComputeInstanceID() error = %v", err)
		}
		if gotID != "server-1" {
			t.Fatalf("waitForRestoredComputeInstanceID() ID = %q, want %q", gotID, "server-1")
		}
		if calls != 2 {
			t.Fatalf("restore getter called %d times, want 2", calls)
		}
	})

	t.Run("returns API errors", func(t *testing.T) {
		wantErr := errors.New("backup API unavailable")
		getter := func(context.Context, string) (*restores.Restore, error) {
			return nil, wantErr
		}

		_, err := waitForRestoredComputeInstanceID(context.Background(), "restore-2", time.Second, time.Millisecond, getter)
		if !errors.Is(err, wantErr) {
			t.Fatalf("waitForRestoredComputeInstanceID() error = %v, want wrapped %v", err, wantErr)
		}
	})

	t.Run("times out while the restore is pending", func(t *testing.T) {
		getter := func(_ context.Context, restoreID string) (*restores.Restore, error) {
			return &restores.Restore{
				ID:              restoreID,
				ResourcesStatus: map[string]string{"OS::Nova::Server#server-3": "restoring"},
			}, nil
		}

		_, err := waitForRestoredComputeInstanceID(context.Background(), "restore-3", time.Millisecond, time.Hour, getter)
		if err == nil || !strings.Contains(err.Error(), "restore-3") {
			t.Fatalf("waitForRestoredComputeInstanceID() error = %v, want timeout mentioning restore ID", err)
		}
	})
}
