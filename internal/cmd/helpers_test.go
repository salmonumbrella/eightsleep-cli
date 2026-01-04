package cmd

import (
	"strings"
	"testing"
)

func TestFormatDays(t *testing.T) {
	tests := []struct {
		name string
		days []int
		want string
	}{
		{"empty", nil, ""},
		{"empty slice", []int{}, ""},
		{"sunday", []int{0}, "sun"},
		{"weekdays", []int{1, 2, 3, 4, 5}, "mon,tue,wed,thu,fri"},
		{"weekend", []int{0, 6}, "sun,sat"},
		{"all days", []int{0, 1, 2, 3, 4, 5, 6}, "sun,mon,tue,wed,thu,fri,sat"},
		{"out of range", []int{7, 8}, "7,8"},
		{"negative", []int{-1}, "-1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatDays(tt.days); got != tt.want {
				t.Errorf("formatDays(%v) = %q, want %q", tt.days, got, tt.want)
			}
		})
	}
}

func TestValidAlarmTime(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"07:30", true},
		{"23:59", true},
		{"00:00", true},
		{"07:30:00", true},
		{"7:30", true},   // single digit hour - Go's time.Parse accepts this
		{"07:3", false},  // single digit minute
		{"25:00", false}, // invalid hour
		{"07:60", false}, // invalid minute
		{"", false},
		{"abc", false},
		{"07-30", false}, // wrong separator
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := validAlarmTime(tt.input); got != tt.want {
				t.Errorf("validAlarmTime(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestSuggestAccountName(t *testing.T) {
	tests := []struct {
		email string
		want  string
	}{
		{"user@example.com", "user"},
		{"John.Doe@company.org", "john-doe"},
		{"first_last@domain.com", "first_last"},
		{"test+tag@mail.com", "test-tag"},
		{"UPPER@case.com", "upper"},
		{"", ""},
		{"@nolocal.com", ""},
		{"   spaces@trim.com   ", "spaces"},
		// Test 64-char truncation
		{strings.Repeat("a", 100) + "@example.com", strings.Repeat("a", 64)},
	}
	for _, tt := range tests {
		t.Run(tt.email, func(t *testing.T) {
			if got := suggestAccountName(tt.email); got != tt.want {
				t.Errorf("suggestAccountName(%q) = %q, want %q", tt.email, got, tt.want)
			}
		})
	}
}

func TestParseZoneinfoPath(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{"/usr/share/zoneinfo/America/Los_Angeles", "America/Los_Angeles"},
		{"/var/db/timezone/zoneinfo/Europe/London", "Europe/London"},
		{"/etc/localtime", ""}, // no zoneinfo in path
		{"", ""},
		{"zoneinfo/UTC", "UTC"},
		{"/some/path/zoneinfo/Asia/Tokyo", "Asia/Tokyo"},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			if got := parseZoneinfoPath(tt.path); got != tt.want {
				t.Errorf("parseZoneinfoPath(%q) = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}
