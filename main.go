package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
	"syscall"
)

func getSystemMetrics(nIndex int) int {
	ret, _, _ := syscall.NewLazyDLL(`User32.dll`).NewProc(`GetSystemMetrics`).Call(uintptr(nIndex))
	return int(ret)
}

var mainWindow *walk.MainWindow
var splitLinesView *walk.Composite

func main() {
	if err := initSplitsFile(true); err != nil {
		walk.MsgBox(nil, "错误", "内部错误", walk.MsgBoxIconError)
		return
	}
	screenX, screenY := getSystemMetrics(0), getSystemMetrics(1)
	width, height := 550, 480
	err := MainWindow{
		AssignTo: &mainWindow,
		Title:    "hksplitmaker",
		Bounds:   Rectangle{X: (screenX - width) / 2, Y: (screenY - height) / 2, Width: width, Height: height},
		Layout:   VBox{},
		Children: []Widget{
			//Composite{
			//	MaxSize: Size{Width: 0, Height: 20},
			//	Layout:  HBox{},
			//	Children: []Widget{
			//		Composite{
			//			Layout: HBox{},
			//			Children: []Widget{
			//				TextLabel{Text: "你可以创建新的Splits文件，也可以"},
			//				PushButton{Text: "导入已有的文件"},
			//			},
			//		},
			//		HSeparator{},
			//		PushButton{AssignTo: &updateBtn, Text: "已是最新", Enabled: false},
			//	},
			//},
			ScrollView{
				HorizontalFixed: true,
				Layout:          VBox{},
				Children: []Widget{
					Composite{
						AssignTo:  &splitLinesView,
						Alignment: AlignHCenterVNear,
						Layout:    VBox{},
					},
					Composite{
						Layout: HBox{},
						Children: []Widget{
							LineEdit{AssignTo: &finalLine.name, Text: "空洞骑士"},
							ComboBox{AssignTo: &finalLine.splitId, Visible: false, Editable: true, Model: splitDescriptions, MaxSize: Size{Width: 200}, Value: splitDescriptions[0]},
							ComboBox{AssignTo: &finalLine.splitId2, Editable: true, Model: []string{"空洞骑士", "辐光", "无上辐光"}, MaxSize: Size{Width: 200}, Value: "空洞骑士"},
							CheckBox{AssignTo: &endTriggerCheckBox, Checked: true, Text: "以游戏结束停止计时",
								OnCheckedChanged: func() {
									finalLine.splitId.SetVisible(!endTriggerCheckBox.Checked())
									finalLine.splitId2.SetVisible(endTriggerCheckBox.Checked())
								},
							},
						},
					},
				},
			},
			PushButton{Text: "另存为", OnClicked: saveSplitsFile},
		},
	}.Create()
	addLine()
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
	//				err := ioutil.WriteFile("splits.txt", buf, 0644)
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
