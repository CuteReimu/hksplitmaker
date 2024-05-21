package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
	"os"
	"path/filepath"
	"sort"
	"syscall"
)

const hkSplitMakerDir = "hk-split-maker/src/asset/hollowknight/"

func getSystemMetrics(nIndex int) int {
	ret, _, _ := syscall.NewLazyDLL(`User32.dll`).NewProc(`GetSystemMetrics`).Call(uintptr(nIndex))
	return int(ret)
}

var mainWindow *walk.MainWindow
var splitLinesViewContainer *walk.Composite
var splitLinesView *walk.Composite
var categoriesComboBox *walk.ComboBox
var startTriggerCheckBox *walk.CheckBox
var startTriggerComboBox *walk.ComboBox
var saveTimeCheckBox *walk.CheckBox

func main() {
	initCategories()
	if err := initSplitsFile(true); err != nil {
		walk.MsgBox(nil, "错误", "内部错误", walk.MsgBoxIconError)
		return
	}
	screenX, screenY := getSystemMetrics(0), getSystemMetrics(1)
	width, height := 720, 960
	err := MainWindow{
		OnDropFiles: func(f []string) {
			if len(f) > 0 {
				file := f[0]
				if filepath.Ext(file) != ".lss" {
					return
				}
				buf, err := os.ReadFile(file)
				if err != nil {
					walk.MsgBox(mainWindow, "内部错误", err.Error(), walk.MsgBoxIconError)
					return
				}
				loadSplitFile(buf)
			}
		},
		AssignTo: &mainWindow,
		Title:    "计时器生成器",
		Bounds:   Rectangle{X: (screenX - width) / 2, Y: (screenY - height) / 2, Width: width, Height: height},
		Layout:   VBox{},
		Children: []Widget{
			Composite{
				MaxSize: Size{Height: 20},
				Layout:  HBox{},
				Children: []Widget{
					TextLabel{TextAlignment: AlignHFarVCenter, Text: "你可以"},
					PushButton{Text: "打开已有的lss文件", OnClicked: onClickLoadSplitFile},
					TextLabel{TextAlignment: AlignHFarVCenter, Text: "或者把文件拖拽进来，也可以使用现有模板"},
					ComboBox{
						AssignTo: &categoriesComboBox,
						Model: func() []string {
							var keys []string
							for key := range categoriesCache {
								keys = append(keys, key)
							}
							sort.Strings(keys)
							return keys
						}(),
						OnCurrentIndexChanged: onSelectCategory,
					},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					TextLabel{
						TextAlignment: AlignHFarVCenter,
						Text:          "还有一些其他玩家友情提供的",
					},
					GetUserDefinedComboBox(),
					TextLabel{
						TextAlignment: AlignHFarVCenter,
						Text:          "Auto Splitter Version: 3.1.13.0",
					},
					PushButton{
						Text:      "更新LiveSplit",
						OnClicked: fixLiveSplit,
					},
				},
			},
			ScrollView{
				HorizontalFixed: true,
				Layout:          VBox{},
				Children: []Widget{
					Composite{
						MaxSize: Size{Width: 0, Height: 25},
						Layout:  HBox{},
						Children: []Widget{
							CheckBox{
								AssignTo:    &startTriggerCheckBox,
								Text:        "自动开始",
								ToolTipText: "对于全关的速通和万神殿某一门的速通，不要勾选",
								OnClicked: func() {
									startTriggerComboBox.SetEnabled(startTriggerCheckBox.Checked())
								},
							},
							ComboBox{
								AssignTo: &startTriggerComboBox,
								Model:    &splitIdModel{},
								Enabled:  false,
								Editable: true,
								Value:    splitDescriptions[0],
								OnTextChanged: func() {
									onSearchSplitId(true, &lineData{splitId: startTriggerComboBox})
								},
							},
						},
					},
					Composite{
						AssignTo: &splitLinesViewContainer,
						Layout:   Flow{},
						Children: []Widget{
							Composite{
								AssignTo:  &splitLinesView,
								Alignment: AlignHCenterVNear,
								Layout:    VBox{},
							},
						},
					},
					Composite{
						Layout: HBox{},
						Children: []Widget{
							LineEdit{AssignTo: &finalLine.name, Text: "空洞骑士", ToolTipText: "片段名"},
							ComboBox{AssignTo: &finalLine.splitId, Visible: false, Editable: true, Model: splitDescriptions, MaxSize: Size{Width: 200}, Value: splitDescriptions[0]},
							ComboBox{AssignTo: &finalLine.splitId2, Editable: true, Model: []string{"空洞骑士", "辐光", "无上辐光"}, MaxSize: Size{Width: 200}, Value: "空洞骑士"},
							CheckBox{AssignTo: &finalLine.endTrigger, Checked: true, Text: "以游戏结束停止计时", ToolTipText: "如果是以游戏结束或者万神殿某一门结束停止计时，不要勾选",
								OnCheckedChanged: func() {
									finalLine.splitId.SetVisible(!finalLine.endTrigger.Checked())
									finalLine.splitId2.SetVisible(finalLine.endTrigger.Checked())
								},
							},
						},
					},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					CheckBox{AssignTo: &saveTimeCheckBox, Text: "保留原lss文件中的时间数据", Enabled: false},
					PushButton{Text: "另存为", OnClicked: saveSplitsFile},
					PushButton{
						MaxSize:   Size{Width: 100},
						Alignment: AlignHFarVCenter,
						Text:      "帮助",
						OnClicked: func() {
							walk.MsgBox(mainWindow, "帮助", readme, walk.MsgBoxIconInformation)
						},
					}, PushButton{
						MaxSize:   Size{Width: 100},
						Alignment: AlignHFarVCenter,
						Text:      "FAQ",
						OnClicked: func() {
							walk.MsgBox(mainWindow, "FAQ", faq, walk.MsgBoxIconInformation)
						},
					},
				},
			},
		},
	}.Create()
	addLine(true)
	if err != nil {
		walk.MsgBox(nil, "错误", err.Error(), walk.MsgBoxIconError)
		return
	}
	hWnd := mainWindow.Handle()
	currStyle := win.GetWindowLong(hWnd, win.GWL_STYLE)
	win.SetWindowLong(hWnd, win.GWL_STYLE, currStyle & ^win.WS_SIZEBOX)
	mainWindow.Run()
}
