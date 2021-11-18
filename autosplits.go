package main

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"io"
	"io/ioutil"
	"regexp"
)

type splitData struct {
	description, tooltip, id string
}

var splitsFileBuf []byte
var splitDescriptions []string
var splitsDict = make(map[string]*splitData)

func initSplitsFile(init bool) error {
	var err error
	splitsFileBuf, err = ioutil.ReadFile("splits.txt")
	if err != nil && !init {
		return err
	}
	rd := bufio.NewReader(bytes.NewReader(splitsFileBuf))
	re, err := regexp.Compile("\\[Description\\(\"(.*?)\"\\)\\s*,\\s*ToolTip\\(\"(.*?)\"\\)]")
	if err != nil {
		return err
	}
	splitDescriptions = nil
	splitsDict = make(map[string]*splitData)
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
				splitDescriptions = append(splitDescriptions, description)
				splitsDict[description] = &splitData{description: description, tooltip: result[2], id: string(line)}
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
	if err == io.EOF {
		return nil
	}
	return err
}

type lineData struct {
	line    *walk.Composite
	name    *walk.LineEdit
	splitId *walk.ComboBox
}

var lines []*lineData
var finalLine lineData
var endTriggerCheckBox *walk.CheckBox

func addLine() {
	line := new(lineData)
	c := Composite{
		AssignTo: &line.line,
		Layout:   HBox{},
		MaxSize:  Size{Width: 0, Height: 25},
		Children: []Widget{
			LineEdit{AssignTo: &line.name, MinSize: Size{Width: 200}},
			ComboBox{AssignTo: &line.splitId, Editable: true, Model: splitDescriptions, MaxSize: Size{Width: 200}, Value: splitDescriptions[0]},
			PushButton{Text: "✘", MaxSize: Size{Width: 25}, OnClicked: func() {
				if len(lines) > 1 {
					idx := splitLinesView.Children().Index(line.line)
					if idx < 0 {
						walk.MsgBox(mainWindow, "错误", "无法删除这一行", walk.MsgBoxIconError)
						return
					}
					err := splitLinesView.Children().RemoveAt(idx)
					if err != nil {
						walk.MsgBox(mainWindow, "错误", err.Error(), walk.MsgBoxIconError)
						return
					}
					line.line.Dispose()
					lines = append(lines[:idx], lines[:idx+1]...)
				}
			}},
			PushButton{Text: "↑+", MaxSize: Size{Width: 25},
				OnClicked: func() {
					addLine()
				},
			},
			PushButton{Text: "↓+", MaxSize: Size{Width: 25},
				OnClicked: func() {
					addLine()
				},
			},
		},
	}
	err := c.Create(NewBuilder(splitLinesView))
	if err != nil {
		walk.MsgBox(nil, "错误", err.Error(), walk.MsgBoxIconError)
	}
	lines = append(lines, line)
}
