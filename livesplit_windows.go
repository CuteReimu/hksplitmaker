//go:build windows

package main

import (
	_ "embed"
	"errors"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed LiveSplit.AutoSplitters.xml
var liveSplitAutoSplittersXml []byte

//go:embed LiveSplit.HollowKnight.dll
var wasmFile []byte

var ErrFixLiveSplitIgnore = errors.New("fix live split ignore")

func (a *App) FixLiveSplit() {
	err := a.fixLiveSplit()
	if err == nil {
		_, _ = runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.InfoDialog,
			Title:   "成功",
			Message: "修复成功，建议使用管理员模式打开计时器",
		})
	} else if !errors.Is(err, ErrFixLiveSplitIgnore) {
		_, _ = runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Title:   "错误",
			Message: err.Error(),
		})
	}
}

func (a *App) fixLiveSplit() error {
	file, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "请选择LiveSplit.exe",
		Filters: []runtime.FileFilter{
			{DisplayName: "LiveSplit.exe", Pattern: "LiveSplit.exe"},
		},
	})
	if err != nil {
		return err
	}
	if file == "" {
		return ErrFixLiveSplitIgnore
	}
	if filepath.Base(file) != "LiveSplit.exe" {
		return errors.New("您选择的并不是LiveSplit.exe")
	}
	dir := filepath.Dir(file)
	if _, err := os.Stat(filepath.Join(dir, "Components")); err != nil {
		return errors.New("计时器目录似乎有些问题，无法一键修复")
	}
	if err := os.WriteFile(filepath.Join(dir, "LiveSplit.AutoSplitters.xml"), liveSplitAutoSplittersXml, 0644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "Components", "LiveSplit.HollowKnight.dll"), wasmFile, 0644); err != nil {
		return err
	}
	return nil
}
