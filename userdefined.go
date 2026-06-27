package main

import (
	"embed"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed userdefined
var userDefined embed.FS

func (a *App) GetUserDefinedFiles() []Option {
	files, err := userDefined.ReadDir("userdefined")
	if err != nil {
		panic(err)
	}
	models := make([]Option, 0, len(files))
	for _, file := range files {
		fileName := file.Name()
		name := fileName[:len(fileName)-len(".lss")]
		models = append(models, Option{
			Value: name,
			Label: name,
		})
	}
	return models
}

func (a *App) OnSelectUserDefinedFile(fileName string) (*GetSplitsResult, error) {
	buf, err := userDefined.ReadFile("userdefined/" + fileName + ".lss")
	if err != nil {
		_, _ = runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Title:   "内部错误",
			Message: err.Error(),
		})
		return nil, err
	}
	ret, err := a.loadSplitFile(string(buf))
	if err != nil {
		return nil, err
	}
	for i := range ret.Splits {
		if i == 0 {
			continue
		}
		ret.Splits[i].Icon = getIconHtmlFormat(ret.Splits[i].Event)
	}
	return ret, nil
}
