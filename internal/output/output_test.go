package output

import "testing"

func TestFilterFields(t *testing.T) {
	rows := []map[string]any{{"a": 1, "b": 2}, {"a": 3, "b": 4}}
	got := FilterFields(rows, []string{"b"})
	if len(got) != 2 {
		t.Fatalf("expected 2 rows")
	}
	if _, ok := got[0]["a"]; ok {
		t.Fatalf("field a should be filtered out")
	}
	if got[0]["b"].(int) != 2 || got[1]["b"].(int) != 4 {
		t.Fatalf("unexpected values: %+v", got)
	}
}
