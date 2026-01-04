package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// Alarm represents alarm payload (flattened from routines).
type Alarm struct {
	ID          string  `json:"id"`
	RoutineID   string  `json:"routineId,omitempty"`
	RoutineName string  `json:"routineName,omitempty"`
	Enabled     bool    `json:"enabled"`
	Time        string  `json:"time"`
	DaysOfWeek  []int   `json:"daysOfWeek"`
	Vibration   bool    `json:"vibration"`
	Sound       *string `json:"sound,omitempty"`
}

type Routine struct {
	ID       string           `json:"id"`
	Name     string           `json:"name,omitempty"`
	Days     []int            `json:"days,omitempty"`
	Alarms   []RoutineAlarm   `json:"alarms"`
	Override *RoutineOverride `json:"override,omitempty"`
}

type RoutineOverride struct {
	RoutineEnabled bool           `json:"routineEnabled"`
	Alarms         []RoutineAlarm `json:"alarms"`
}

type RoutineAlarm struct {
	AlarmID              string                `json:"alarmId"`
	Enabled              bool                  `json:"enabled"`
	DisabledIndividually bool                  `json:"disabledIndividually,omitempty"`
	EnabledSince         string                `json:"enabledSince,omitempty"`
	Time                 string                `json:"time,omitempty"`
	TimeWithOffset       *RoutineAlarmTime     `json:"timeWithOffset,omitempty"`
	Settings             *RoutineAlarmSettings `json:"settings,omitempty"`
}

type RoutineAlarmTime struct {
	Time string `json:"time"`
}

type RoutineAlarmSettings struct {
	Vibration *RoutineAlarmVibration `json:"vibration,omitempty"`
	Thermal   *RoutineAlarmThermal   `json:"thermal,omitempty"`
}

type RoutineAlarmVibration struct {
	Enabled    bool   `json:"enabled"`
	PowerLevel int    `json:"powerLevel,omitempty"`
	Pattern    string `json:"pattern,omitempty"`
}

type RoutineAlarmThermal struct {
	Enabled bool `json:"enabled"`
	Level   int  `json:"level,omitempty"`
}

type RoutineNextAlarm struct {
	AlarmID       string `json:"alarmId"`
	NextTimestamp string `json:"nextTimestamp"`
}

type RoutineState struct {
	NextAlarm         *RoutineNextAlarm `json:"nextAlarm,omitempty"`
	UpcomingRoutineID string            `json:"upcomingRoutineId,omitempty"`
}

type RoutinesState struct {
	Routines []Routine
	State    RoutineState
}

type routinesResponse struct {
	Settings struct {
		Routines []Routine `json:"routines"`
	} `json:"settings"`
	State RoutineState `json:"state"`
}

type OneOffAlarm struct {
	Time             string
	Enabled          bool
	VibrationEnabled bool
	VibrationLevel   int
	VibrationPattern string
	ThermalEnabled   bool
	ThermalLevel     int
}

func (c *Client) ListRoutines(ctx context.Context) (*RoutinesState, error) {
	if err := c.requireUser(ctx); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/v2/users/%s/routines", c.UserID)
	var res routinesResponse
	if err := c.doApp(ctx, http.MethodGet, path, nil, nil, &res); err != nil {
		return nil, err
	}
	return &RoutinesState{Routines: res.Settings.Routines, State: res.State}, nil
}

func (c *Client) UpdateRoutine(ctx context.Context, routineID string, routine Routine) error {
	if err := c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/v2/users/%s/routines/%s", c.UserID, routineID)
	return c.doApp(ctx, http.MethodPut, path, nil, routine, nil)
}

func (c *Client) SetOneOffAlarm(ctx context.Context, alarm OneOffAlarm) error {
	if err := c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/v2/users/%s/routines", c.UserID)
	body := map[string]any{
		"oneOffAlarms": []map[string]any{
			{
				"time":    alarm.Time,
				"enabled": alarm.Enabled,
				"settings": map[string]any{
					"vibration": map[string]any{
						"enabled":    alarm.VibrationEnabled,
						"powerLevel": alarm.VibrationLevel,
						"pattern":    alarm.VibrationPattern,
					},
					"thermal": map[string]any{
						"enabled": alarm.ThermalEnabled,
						"level":   alarm.ThermalLevel,
					},
				},
			},
		},
	}
	q := url.Values{}
	q.Set("ignoreDeviceErrors", "false")
	return c.doApp(ctx, http.MethodPut, path, q, body, nil)
}

func (c *Client) ListAlarms(ctx context.Context) ([]Alarm, error) {
	routinesState, err := c.ListRoutines(ctx)
	if err != nil {
		return nil, err
	}
	routines := routinesState.Routines
	alarms := []Alarm{}
	seen := map[string]bool{}
	for _, r := range routines {
		for _, a := range r.Alarms {
			effective := routineEffectiveAlarm(r, a)
			alarms = append(alarms, Alarm{
				ID:          effective.AlarmID,
				RoutineID:   r.ID,
				RoutineName: r.Name,
				Enabled:     routineAlarmEnabled(effective),
				Time:        routineAlarmTime(effective),
				DaysOfWeek:  r.Days,
				Vibration:   routineAlarmVibration(effective),
			})
			seen[effective.AlarmID] = true
		}
		if r.Override != nil {
			for _, a := range r.Override.Alarms {
				if seen[a.AlarmID] {
					continue
				}
				alarms = append(alarms, Alarm{
					ID:          a.AlarmID,
					RoutineID:   r.ID,
					RoutineName: r.Name,
					Enabled:     routineAlarmEnabled(a),
					Time:        routineAlarmTime(a),
					DaysOfWeek:  r.Days,
					Vibration:   routineAlarmVibration(a),
				})
				seen[a.AlarmID] = true
			}
		}
	}
	return alarms, nil
}

func routineEffectiveAlarm(r Routine, alarm RoutineAlarm) RoutineAlarm {
	if r.Override == nil {
		return alarm
	}
	for _, a := range r.Override.Alarms {
		if a.AlarmID == alarm.AlarmID {
			return a
		}
	}
	return alarm
}

func routineAlarmTime(a RoutineAlarm) string {
	if a.Time != "" {
		return a.Time
	}
	if a.TimeWithOffset != nil {
		return a.TimeWithOffset.Time
	}
	return ""
}

func routineAlarmEnabled(a RoutineAlarm) bool {
	if a.DisabledIndividually {
		return false
	}
	return a.Enabled
}

func routineAlarmVibration(a RoutineAlarm) bool {
	if a.Settings == nil || a.Settings.Vibration == nil {
		return false
	}
	return a.Settings.Vibration.Enabled
}
