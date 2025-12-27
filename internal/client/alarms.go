package client

import (
	"context"
	"fmt"
	"net/http"
)

// AlarmActions groups alarm endpoints.
type AlarmActions struct {
	c *Client
}

// Alarms helper accessor.
func (c *Client) Alarms() *AlarmActions { return &AlarmActions{c: c} }

func (a *AlarmActions) Snooze(ctx context.Context, alarmID string) error {
	if err := a.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/users/%s/alarms/%s/snooze", a.c.UserID, alarmID)
	return a.c.do(ctx, http.MethodPost, path, nil, map[string]string{}, nil)
}

func (a *AlarmActions) Dismiss(ctx context.Context, alarmID string) error {
	if err := a.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/users/%s/alarms/%s/dismiss", a.c.UserID, alarmID)
	return a.c.do(ctx, http.MethodPost, path, nil, map[string]string{}, nil)
}

func (a *AlarmActions) DismissAll(ctx context.Context) error {
	if err := a.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/users/%s/alarms/active/dismiss-all", a.c.UserID)
	return a.c.do(ctx, http.MethodPost, path, nil, map[string]string{}, nil)
}

func (a *AlarmActions) VibrationTest(ctx context.Context) error {
	if err := a.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/users/%s/vibration-test", a.c.UserID)
	return a.c.do(ctx, http.MethodPost, path, nil, map[string]string{}, nil)
}
