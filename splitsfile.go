package main

import (
	"encoding/xml"
	"errors"
	"github.com/lxn/walk"
	"github.com/lxn/win"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type xmlIcon struct {
	Icon string `xml:",cdata"`
}

type xmlRun struct {
	XMLName              xml.Name `xml:"Run"`
	Version              string   `xml:"version,attr"`
	GameIcon             xmlIcon
	GameName             string
	CategoryName         string
	Metadata             xmlMetadata
	Offset               string
	AttemptCount         int
	AttemptHistory       string
	Segments             xmlSegments
	AutoSplitterSettings xmlAutoSplitterSettings
}

type xmlMetadata struct {
	Run       xmlMetadataRun
	Platform  xmlMetadataPlatform
	Variables xmlVariables
}

type xmlVariables struct {
	Variable []xmlVariable
}

type xmlVariable struct {
	Name  string `xml:"name,attr"`
	Value string `xml:",chardata"`
}

type xmlMetadataRun struct {
	Id string `xml:"id,attr"`
}

type xmlMetadataPlatform struct {
	UsesEmulator string `xml:"usesEmulator,attr"`
	Platform     string `xml:",chardata"`
}

type xmlSegments struct {
	Segment []*xmlSegment
}

type xmlSegment struct {
	Name            string
	Icon            xmlIcon
	SplitTimes      xmlSplitTimes
	BestSegmentTime xmlSplitTime
	SegmentHistory  xmlSegmentHistory
}

type xmlSegmentHistory struct {
	Time []xmlSplitTime
}

type xmlSplitTimes struct {
	SplitTime []xmlSplitTime
}

type xmlSplitTime struct {
	Id       string `xml:"id,attr,omitempty"`
	Name     string `xml:"name,attr,omitempty"`
	RealTime string `xml:"RealTime,omitempty"`
	GameTime string `xml:"GameTime,omitempty"`
}

type xmlAutoSplitterSettings struct {
	Ordered            string
	AutosplitEndRuns   string
	AutosplitStartRuns string
	Splits             xmlSplits
}

type xmlSplits struct {
	Split []string
}

func onClickLoadSplitFile() {
	dlg := new(walk.FileDialog)
	dlg.Title = "打开Splits文件"
	dlg.Filter = "Splits文件（*.lss）|*.lss"
	dlg.Flags = win.OFN_FILEMUSTEXIST
	if ok, err := dlg.ShowOpen(mainWindow); err != nil {
		walk.MsgBox(mainWindow, "内部错误", err.Error(), walk.MsgBoxIconError)
	} else if ok {
		loadSplitFile(dlg.FilePath)
	}
}

func loadSplitFile(file string) {
	if filepath.Ext(file) != ".lss" {
		return
	}
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		walk.MsgBox(mainWindow, "内部错误", err.Error(), walk.MsgBoxIconError)
		return
	}
	run := &xmlRun{}
	err = xml.Unmarshal(buf, run)
	if err != nil {
		walk.MsgBox(mainWindow, "内部错误", err.Error(), walk.MsgBoxIconError)
		return
	}
	count := len(run.AutoSplitterSettings.Splits.Split)
	if run.AutoSplitterSettings.AutosplitEndRuns != "True" {
		count++
	}
	if count <= 1 {
		walk.MsgBox(mainWindow, "错误", "暂不支持只有一个片段或者无片段的文件", walk.MsgBoxIconError)
		return
	}
	resetLines(count - 1)
	if startTrigger, ok := splitsDictIdToDescriptions[run.AutoSplitterSettings.AutosplitStartRuns]; ok {
		startTriggerCheckBox.SetChecked(true)
		err := startTriggerComboBox.SetText(startTrigger)
		if err != nil {
			walk.MsgBox(mainWindow, "错误", err.Error(), walk.MsgBoxIconError)
			return
		}
		startTriggerComboBox.SetEnabled(true)
	} else {
		startTriggerCheckBox.SetChecked(false)
		startTriggerComboBox.SetEnabled(false)
	}
	if run.AutoSplitterSettings.AutosplitEndRuns == "True" {
		for i, splitId := range run.AutoSplitterSettings.Splits.Split {
			if i < len(run.AutoSplitterSettings.Splits.Split)-1 {
				description := splitsDictIdToDescriptions[splitId]
				err := lines[i].splitId.SetText(description)
				if err != nil {
					walk.MsgBox(nil, "错误", err.Error(), walk.MsgBoxIconError)
					return
				}
				if i < len(run.Segments.Segment) {
					err = lines[i].name.SetText(run.Segments.Segment[i].Name)
					if err != nil {
						walk.MsgBox(nil, "错误", err.Error(), walk.MsgBoxIconError)
						return
					}
				}
			} else {
				description := splitsDictIdToDescriptions[splitId]
				finalLine.endTrigger.SetChecked(false)
				err := finalLine.splitId.SetText(description)
				if err != nil {
					walk.MsgBox(nil, "错误", err.Error(), walk.MsgBoxIconError)
					return
				}
				if i < len(run.Segments.Segment) {
					err = finalLine.name.SetText(run.Segments.Segment[i].Name)
					if err != nil {
						walk.MsgBox(nil, "错误", err.Error(), walk.MsgBoxIconError)
						return
					}
				}
			}
		}
	} else {
		for i, splitId := range run.AutoSplitterSettings.Splits.Split {
			description := splitsDictIdToDescriptions[splitId]
			err := lines[i].splitId.SetText(description)
			if err != nil {
				walk.MsgBox(nil, "错误", err.Error(), walk.MsgBoxIconError)
				return
			}
			if i < len(run.Segments.Segment) {
				err = lines[i].name.SetText(run.Segments.Segment[i].Name)
				if err != nil {
					walk.MsgBox(nil, "错误", err.Error(), walk.MsgBoxIconError)
					return
				}
			}
		}
		finalLine.endTrigger.SetChecked(true)
		i := len(run.AutoSplitterSettings.Splits.Split)
		if i < len(run.Segments.Segment) {
			text := "空洞骑士"
			name := run.Segments.Segment[i].Name
			if strings.Contains(name, "无上辐光") || strings.Contains(name, "Absolute Radiance") {
				text = "无上辐光"
			} else if strings.Contains(name, "辐光") || strings.Contains(name, "Radiance") || strings.Contains(name, "radiance") {
				text = "辐光"
			}
			err := finalLine.splitId2.SetText(text)
			if err != nil {
				walk.MsgBox(nil, "错误", err.Error(), walk.MsgBoxIconError)
				return
			}
			err = finalLine.name.SetText(run.Segments.Segment[i].Name)
			if err != nil {
				walk.MsgBox(nil, "错误", err.Error(), walk.MsgBoxIconError)
				return
			}
		}
	}
}

func saveSplitsFile() {
	if err := checkSplitsFile(); err != nil {
		walk.MsgBox(mainWindow, "内部错误", err.Error(), walk.MsgBoxIconError)
		return
	}
	dlg := new(walk.FileDialog)
	dlg.Title = "保存Splits文件"
	dlg.Filter = "Splits文件（*.lss）|*.lss"
	dlg.ShowReadOnlyCB = true
	dlg.Flags = win.OFN_OVERWRITEPROMPT | win.OFN_NOREADONLYRETURN
	if ok, err := dlg.ShowSave(mainWindow); err != nil {
		walk.MsgBox(mainWindow, "内部错误", err.Error(), walk.MsgBoxIconError)
		return
	} else if ok {
		filePath := dlg.FilePath
		if filepath.Ext(filePath) != ".lss" {
			filePath += ".lss"
		}
		run := &xmlRun{
			Version:  "1.7.0",
			GameName: "Hollow Knight",
			Metadata: xmlMetadata{Platform: xmlMetadataPlatform{UsesEmulator: "False"}},
			Offset:   "00:00:00",
			AutoSplitterSettings: xmlAutoSplitterSettings{
				Ordered:          "True",
				AutosplitEndRuns: "False",
				Splits:           xmlSplits{},
			},
		}
		for _, line := range lines {
			splitId := splitsDict[line.splitId.Text()].id
			run.Segments.Segment = append(run.Segments.Segment, &xmlSegment{
				Name:       line.name.Text(),
				Icon:       xmlIcon{getIcon(splitId)},
				SplitTimes: xmlSplitTimes{SplitTime: []xmlSplitTime{{Name: "Personal Best"}}},
			})
			run.AutoSplitterSettings.Splits.Split = append(run.AutoSplitterSettings.Splits.Split, splitId)
		}
		run.Segments.Segment = append(run.Segments.Segment, &xmlSegment{
			Name:       finalLine.name.Text(),
			SplitTimes: xmlSplitTimes{SplitTime: []xmlSplitTime{{Name: "Personal Best"}}},
		})
		if startTriggerCheckBox.Checked() {
			text := startTriggerComboBox.Text()
			run.AutoSplitterSettings.AutosplitStartRuns = splitsDict[text].id
		}
		if finalLine.endTrigger.Checked() {
			switch finalLine.splitId2.Text() {
			case "空洞骑士":
				run.Segments.Segment[len(run.Segments.Segment)-1].Icon.Icon = getIcon("HollowKnightBoss")
			case "辐光":
				fallthrough
			case "无上辐光":
				run.Segments.Segment[len(run.Segments.Segment)-1].Icon.Icon = getIcon("RadianceBoss")
			}
		} else {
			splitId := splitsDict[finalLine.splitId.Text()].id
			run.AutoSplitterSettings.AutosplitEndRuns = "True"
			run.Segments.Segment[len(run.Segments.Segment)-1].Icon.Icon = getIcon(splitId)
			run.AutoSplitterSettings.Splits.Split = append(run.AutoSplitterSettings.Splits.Split, splitId)
		}
		buf, err := xml.MarshalIndent(run, "", "  ")
		if err != nil {
			walk.MsgBox(mainWindow, "错误", err.Error(), walk.MsgBoxIconError)
			return
		}
		err = ioutil.WriteFile(filePath, append([]byte(xml.Header), buf...), 0644)
		if err != nil {
			walk.MsgBox(mainWindow, "错误", err.Error(), walk.MsgBoxIconError)
			return
		}
		walk.MsgBox(mainWindow, "成功", "保存成功", walk.MsgBoxIconInformation)
	}
}

func checkSplitsFile() error {
	if startTriggerCheckBox.Checked() {
		text := startTriggerComboBox.Text()
		if _, ok := splitsDict[text]; !ok {
			return errors.New(`自动分割配置"` + text + `"不存在，请检查`)
		}
	}
	for _, line := range lines {
		text := line.splitId.Text()
		if _, ok := splitsDict[text]; !ok {
			return errors.New(`自动分割配置"` + text + `"不存在，请检查`)
		}
	}
	if !finalLine.endTrigger.Checked() {
		text := finalLine.splitId.Text()
		if _, ok := splitsDict[text]; !ok {
			return errors.New(`自动分割配置"` + text + `"不存在，请检查`)
		}
	} else {
		text := finalLine.splitId2.Text()
		if text != "无上辐光" && text != "辐光" && text != "空洞骑士" {
			return errors.New(`最后一个片段设置有误`)
		}
	}
	return nil
}
