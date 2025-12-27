package client

import (
	"context"
	"fmt"
	"net/http"
)

type DeviceActions struct{ c *Client }

func (c *Client) Device() *DeviceActions { return &DeviceActions{c: c} }

func (d *DeviceActions) Info(ctx context.Context) (any, error) {
	id, err := d.c.EnsureDeviceID(ctx)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/devices/%s", id)
	var res any
	err = d.c.do(ctx, http.MethodGet, path, nil, nil, &res)
	return res, err
}

func (d *DeviceActions) Peripherals(ctx context.Context) (any, error) {
	id, err := d.c.EnsureDeviceID(ctx)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/devices/%s/peripherals", id)
	var res any
	err = d.c.do(ctx, http.MethodGet, path, nil, nil, &res)
	return res, err
}

func (d *DeviceActions) Owner(ctx context.Context) (any, error) {
	id, err := d.c.EnsureDeviceID(ctx)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/devices/%s/owner", id)
	var res any
	err = d.c.do(ctx, http.MethodGet, path, nil, nil, &res)
	return res, err
}

func (d *DeviceActions) Warranty(ctx context.Context) (any, error) {
	id, err := d.c.EnsureDeviceID(ctx)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/devices/%s/warranty", id)
	var res any
	err = d.c.do(ctx, http.MethodGet, path, nil, nil, &res)
	return res, err
}

func (d *DeviceActions) Online(ctx context.Context) (any, error) {
	id, err := d.c.EnsureDeviceID(ctx)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/devices/%s/online", id)
	var res any
	err = d.c.do(ctx, http.MethodGet, path, nil, nil, &res)
	return res, err
}

func (d *DeviceActions) PrimingTasks(ctx context.Context) (any, error) {
	id, err := d.c.EnsureDeviceID(ctx)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/devices/%s/priming/tasks", id)
	var res any
	err = d.c.do(ctx, http.MethodGet, path, nil, nil, &res)
	return res, err
}

func (d *DeviceActions) PrimingSchedule(ctx context.Context) (any, error) {
	id, err := d.c.EnsureDeviceID(ctx)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/devices/%s/priming/schedule", id)
	var res any
	err = d.c.do(ctx, http.MethodGet, path, nil, nil, &res)
	return res, err
}
