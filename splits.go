package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"errors"
	"io"
	"os"
	"regexp"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed hk-split-maker/src/asset/hollowknight/splits.txt
var splitsFileBuf []byte
var splitsDictIdToDescriptions = make(map[string]string)

func initSplitsFile() error {
	rd := bufio.NewReader(bytes.NewReader(splitsFileBuf))
	re, err := regexp.Compile(`\[Description\("(.*?)"\)\s*,\s*ToolTip\("(.*?)"\)]`)
	if err != nil {
		return err
	}
	var isNameLine bool
	line, isPrefix, err := rd.ReadLine()
	var result []string
	for ; err == nil; line, isPrefix, err = rd.ReadLine() {
		if isPrefix {
			err = errors.New("尚未支持这种文件，尽情期待更新")
			break
		}
		line = bytes.Trim(bytes.TrimSpace(line), ",")
		if len(line) == 0 {
			continue
		}
		if isNameLine {
			if len(result) == 3 {
				description := translate(result[1])
				splitsDictIdToDescriptions[string(line)] = description
				isNameLine = false
			} else {
				err = errors.New("splits.txt文件格式错误")
				break
			}
		} else {
			result = re.FindStringSubmatch(string(line))
			if result == nil {
				err = errors.New("splits.txt文件格式错误")
				break
			}
			isNameLine = true
		}
	}
	if errors.Is(err, io.EOF) {
		return nil
	}
	return err
}

//go:embed blank-colo-save.dat
var blankColoSave []byte

func (a *App) coloNotice(...any) {
	if v, err := runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Type:          runtime.QuestionDialog,
		Title:         "提示",
		Message:       "愚人斗兽场相关的Splits需要专用的存档才能正常使用。是否要生成专用存档？",
		Buttons:       []string{"Yes", "No"},
		DefaultButton: "Yes",
		CancelButton:  "No",
	}); err != nil {
		runtime.LogError(a.ctx, err.Error())
	} else if v == "Yes" {
		if err = os.WriteFile("user4.dat", blankColoSave, 0644); err != nil {
			_, _ = runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
				Type:    runtime.ErrorDialog,
				Title:   "错误",
				Message: err.Error(),
			})
		} else {
			_, _ = runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
				Type:    runtime.InfoDialog,
				Title:   "提示",
				Message: "已成功生成 user4.dat ，请自行放入空洞骑士的存档目录",
			})
		}
	}
}
