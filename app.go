package main

import (
	"context"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx         context.Context
	uploadedRun *xmlRun
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

type Option struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

func (a *App) GetOptions() []Option {
	options := make([]Option, 0, len(splitsDictIdToDescriptions))
	for id, desc := range splitsDictIdToDescriptions {
		options = append(options, Option{Value: id, Label: desc})
	}
	slices.SortStableFunc(options, func(a, b Option) int {
		i := strings.LastIndex(a.Label, "(")
		j := strings.LastIndex(b.Label, "(")
		if i == -1 {
			i = len(a.Label)
		}
		if j == -1 {
			j = len(b.Label)
		}
		ad, bd := a.Label[i:], b.Label[j:]
		d := strings.Compare(ad, bd)
		if d != 0 {
			return d
		}
		return strings.Compare(a.Label, b.Label)
	})
	return options
}

func (a *App) GetTemplates() []Option {
	return GetAllFiles()
}

type SplitLine struct {
	Name  string        `json:"name"`
	Event string        `json:"event"`
	Icon  string        `json:"icon"`
	Other []*xmlElement `json:"other"`
}

type GetSplitsResult struct {
	Name            string      `json:"name"`
	Splits          []SplitLine `json:"splits"`
	StartTriggering bool        `json:"startTriggering"`
	EndTriggering   bool        `json:"endTriggering"`
}

// LoadSplitFile parses an .lss file from its XML text content
func (a *App) LoadSplitFile(filePath string) (*GetSplitsResult, error) {
	if filepath.Ext(filePath) != ".lss" {
		return nil, errors.New("只支持*.lss文件")
	}
	buf, err := os.ReadFile(filePath)
	if err != nil {
		slog.Error("找不到文件", "filePath", filePath, "error", err)
		_, _ = runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Title:   "错误",
			Message: "找不到文件",
		})
		return nil, err
	}
	return a.loadSplitFile(string(buf))
}

func (a *App) loadSplitFile(content string) (*GetSplitsResult, error) {
	run := &xmlRun{}
	if err := xml.Unmarshal([]byte(content), run); err != nil {
		return nil, err
	}

	result := []SplitLine{{
		Name:  "开始",
		Event: strings.TrimSpace(run.AutoSplitterSettings.AutosplitStartRuns),
	}}
	if result[0].Event == "" {
		result[0].Event = "StartNewGame"
	}
	for i := range max(len(run.Segments), len(run.AutoSplitterSettings.Splits)) {
		result = append(result, SplitLine{})
		if i < len(run.Segments) {
			result[i+1].Name = run.Segments[i].Name
			result[i+1].Icon = convertIconToHtmlFormat(run.Segments[i].Icon.Icon)
			result[i+1].Other = run.Segments[i].Other
		}
		if i < len(run.AutoSplitterSettings.Splits) {
			result[i+1].Event = run.AutoSplitterSettings.Splits[i]
		}
	}
	if run.AutoSplitterSettings.AutosplitEndRuns != "True" {
		result[len(result)-1].Event = "EndingSplit"
	}

	a.uploadedRun = run
	return &GetSplitsResult{
		Name:            run.CategoryName,
		Splits:          result,
		StartTriggering: strings.TrimSpace(run.AutoSplitterSettings.AutosplitStartRuns) != "",
		EndTriggering:   run.AutoSplitterSettings.AutosplitEndRuns == "True",
	}, nil
}

// GetSplits returns split lines from a template file
func (a *App) GetSplits(name string) (*GetSplitsResult, error) {
	c, err := GetSplitMakerConfig(name)
	if err != nil {
		return nil, err
	}

	splits := make([]SplitLine, 0, len(c.Ids)+2)
	splits = append(splits, SplitLine{
		Name:  "开始",
		Icon:  getIconHtmlFormat("StartNewGame"),
		Event: "StartNewGame",
	})
	if len(c.StartTriggering) > 0 {
		splits[0].Event = c.StartTriggering
	}
	nameCache := make(map[string]int, len(c.Names))
	iconCache := make(map[string]int, len(c.Icons))
	re := regexp.MustCompile(`\{.*?}`)
	for _, id := range c.Ids {
		isSubSplit := strings.HasPrefix(id, "-")
		id = strings.TrimPrefix(id, "-")
		id = re.ReplaceAllString(id, "")
		var splitName, icon string
		splitName = splitsDictIdToDescriptions[id]
		if idx := strings.LastIndex(splitName, "("); idx > 0 {
			splitName = strings.TrimSpace(splitName[:idx])
		}
		if nameCache[id] < len(c.Names[id]) {
			splitName = translate(strings.ReplaceAll(c.Names[id][nameCache[id]], "%s", splitName))
			nameCache[id]++
		}
		if iconCache[id] < len(c.Icons[id]) {
			icon = getIconHtmlFormat(c.Icons[id][iconCache[id]])
			iconCache[id]++
		} else {
			icon = getIconHtmlFormat(id)
		}
		if isSubSplit {
			splitName = "-" + splitName
		}
		splits = append(splits, SplitLine{Name: splitName, Event: id, Icon: icon})
	}

	if !c.EndTriggering {
		splits = append(splits, SplitLine{
			Name:  "游戏结束",
			Event: "EndingSplit",
			Icon:  getIconHtmlFormat("EndingSplit"),
		})
		if c.EndingSplit != nil {
			if c.EndingSplit.Name != "Completion" {
				splits[len(splits)-1].Name = translate(c.EndingSplit.Name)
			}
			splits[len(splits)-1].Icon = getIconHtmlFormat(c.EndingSplit.Icon)
			if _, ok := splitsDictIdToDescriptions[c.EndingSplit.Icon]; ok {
				splits[len(splits)-1].Event = c.EndingSplit.Icon
			}
		}
	}

	a.uploadedRun = &xmlRun{
		Version:      "1.7.0",
		GameName:     "Hollow Knight",
		CategoryName: c.CategoryName,
		Offset:       "00:00:00",
	}

	return &GetSplitsResult{
		Name:            c.CategoryName,
		Splits:          splits,
		StartTriggering: len(c.StartTriggering) > 0,
		EndTriggering:   c.EndTriggering,
	}, nil
}

// GetIcon returns the icon in HTML img-src format for a split ID
func (a *App) GetIcon(splitId string) string {
	return getIconHtmlFormat(splitId)
}

// buildSplits builds the LSS XML and returns base64-encoded bytes
func (a *App) buildSplits(data []SplitLine, includeTimeRecords, startTriggering, endTriggering bool) (string, error) {
	var fileRunData xmlRun
	if a.uploadedRun != nil {
		fileRunData = *a.uploadedRun
	}
	if !includeTimeRecords {
		fileRunData.Other = nil
	}
	if fileRunData.Version == "" {
		fileRunData.Version = "1.7.0"
	}
	if fileRunData.GameName == "" {
		fileRunData.GameName = "Hollow Knight"
	}
	fileRunData.Offset = "00:00:00"
	splits := make([]string, 0, len(data))
	fileRunData.Segments = make([]*xmlSegment, 0, len(data))
	for i, line := range data {
		if i == 0 {
			continue
		}
		if i < len(data)-1 || endTriggering {
			splits = append(splits, line.Event)
		}
		var other []*xmlElement
		if includeTimeRecords {
			other = line.Other
		}
		fileRunData.Segments = append(fileRunData.Segments, &xmlSegment{
			Name:  line.Name,
			Other: other,
			Icon:  xmlIcon{convertIconToLiveSplitFormat(line.Icon)},
		})
	}
	fileRunData.AutoSplitterSettings = autoSplittingRuntimeSettings{
		Ordered:          "True",
		AutosplitEndRuns: "False",
		Splits:           splits,
	}
	if startTriggering {
		fileRunData.AutoSplitterSettings.AutosplitStartRuns = data[0].Event
	}
	if endTriggering {
		fileRunData.AutoSplitterSettings.AutosplitEndRuns = "True"
	}

	buf, err := xml.MarshalIndent(fileRunData, "", "  ")
	if err != nil {
		return "", err
	}
	buf = append([]byte(`<?xml version="1.0" encoding="UTF-8"?>`+"\n"), buf...)
	return base64.StdEncoding.EncodeToString(buf), nil
}

// SaveSplitsFile shows a native save dialog and writes the LSS file
func (a *App) SaveSplitsFile(data []SplitLine, includeTimeRecords, startTriggering, endTriggering bool) error {
	b64, err := a.buildSplits(data, includeTimeRecords, startTriggering, endTriggering)
	if err != nil {
		return err
	}
	raw, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return err
	}
	dest, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		DefaultFilename: "splits.lss",
		Filters: []runtime.FileFilter{
			{DisplayName: "LiveSplit splits (*.lss)", Pattern: "*.lss"},
		},
	})
	if err != nil || dest == "" {
		return err
	}
	return os.WriteFile(dest, raw, 0644)
}

// SaveIconsZip shows a native save dialog and writes the icons zip
func (a *App) SaveIconsZip() error {
	raw, err := zipIcons()
	if err != nil {
		return err
	}
	dest, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		DefaultFilename: "icons.zip",
		Filters: []runtime.FileFilter{
			{DisplayName: "ZIP archive (*.zip)", Pattern: "*.zip"},
		},
	})
	if err != nil || dest == "" {
		return err
	}
	return os.WriteFile(dest, raw, 0644)
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	runtime.EventsOn(ctx, "onSelectColo", a.coloNotice)
	a.setWindowSize()
}

func (a *App) setWindowSize() {
	// 1. 获取所有屏幕的信息
	screens, _ := runtime.ScreenGetAll(a.ctx)
	if len(screens) == 0 {
		return // 获取失败，使用默认配置
	}

	// 2. 找到主屏幕 (IsPrimary) 或当前屏幕 (IsCurrent)
	var targetScreen *runtime.Screen
	for _, screen := range screens {
		if screen.IsPrimary || screen.IsCurrent {
			targetScreen = &screen
			if screen.IsCurrent {
				break
			}
		}
	}
	// 如果没找到主屏幕，就回退使用第一个屏幕
	if targetScreen == nil && len(screens) > 0 {
		targetScreen = &screens[0]
	}
	if targetScreen == nil {
		return
	}

	// 3. 计算高度：屏幕总高度减去 100 像素
	newHeight := targetScreen.Size.Height - 100

	// 防止负高度（比如屏幕高度本身就小于 100，虽然很少见）
	if newHeight < 200 {
		newHeight = targetScreen.Size.Height
	}

	runtime.WindowSetSize(a.ctx, 1000, newHeight)
	runtime.WindowCenter(a.ctx)
}

// --- XML types ---

type autoSplittingRuntimeSettings struct {
	Ordered            string
	AutosplitEndRuns   string
	AutosplitStartRuns string
	Splits             []string `xml:"Splits>Split"`
}

type xmlRun struct {
	XMLName              xml.Name      `xml:"Run"`
	Version              string        `xml:"version,attr"`
	GameIcon             string        `xml:"GameIcon"`
	GameName             string        `xml:"GameName"`
	CategoryName         string        `xml:"CategoryName"`
	Offset               string        `xml:"Offset"`
	AttemptCount         int           `xml:"AttemptCount"`
	Segments             []*xmlSegment `xml:"Segments>Segment"`
	AutoSplitterSettings autoSplittingRuntimeSettings
	Other                []*xmlElement `xml:",any"`
}

type xmlSegment struct {
	Name  string
	Icon  xmlIcon
	Other []*xmlElement `xml:",any"`
}

type xmlElement struct {
	XMLName xml.Name
	Content string `xml:",innerxml"`
}

type xmlIcon struct {
	Icon string `xml:",cdata"`
}
