//go:build !windows

package main

import "github.com/wailsapp/wails/v2/pkg/runtime"

func (a *App) FixLiveSplit() {
	runtime.EventsEmit(a.ctx, "ElMessage", "error", "LiveSplit 仅支持 Windows")
}
