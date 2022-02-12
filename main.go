package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
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
	width, height := 550, 750
	err := MainWindow{
		OnDropFiles: func(f []string) {
			if len(f) > 0 {
				loadSplitFile(f[0])
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
					PushButton{Text: "打开已有的Splits文件", OnClicked: onClickLoadSplitFile},
					TextLabel{TextAlignment: AlignHFarVCenter, Text: "，也可以使用现有的模板"},
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
						Text:          "Auto Splitter Version: 3.1.2.0",
					},
					PushButton{
						MaxSize:   Size{Width: 100},
						Alignment: AlignHFarVCenter,
						Text:      "帮助",
						OnClicked: func() {
							walk.MsgBox(mainWindow, "帮助", readme, walk.MsgBoxIconInformation)
						},
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
								AssignTo: &startTriggerCheckBox,
								Text:     "自动开始",
								OnClicked: func() {
									startTriggerComboBox.SetEnabled(startTriggerCheckBox.Checked())
								},
							},
							ComboBox{
								AssignTo: &startTriggerComboBox,
								Model:    splitDescriptions,
								Enabled:  false,
								Editable: true,
								Value:    splitDescriptions[0],
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
							LineEdit{AssignTo: &finalLine.name, Text: "空洞骑士"},
							ComboBox{AssignTo: &finalLine.splitId, Visible: false, Editable: true, Model: splitDescriptions, MaxSize: Size{Width: 200}, Value: splitDescriptions[0]},
							ComboBox{AssignTo: &finalLine.splitId2, Editable: true, Model: []string{"空洞骑士", "辐光", "无上辐光"}, MaxSize: Size{Width: 200}, Value: "空洞骑士"},
							CheckBox{AssignTo: &finalLine.endTrigger, Checked: true, Text: "以游戏结束停止计时",
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
				},
			},
		},
	}.Create()
	addLine(true)
	if err != nil {
		walk.MsgBox(nil, "错误", err.Error(), walk.MsgBoxIconError)
		return
	}
	//go func() {
	//	resp, err := http.Get("https://raw.githubusercontent.com/slaurent22/hk-split-maker/main/src/asset/splits.txt")
	//	if err != nil || resp.StatusCode != 200 {
	//		return
	//	}
	//	buf, err := ioutil.ReadAll(resp.Body)
	//	if err != nil {
	//		return
	//	}
	//	if !bytes.Equal(buf, splitsFileBuf) {
	//		err := updateBtn.SetText("更新")
	//		if err != nil {
	//			walk.MsgBox(mainWindow, "错误", err.Error(), walk.MsgBoxIconError)
	//		} else {
	//			var a int
	//			a = updateBtn.Clicked().Attach(func() {
	//				updateBtn.SetEnabled(false)
	//				err := ioutil.WriteFile(filepath.Join(hkSplitMakerDir, "splits.txt"), buf, 0644)
	//				if err != nil {
	//					walk.MsgBox(mainWindow, "错误", err.Error(), walk.MsgBoxIconError)
	//					return
	//				}
	//				initSplitsFile(false)
	//				updateBtn.Clicked().Detach(a)
	//				err = updateBtn.SetText("已是最新")
	//				if err != nil {
	//					walk.MsgBox(mainWindow, "错误", err.Error(), walk.MsgBoxIconError)
	//					return
	//				}
	//			})
	//			updateBtn.SetEnabled(true)
	//		}
	//	}
	//}()
	hWnd := mainWindow.Handle()
	currStyle := win.GetWindowLong(hWnd, win.GWL_STYLE)
	win.SetWindowLong(hWnd, win.GWL_STYLE, currStyle & ^win.WS_SIZEBOX)
	mainWindow.Run()
}
