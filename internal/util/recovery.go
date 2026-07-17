package util

import (
	"context"
	"log/slog"
	"runtime/debug"
)
func SafeGo(ctx context.Context, name string, fn func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				slog.ErrorContext(ctx, "goroutine panic recovered",
					"goroutine", name,
					"error", r,
					"stack", string(debug.Stack()),
				)
			}
		}()
		
	}()
}