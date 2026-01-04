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

func (a *AlarmActions) Snooze(ctx context.Context, alarmID string, minutes int) error {
	if err := a.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/v1/users/%s/routines", a.c.UserID)
	body := map[string]any{
		"alarm": map[string]any{
			"alarmId":          alarmID,
			"snoozeForMinutes": minutes,
		},
	}
	return a.c.doApp(ctx, http.MethodPut, path, nil, body, nil)
}

func (a *AlarmActions) Dismiss(ctx context.Context, alarmID string) error {
	if err := a.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/v1/users/%s/routines", a.c.UserID)
	body := map[string]any{
		"alarm": map[string]any{
			"alarmId":   alarmID,
			"dismissed": true,
		},
	}
	return a.c.doApp(ctx, http.MethodPut, path, nil, body, nil)
}

func (a *AlarmActions) DismissAll(ctx context.Context) error {
	if err := a.c.requireUser(ctx); err != nil {
		return err
	}
	state, err := a.c.ListRoutines(ctx)
	if err != nil {
		return err
	}
	if state.State.NextAlarm == nil || state.State.NextAlarm.AlarmID == "" {
		return fmt.Errorf("no active alarm to dismiss")
	}
	return a.Dismiss(ctx, state.State.NextAlarm.AlarmID)
}

func (a *AlarmActions) VibrationTest(ctx context.Context) error {
	if err := a.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/users/%s/vibration-test", a.c.UserID)
	return a.c.do(ctx, http.MethodPost, path, nil, map[string]string{}, nil)
}
