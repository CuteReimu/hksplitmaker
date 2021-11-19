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
var splitsDictIdToDescriptions = make(map[string]string)
var splitsDict = make(map[string]*splitData)
var splitsSearchDict = make(map[string][]string)

func initSplitsSearchDict(content string) {
	rs := ([]rune)(content)
	for i := 0; i < len(rs); i++ {
		for j := i + 1; j <= len(rs); j++ {
			s := string(rs[i:j])
			v, ok := splitsSearchDict[s]
			if !ok {
				v = nil
			}
			splitsSearchDict[s] = append(v, content)
		}
	}
}

func initSplitsFile(init bool) error {
	var err error
	splitsFileBuf, err = ioutil.ReadFile("splits.txt")
	if err != nil && !init {
		return err
	}
	rd := bufio.NewReader(bytes.NewReader(splitsFileBuf))
	re, err := regexp.Compile(`\[Description\("(.*?)"\)\s*,\s*ToolTip\("(.*?)"\)]`)
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
				splitsDictIdToDescriptions[string(line)] = description
				initSplitsSearchDict(description)
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

type splitIdModel struct {
	walk.ListModelBase
	items []string
}

func (s *splitIdModel) Value(index int) interface{} {
	return s.items[index]
}

func (s *splitIdModel) ItemCount() int {
	return len(s.items)
}

type lineData struct {
	line    *walk.Composite
	name    *walk.LineEdit
	splitId *walk.ComboBox
}

type finalLineData struct {
	lineData
	splitId2   *walk.ComboBox
	endTrigger *walk.CheckBox
}

var lines []*lineData
var finalLine finalLineData

func addLine(initAll bool) {
	line := new(lineData)
	c := Composite{
		AssignTo: &line.line,
		Layout:   HBox{},
		MaxSize:  Size{Width: 0, Height: 25},
		Children: []Widget{
			LineEdit{AssignTo: &line.name, MinSize: Size{Width: 200}},
			ComboBox{AssignTo: &line.splitId, Editable: true, MinSize: Size{Width: 200},
				Model: &splitIdModel{}, Value: splitDescriptions[0],
				OnTextChanged: func() {
					onSearchSplitId(initAll, line.splitId)
				},
			},
			PushButton{Text: "✘", MaxSize: Size{Width: 25}, OnClicked: func() {
				if len(lines) > 1 {
					removeLine(line)
				}
			}},
			PushButton{Text: "↑+", MaxSize: Size{Width: 25},
				OnClicked: func() {
					idx := splitLinesView.Children().Index(line.line)
					addLine(true)
					moveLine(idx)
				},
			},
			PushButton{Text: "↓+", MaxSize: Size{Width: 25},
				OnClicked: func() {
					idx := splitLinesView.Children().Index(line.line)
					addLine(true)
					moveLine(idx + 1)
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

func removeLine(line *lineData) {
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
	lines = append(lines[:idx], lines[idx+1:]...)
}

func cleanAllLines() {
	err := splitLinesViewContainer.Children().RemoveAt(0)
	if err != nil {
		walk.MsgBox(mainWindow, "错误", err.Error(), walk.MsgBoxIconError)
		return
	}
	splitLinesView.Dispose()
	err = Composite{
		AssignTo:  &splitLinesView,
		Alignment: AlignHCenterVNear,
		Layout:    VBox{},
	}.Create(NewBuilder(splitLinesViewContainer))
	if err != nil {
		walk.MsgBox(mainWindow, "错误", err.Error(), walk.MsgBoxIconError)
		panic(err)
	}
	lines = nil
}

func moveLine(index int) {
	last := lines[len(lines)-1]
	if len(last.name.Text()) != 0 {
		walk.MsgBox(nil, "错误", "内部错误", walk.MsgBoxIconError)
		return
	}
	for i := len(lines) - 2; i >= index; i-- {
		err := lines[i+1].splitId.SetText(lines[i].splitId.Text())
		if err != nil {
			walk.MsgBox(nil, "错误", err.Error(), walk.MsgBoxIconError)
			return
		}
		err = lines[i+1].name.SetText(lines[i].name.Text())
		if err != nil {
			walk.MsgBox(nil, "错误", err.Error(), walk.MsgBoxIconError)
			return
		}
	}
	err := lines[index].name.SetText("")
	if err != nil {
		walk.MsgBox(nil, "错误", err.Error(), walk.MsgBoxIconError)
		return
	}
}

func onSearchSplitId(initAll bool, splitIdComboBox *walk.ComboBox) {
	s := splitIdComboBox.Text()
	model := splitIdComboBox.Model().(*splitIdModel)
	if len(model.items) == 0 {
		if initAll {
			model.items = splitDescriptions
		} else {
			model.items = []string{s}
		}
		model.PublishItemsReset()
		return
	}
	if len(s) > 0 {
		for _, text := range model.items {
			if text == s {
				return
			}
		}
		v, ok := splitsSearchDict[s]
		if ok && len(v) > 0 {
			if len(model.items) > 0 {
				model.PublishItemsRemoved(0, len(model.items)-1)
			}
			model.items = v
			if len(v) > 0 {
				model.PublishItemsInserted(0, len(v)-1)
			}
		}
	}
}
