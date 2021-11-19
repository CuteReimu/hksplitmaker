package main

import (
	"encoding/xml"
	"errors"
	"github.com/lxn/walk"
	"github.com/lxn/win"
	"io/ioutil"
	"path/filepath"
)

type xmlRun struct {
	XMLName              xml.Name `xml:"Run"`
	Version              string   `xml:"version,attr"`
	GameIcon             string
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
	Variables string
}

type xmlMetadataRun struct {
	Id string `xml:"id,attr"`
}

type xmlMetadataPlatform struct {
	UsesEmulator string `xml:"usesEmulator,attr"`
}

type xmlSegments struct {
	Segment []*xmlSegment
}

type xmlSegment struct {
	Name            string
	Icon            string `xml:",innerxml"`
	SplitTimes      xmlSplitTimes
	BestSegmentTime string
	SegmentHistory  string
}

type xmlSplitTimes struct {
	SplitTime []xmlSplitTime
}

type xmlSplitTime struct {
	Name string `xml:"name,attr"`
}

type xmlAutoSplitterSettings struct {
	Ordered          string
	AutosplitEndRuns string
	Splits           xmlSplits
}

type xmlSplits struct {
	Split []string
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
				Icon:       getIcon(splitId),
				SplitTimes: xmlSplitTimes{SplitTime: []xmlSplitTime{{Name: "Personal Best"}}},
			})
			run.AutoSplitterSettings.Splits.Split = append(run.AutoSplitterSettings.Splits.Split, splitId)
		}
		run.Segments.Segment = append(run.Segments.Segment, &xmlSegment{
			Name:       finalLine.name.Text(),
			SplitTimes: xmlSplitTimes{SplitTime: []xmlSplitTime{{Name: "Personal Best"}}},
		})
		if finalLine.endTrigger.Checked() {
			switch finalLine.splitId2.Text() {
			case "空洞骑士":
				run.Segments.Segment[len(run.Segments.Segment)-1].Icon = getIcon("HollowKnightBoss")
			case "辐光":
				fallthrough
			case "无上辐光":
				run.Segments.Segment[len(run.Segments.Segment)-1].Icon = getIcon("RadianceBoss")
			}
		} else {
			splitId := splitsDict[finalLine.splitId.Text()].id
			run.AutoSplitterSettings.AutosplitEndRuns = "True"
			run.Segments.Segment[len(run.Segments.Segment)-1].Icon = getIcon(splitId)
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
