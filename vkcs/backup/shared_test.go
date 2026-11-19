package backup

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
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
