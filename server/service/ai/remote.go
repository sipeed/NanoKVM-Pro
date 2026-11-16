package ai

import (
	"context"
	"errors"
	"fmt"
	"time"

	"NanoKVM-Server/service/hid"
)

type HIDWriter interface {
	WriteKeyboardReport([]byte) error
	WriteRelativeMouseReport([]byte) error
	WriteAbsoluteMouseReport([]byte) error
}

type Remote struct {
	action func(func(HIDWriter) error) error
	gate   chan struct{}
}

func NewRemote() *Remote {
	return &Remote{gate: make(chan struct{}, 1), action: func(action func(HIDWriter) error) error {
		return hid.GetHid().WithInputReports(func(writer *hid.ReportWriter) error { return action(writer) })
	}}
}

func waitInput(ctx context.Context, delay time.Duration) error {
	if delay <= 0 {
		return ctx.Err()
	}
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func (r *Remote) run(ctx context.Context, action func(context.Context, HIDWriter) error) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	select {
	case r.gate <- struct{}{}:
		defer func() { <-r.gate }()
	case <-ctx.Done():
		return ctx.Err()
	}
	return r.action(func(writer HIDWriter) error {
		if err := ctx.Err(); err != nil {
			return err
		}
		return action(ctx, writer)
	})
}

func (r *Remote) TypeText(ctx context.Context, text string, delay time.Duration) ([]rune, error) {
	reports, skipped := buildTypeReports(text)
	err := r.run(ctx, func(ctx context.Context, writer HIDWriter) (err error) {
		defer func() { err = errors.Join(err, writer.WriteKeyboardReport(keyUpReport)) }()
		for _, report := range reports {
			if err := ctx.Err(); err != nil {
				return err
			}
			if err := writer.WriteKeyboardReport(report); err != nil {
				return err
			}
			if err := waitInput(ctx, delay); err != nil {
				return err
			}
		}
		return nil
	})
	return skipped, err
}

func (r *Remote) PressKeys(ctx context.Context, keys []string, hold time.Duration) error {
	report, err := buildKeyComboReport(keys)
	if err != nil {
		return err
	}
	return r.run(ctx, func(ctx context.Context, writer HIDWriter) (err error) {
		defer func() { err = errors.Join(err, writer.WriteKeyboardReport(keyUpReport)) }()
		if err := writer.WriteKeyboardReport(report); err != nil {
			return err
		}
		return waitInput(ctx, hold)
	})
}

func (r *Remote) MoveAbsolute(ctx context.Context, x, y float64) error {
	return r.run(ctx, func(_ context.Context, writer HIDWriter) error {
		return writer.WriteAbsoluteMouseReport(buildAbsolutePointerReport(absoluteCoordinate(x), absoluteCoordinate(y), false, false))
	})
}

func (r *Remote) MoveRelative(ctx context.Context, x, y int) error {
	reports, err := buildRelativeMoveReports(x, y)
	if err != nil {
		return err
	}
	return r.run(ctx, func(ctx context.Context, writer HIDWriter) error {
		for _, report := range reports {
			if err := ctx.Err(); err != nil {
				return err
			}
			if err := writer.WriteRelativeMouseReport(report); err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *Remote) Click(ctx context.Context, button string, clicks int, absolute bool, x, y float64) error {
	bit, err := mouseButtonBit(button)
	if err != nil {
		return err
	}
	return r.run(ctx, func(ctx context.Context, writer HIDWriter) (err error) {
		write := func(buttons byte) error {
			if absolute {
				report := buildAbsolutePointerReport(absoluteCoordinate(x), absoluteCoordinate(y), false, false)
				report[0] = buttons
				return writer.WriteAbsoluteMouseReport(report)
			}
			return writer.WriteRelativeMouseReport(buildRelativeMouseReport(0, 0, buttons, 0))
		}
		defer func() { err = errors.Join(err, write(0)) }()
		for i := 0; i < clampInt(clicks, 1, 5); i++ {
			if err := ctx.Err(); err != nil {
				return err
			}
			if err := write(bit); err != nil {
				return err
			}
			if err := waitInput(ctx, 50*time.Millisecond); err != nil {
				return err
			}
			if err := write(0); err != nil {
				return err
			}
			if err := waitInput(ctx, 50*time.Millisecond); err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *Remote) Scroll(ctx context.Context, delta int) error {
	if reportsForDelta(delta) > maxRelativeMouseReports {
		return fmt.Errorf("scroll exceeds %d reports", maxRelativeMouseReports)
	}
	return r.run(ctx, func(ctx context.Context, writer HIDWriter) error {
		for delta != 0 {
			if err := ctx.Err(); err != nil {
				return err
			}
			wheel := clampInt(delta, -127, 127)
			if err := writer.WriteRelativeMouseReport(buildRelativeMouseReport(0, 0, 0, wheel)); err != nil {
				return err
			}
			delta -= wheel
		}
		return nil
	})
}
