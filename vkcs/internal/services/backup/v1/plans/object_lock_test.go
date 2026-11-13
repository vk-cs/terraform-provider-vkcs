package plans

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestPlanCreateMapsObjectLock(t *testing.T) {
	days := int64(30)
	request := PlanCreate{
		Plan: &CreateOpts{
			Name:       "daily-backup",
			Resources:  []*BackupPlanResource{},
			ProviderID: "provider-id",
			ObjectLock: &ObjectLock{Type: "compliance", Days: &days},
		},
	}

	assertObjectLockRequest(t, request.Map, map[string]interface{}{
		"type": "compliance",
		"days": float64(30),
	})
}

func TestPlanUpdateMapsLegalHoldWithoutRetention(t *testing.T) {
	request := PlanUpdate{
		Plan: &UpdateOpts{
			Name:      "daily-backup",
			Status:    "running",
			Resources: []*BackupPlanResource{},
			ObjectLock: &ObjectLock{
				Type: "legalhold",
			},
		},
	}

	assertObjectLockRequest(t, request.Map, map[string]interface{}{
		"type": "legalhold",
	})
}

func TestPlanUpdateMapsRetentionEndTimestamp(t *testing.T) {
	retainUntil := int64(1900000000)
	request := PlanUpdate{
		Plan: &UpdateOpts{
			Name:       "daily-backup",
			Status:     "running",
			Resources:  []*BackupPlanResource{},
			ObjectLock: &ObjectLock{Type: "governance", RetainUntil: &retainUntil},
		},
	}

	assertObjectLockRequest(t, request.Map, map[string]interface{}{
		"type":         "governance",
		"retain_until": float64(1900000000),
	})
}

func TestPlanResponseExtractsObjectLock(t *testing.T) {
	var response PlanResp
	if err := json.Unmarshal([]byte(`{"plan":{"id":"plan-id","name":"daily-backup","resources":[],"provider_id":"provider-id","status":"running","object_lock":{"type":"governance","retain_until":1900000000}}}`), &response); err != nil {
		t.Fatalf("unmarshal plan response: %v", err)
	}
	if response.Plan.ObjectLock == nil || response.Plan.ObjectLock.Type != "governance" || response.Plan.ObjectLock.RetainUntil == nil || *response.Plan.ObjectLock.RetainUntil != 1900000000 {
		t.Fatalf("unexpected object lock in plan response: %#v", response.Plan.ObjectLock)
	}
}

func assertObjectLockRequest(t *testing.T, mapRequest func() (map[string]interface{}, error), want map[string]interface{}) {
	t.Helper()

	requestBody, err := mapRequest()
	if err != nil {
		t.Fatalf("map request: %v", err)
	}
	encoded, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("marshal request body: %v", err)
	}

	var decoded struct {
		Plan struct {
			ObjectLock map[string]interface{} `json:"object_lock"`
		} `json:"plan"`
	}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatalf("unmarshal request body: %v", err)
	}
	if !reflect.DeepEqual(decoded.Plan.ObjectLock, want) {
		t.Fatalf("unexpected object_lock body: got %#v, want %#v", decoded.Plan.ObjectLock, want)
	}
}
