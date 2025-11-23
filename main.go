package main

import (
	"embed"
	"encoding/base64"
	"net/http"
	"os"
	"strings"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

type FileLoader struct {
	http.Handler
}

func NewFileLoader() *FileLoader {
	return &FileLoader{}
}

func (h *FileLoader) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	encodedFilePath := req.URL.Path[1:] // Remove leading '/'

	// 2. Remove the extension (e.g., .png) if it exists
	if idx := strings.LastIndex(encodedFilePath, "."); idx != -1 {
		encodedFilePath = encodedFilePath[:idx]
	}
	decodedFilePath, err := base64.StdEncoding.DecodeString(encodedFilePath)
	if err != nil {
		http.Error(res, "Invalid path", 400)
		return
	}

	filePath := string(decodedFilePath)
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("Error getting file(s)"))
	}
	res.Write(fileData)
}

func main() {
	app := NewApp()
	err := wails.Run(&options.App{
		Title:  "Add-On Compiler",
		Width:  1850,
		Height: 800,
		AssetServer: &assetserver.Options{
			Assets:  assets,
			Handler: NewFileLoader(),
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})
	if err != nil {
		println("Error:", err.Error())
	}
}
