package valid

import "testing"

func TestClusterName(t *testing.T) {
	tests := map[string]struct {
		name string
		err  error
	}{
		// errors
		"no name":              {name: "", err: ErrInvalidClusterName},
		"invalid first symbol": {name: "0abc", err: ErrInvalidClusterName},
		"invalid name":         {name: "abc-abc-/d", err: ErrInvalidClusterName},
		// ok
		"one symbol name": {name: "a", err: nil},
		"normal name":     {name: "abc-123_dc.22-", err: nil},
	}

	for name := range tests {
		tt := tests[name]
		t.Run(name, func(t *testing.T) {
			if err := ClusterName(tt.name); err != tt.err {
				t.Errorf("err got=%s; want=%s", err, tt.err)
			}
		})
	}
}
