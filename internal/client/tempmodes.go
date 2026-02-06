package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type TempModeActions struct{ c *Client }

func (c *Client) TempModes() *TempModeActions { return &TempModeActions{c: c} }

func (t *TempModeActions) NapActivate(ctx context.Context) error {
	return t.simplePost(ctx, "/temperature/nap-mode/activate")
}

func (t *TempModeActions) NapDeactivate(ctx context.Context) error {
	return t.simplePost(ctx, "/temperature/nap-mode/deactivate")
}

func (t *TempModeActions) NapExtend(ctx context.Context) error {
	return t.simplePost(ctx, "/temperature/nap-mode/extend")
}

func (t *TempModeActions) NapStatus(ctx context.Context, out any) error {
	return t.simpleGet(ctx, "/temperature/nap-mode/status", out)
}

func (t *TempModeActions) HotFlashActivate(ctx context.Context) error {
	return t.simplePost(ctx, "/temperature/hot-flash-mode/activate")
}

func (t *TempModeActions) HotFlashDeactivate(ctx context.Context) error {
	return t.simplePost(ctx, "/temperature/hot-flash-mode/deactivate")
}

func (t *TempModeActions) HotFlashStatus(ctx context.Context, out any) error {
	return t.simpleGet(ctx, "/temperature/hot-flash-mode", out)
}

func (t *TempModeActions) TempEvents(ctx context.Context, from, to string, out any) error {
	if err := t.c.requireUser(ctx); err != nil {
		return err
	}
	q := url.Values{}
	if from != "" {
		q.Set("from", from)
	}
	if to != "" {
		q.Set("to", to)
	}
	path := fmt.Sprintf("/users/%s/temp-events", t.c.UserID)
	return t.c.doApp(ctx, http.MethodGet, path, q, nil, out)
}

func (t *TempModeActions) simplePost(ctx context.Context, suffix string) error {
	if err := t.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/users/%s%s", t.c.UserID, suffix)
	return t.c.doApp(ctx, http.MethodPost, path, nil, map[string]string{}, nil)
}

func (t *TempModeActions) simpleGet(ctx context.Context, suffix string, out any) error {
	if err := t.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/users/%s%s", t.c.UserID, suffix)
	return t.c.doApp(ctx, http.MethodGet, path, nil, nil, out)
}
