package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func resolveTimezone(name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" || strings.EqualFold(name, "local") {
		return localTimezoneName()
	}
	if _, err := time.LoadLocation(name); err != nil {
		return "", fmt.Errorf("invalid timezone %q: %w", name, err)
	}
	return name, nil
}

func localTimezoneName() (string, error) {
	if tz := strings.TrimSpace(os.Getenv("TZ")); tz != "" && !strings.EqualFold(tz, "local") {
		if _, err := time.LoadLocation(tz); err == nil {
			return tz, nil
		}
	}

	if name := time.Now().Location().String(); name != "" && name != "Local" {
		if _, err := time.LoadLocation(name); err == nil {
			return name, nil
		}
	}

	if tz := timezoneFromLocaltimeLink(); tz != "" {
		return tz, nil
	}

	if data, err := os.ReadFile("/etc/timezone"); err == nil {
		if tz := strings.TrimSpace(string(data)); tz != "" {
			if _, err := time.LoadLocation(tz); err == nil {
				return tz, nil
			}
		}
	}

	return "", fmt.Errorf("could not resolve local timezone to IANA name")
}

func timezoneFromLocaltimeLink() string {
	link, err := os.Readlink("/etc/localtime")
	if err != nil {
		link, err = filepath.EvalSymlinks("/etc/localtime")
		if err != nil {
			return ""
		}
	}
	return parseZoneinfoPath(link)
}

func parseZoneinfoPath(path string) string {
	const marker = "zoneinfo/"
	idx := strings.Index(path, marker)
	if idx == -1 {
		return ""
	}
	tz := filepath.ToSlash(path[idx+len(marker):])
	tz = strings.TrimPrefix(tz, "/")
	return tz
}
