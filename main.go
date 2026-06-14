package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var frontEndAssets embed.FS

func main() {
	if err := initSplitsFile(); err != nil {
		panic(err)
	}

	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "空洞骑士计时器生成器",
		Width:  1000,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: frontEndAssets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []any{
			app,
		},
		Logger: appLogger,
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
