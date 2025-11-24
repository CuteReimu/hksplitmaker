package main

import (
	_ "embed"
	"os"
	"path/filepath"

	"github.com/lxn/walk"
	"github.com/lxn/win"
)

//go:embed LiveSplit.AutoSplitters.xml
var liveSplitAutoSplittersXml []byte

//go:embed LiveSplit.HollowKnight.dll
var liveSplitHollowKnightDll []byte

func fixLiveSplit() {
	walk.MsgBox(mainWindow, "提示", "请选择LiveSplit.exe", walk.MsgBoxIconInformation)
	dlg := new(walk.FileDialog)
	dlg.Title = "请选择LiveSplit.exe"
	dlg.Filter = "LiveSplit.exe|LiveSplit.exe"
	dlg.Flags = win.OFN_FILEMUSTEXIST
	if ok, err := dlg.ShowOpen(mainWindow); err != nil {
		walk.MsgBox(mainWindow, "内部错误", err.Error(), walk.MsgBoxIconError)
	} else if ok {
		file := dlg.FilePath
		if filepath.Base(file) != "LiveSplit.exe" {
			walk.MsgBox(mainWindow, "错误", "您选择的并不是LiveSplit.exe", walk.MsgBoxIconError)
			return
		}
		dir := filepath.Dir(file)
		_, err := os.Stat(filepath.Join(dir, "Components"))
		if err != nil {
			walk.MsgBox(mainWindow, "错误", "计时器目录似乎有些问题，无法一键修复", walk.MsgBoxIconError)
			return
		}
		err = os.WriteFile(filepath.Join(dir, "LiveSplit.AutoSplitters.xml"), liveSplitAutoSplittersXml, 0644)
		if err != nil {
			walk.MsgBox(mainWindow, "内部错误", err.Error(), walk.MsgBoxIconError)
			return
		}
		err = os.WriteFile(filepath.Join(dir, "Components", "LiveSplit.HollowKnight.dll"), liveSplitHollowKnightDll, 0644)
		if err != nil {
			walk.MsgBox(mainWindow, "内部错误", err.Error(), walk.MsgBoxIconError)
			return
		}
		walk.MsgBox(mainWindow, "提示", "修复成功，建议使用管理员模式打开计时器", walk.MsgBoxIconInformation)
	}
}
