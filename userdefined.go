package main

import (
	"embed"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

//go:embed userdefined
var userDefined embed.FS

var curUserDefined string

func GetUserDefinedComboBox() ComboBox {
	files, err := userDefined.ReadDir("userdefined")
	if err != nil {
		panic(err)
	}
	models := make([]string, 0, len(files))
	for _, file := range files {
		fileName := file.Name()
		models = append(models, fileName[:len(fileName)-len(".lss")])
	}
	var comboBox *walk.ComboBox
	return ComboBox{
		AssignTo: &comboBox,
		Model:    models,
		OnCurrentIndexChanged: func() {
			newUserDefined := comboBox.Text()
			if len(newUserDefined) == 0 || newUserDefined == curUserDefined {
				return
			}
			curUserDefined = newUserDefined
			buf, err := userDefined.ReadFile("userdefined/" + newUserDefined + ".lss")
			if err != nil {
				walk.MsgBox(mainWindow, "内部错误", err.Error(), walk.MsgBoxIconError)
				return
			}
			loadSplitFile(buf)
		},
	}
}
