package cmd

import "testing"

func TestParseSchedule(t *testing.T) {
	data := []byte(`
schedule:
  - time: "22:00"
    action: "temp"
    temperature: "68F"
`)
	items, err := parseSchedule(data)
	if err != nil {
		t.Fatalf("parseSchedule: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
}

func TestParseScheduleEmpty(t *testing.T) {
	if _, err := parseSchedule([]byte(`schedule: []`)); err == nil {
		t.Fatalf("expected error for empty schedule")
	}
}

func TestDefaultPIDFile(t *testing.T) {
	if got := defaultPIDFile("/tmp/pid"); got != "/tmp/pid" {
		t.Fatalf("expected explicit pid path")
	}
}
