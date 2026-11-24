package backup

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/backup/v1/plans"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/backup/v1/triggers"
)

func TestExpandTriggerScheduleEveryHours(t *testing.T) {
	tests := []struct {
		hours int64
		want  string
	}{
		{hours: 3, want: "0 */3 * * *"},
		{hours: 12, want: "0 */12 * * *"},
		{hours: 24, want: "0 0 * * *"},
	}

	for _, test := range tests {
		t.Run(test.want, func(t *testing.T) {
			plan := PlanResourceModel{
				Schedule: &PlanResourceScheduleModel{EveryHours: types.Int64Value(test.hours)},
			}

			got, err := expandTriggerSchedule(plan)
			if err != nil {
				t.Fatalf("expandTriggerSchedule() error = %v", err)
			}
			if got != test.want {
				t.Fatalf("expandTriggerSchedule() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestExpandAndFlattenTriggerScheduleEveryHours(t *testing.T) {
	tests := []struct {
		hours int64
	}{
		{hours: 3},
		{hours: 12},
		{hours: 24},
	}

	for _, test := range tests {
		t.Run(testHoursName(test.hours), func(t *testing.T) {
			plan := PlanResourceModel{
				Schedule: &PlanResourceScheduleModel{EveryHours: types.Int64Value(test.hours)},
			}
			pattern, err := expandTriggerSchedule(plan)
			if err != nil {
				t.Fatalf("expandTriggerSchedule() error = %v", err)
			}

			flattened := flattenSchedule(
				&plans.PlanResponse{},
				&triggers.TriggerResponse{Properties: &triggers.PropertiesOpts{Pattern: pattern}},
				time.UTC,
			)
			if flattened.EveryHours.IsNull() {
				t.Fatalf("every_hours round trip is null, want %d (pattern %q)", test.hours, pattern)
			}
			if flattened.EveryHours.ValueInt64() != test.hours {
				t.Fatalf("every_hours round trip = %d, want %d (pattern %q)", flattened.EveryHours.ValueInt64(), test.hours, pattern)
			}
			if !flattened.Time.IsNull() || flattened.Date != nil {
				t.Fatalf("daily interval flattened to time=%v date=%v, want time null and date nil", flattened.Time, flattened.Date)
			}
		})
	}
}

func TestFlattenScheduleDoesNotTreatIncrementalMidnightAs24Hours(t *testing.T) {
	fullDay := 0
	flattened := flattenSchedule(
		&plans.PlanResponse{FullDay: &fullDay},
		&triggers.TriggerResponse{Properties: &triggers.PropertiesOpts{Pattern: "0 0 * * *"}},
		time.UTC,
	)
	if !flattened.EveryHours.IsNull() {
		t.Fatalf("incremental midnight schedule has every_hours=%d, want null", flattened.EveryHours.ValueInt64())
	}
	if flattened.Time.ValueString() != "00:00" {
		t.Fatalf("incremental midnight schedule time = %q, want 00:00", flattened.Time.ValueString())
	}
	if len(flattened.Date) != 1 || flattened.Date[0].ValueString() != "Mo" {
		t.Fatalf("incremental midnight schedule days = %v, want [Mo]", flattened.Date)
	}
}

func testHoursName(hours int64) string {
	return fmt.Sprintf("every_%d_hours", hours)
}
