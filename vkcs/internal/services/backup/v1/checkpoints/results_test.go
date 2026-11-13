package checkpoints

import (
	"encoding/json"
	"net/url"
	"testing"
)

func TestListCheckpointsGroupedByPlanIDExtractsObjectLock(t *testing.T) {
	var response ListCheckpointsGroupedByPlanIDResponse
	err := json.Unmarshal([]byte(`{"project_id":"project-id","plan_id":"plan-id","resources":[{"id":"server-id","name":"server","resource_type":"OS::Nova::Server","backups":[{"checkpoint_id":"checkpoint-id","object_lock":{"type":"legalhold"}}]}]}`), &response)
	if err != nil {
		t.Fatalf("unmarshal grouped checkpoints response: %v", err)
	}

	if response.PlanID != "plan-id" {
		t.Fatalf("unexpected plan ID: %q", response.PlanID)
	}
	if len(response.Resources) != 1 || len(response.Resources[0].Backups) != 1 {
		t.Fatalf("unexpected response resources: %#v", response.Resources)
	}
	backup := response.Resources[0].Backups[0]
	if backup.CheckpointID != "checkpoint-id" {
		t.Fatalf("unexpected checkpoint ID: %q", backup.CheckpointID)
	}
	if backup.ObjectLock == nil || backup.ObjectLock.Type != "legalhold" {
		t.Fatalf("unexpected checkpoint object_lock: %#v", backup.ObjectLock)
	}
}

func TestListOptsEncodesPlanID(t *testing.T) {
	query, err := (ListOpts{PlanID: "plan id"}).ToListQuery()
	if err != nil {
		t.Fatalf("build query: %v", err)
	}

	parsed, err := url.Parse("https://example.com" + query)
	if err != nil {
		t.Fatalf("parse query: %v", err)
	}
	if got := parsed.Query().Get("plan_id"); got != "plan id" {
		t.Fatalf("unexpected plan_id query value: %q", got)
	}
}

func TestCheckpointResponseExtractsNovaResource(t *testing.T) {
	var response GetCheckpointResponse
	err := json.Unmarshal([]byte(`{"checkpoint":{"id":"checkpoint-id","resources":[{"id":"resource-id","name":"server","original_id":"server-id","type":"OS::Nova::Server"}]}}`), &response)
	if err != nil {
		t.Fatalf("unmarshal checkpoint response: %v", err)
	}
	if len(response.Checkpoint.Resources) != 1 {
		t.Fatalf("unexpected checkpoint resources: %#v", response.Checkpoint.Resources)
	}
	server := response.Checkpoint.Resources[0]
	if server.Type != "OS::Nova::Server" || server.OriginalID != "server-id" {
		t.Fatalf("unexpected checkpoint resource: %#v", server)
	}
}
