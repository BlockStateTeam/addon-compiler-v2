package main

import (
	"os"
	"io"
	"fmt"
	"strings"
	"context"
	"io/ioutil"
	"archive/zip"
	"path/filepath"
    "encoding/json"
	"encoding/base64"
	"gopkg.in/toast.v1"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx context.Context
}
type PackData struct {
	CleanName		 string `json:"CleanName"`
	PackType         string `json:"PackType"`
	ResoucePackPath  string `json:"ResoucePackPath"`
	BehaviorPackPath string `json:"BehaviorPackPath"`
	ExportPath       string `json:"ExportPath"`
	Format           string `json:"Format"`
}
type ResoucePack struct {
    Name        string
    CleanName   string
    Path        string
}
type BehaviorPack struct {
    Name        string
    CleanName   string
	ScriptState string
    Path        string
}
type Addon struct {
	CleanName     string
	ScriptState   string
	ResourcePath  string
	BehaviorPath  string
}
type Manifest struct {
	Dependencies 	[]map[string]interface{} `json:"dependencies,omitempty"`
}
func NewApp() *App {
	return &App{}
}
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}
func (a *App) SelfVersion() (string) {
	return "2.0.0"
}
func (a *App) GetUserSetting() (string) {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	path := filepath.Join(dir, "userSettings.json")

	if _, err := os.Stat(path); err != nil {
		content := `{"mode": "development","exportLocation": "desktop", "theme": "default", "format": "mcaddon"}`
		_ = ioutil.WriteFile(path, []byte(content), 0644)
		return content
	}
	content, _ := ioutil.ReadFile(path)
	return string(content)
}
func (a *App) SaveUserSetting(data string) {
	err := ioutil.WriteFile(filepath.Join(filepath.Dir(os.Args[0]), "userSettings.json"), []byte(data), 0644);
	fmt.Println(err);
}
/*
func (a *App) OpenFileDialog(dialogOptions interface{}) (string, error) {
	directoryPath, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		DefaultDirectory: fmt.Sprintf("%v", dialogOptions.(map[string]interface{})["defaultDirectory"]),
		Title:            fmt.Sprintf("%v", dialogOptions.(map[string]interface{})["title"]),
		Filters: []runtime.FileFilter{
            {
                DisplayName: "Pack Files",
                Pattern:     "*.zip; *.mcpack; *.mcaddon",
            },
        },
	})
	if err != nil {
		return "", fmt.Errorf("failed opening dialog - %s", err.Error())
	}
	fmt.Println("Selected directory:", directoryPath)
	return directoryPath, nil
}*/
func (a *App) OpenDirectoryDialog(dialogOptions interface{}) (string, error) {
	directoryPath, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		DefaultDirectory: fmt.Sprintf("%v", dialogOptions.(map[string]interface{})["defaultDirectory"]),
		Title:            fmt.Sprintf("%v", dialogOptions.(map[string]interface{})["title"]),
	})
	if err != nil {
		return "", fmt.Errorf("failed opening dialog - %s", err.Error())
	}
	fmt.Println("Selected directory:", directoryPath)
	return directoryPath, nil
}
func renameFile(oldPath string, newPath string) {
	err := os.Rename(oldPath, newPath)
	if err != nil {
		fmt.Println(err)
	}
};
func (a *App) GetData(isDevMode bool) (string) {
	var gamePath string = "C:\\Users\\" + os.Getenv("USERNAME") + "\\AppData\\Local\\Packages\\Microsoft.MinecraftUWP_8wekyb3d8bbwe\\LocalState\\games\\com.mojang";
	var suffix string = "\\development_"
	if !isDevMode {suffix = "\\"};
    var resourcePacksPath string = gamePath + suffix + "resource_packs";
	var behaviorPacksPath string = gamePath + suffix + "behavior_packs";
	exclusiveResourcePack, addonPack, exclusiveBehaviorPack := sortPacks(ListResoucePack(resourcePacksPath), ListBehaviorPack(behaviorPacksPath));
	var duplicateName bool = false;
	for i := 0; i < len(addonPack); i++ {
		var addonPackIndiv = addonPack[i];
		if filepath.Base(addonPackIndiv.ResourcePath) == filepath.Base(addonPackIndiv.BehaviorPath) {
			duplicateName = true;
			renameFile(addonPackIndiv.ResourcePath, addonPackIndiv.ResourcePath + " RP");
			renameFile(addonPackIndiv.BehaviorPath, addonPackIndiv.BehaviorPath + " BP");
		}
	}
	if duplicateName {
		runtime.WindowReload(a.ctx);
	}
	data := []interface{}{exclusiveResourcePack, addonPack, exclusiveBehaviorPack}
    jsonData, _ := json.Marshal(data)
    return string(jsonData)
};
func (a *App) UpdateScriptVersion(behaviorPath string, oldVersion string, newVersion string) {
    content, _ := ioutil.ReadFile(behaviorPath + "\\manifest.json")
    newContent := strings.ReplaceAll(string(content), oldVersion, newVersion)
    _ = ioutil.WriteFile(behaviorPath + "\\manifest.json", []byte(newContent), 0644)
}
func (a *App) CompilePack(packData string) () {
	data := &PackData{}
	if err := json.Unmarshal([]byte(packData), data); err != nil {
		fmt.Println("Error parsing JSON data:", err)
	}
	fmt.Println("Parsed Pack Data:", data)
	var exportPath string = "C:\\Users\\" + os.Getenv("USERNAME") + "\\Desktop\\" + data.CleanName + "." + data.Format;
	if data.ExportPath != "desktop" {
		exportPath = data.ExportPath + "\\" + data.CleanName + "." + data.Format;
	}
	fmt.Println("Export Path:", exportPath)
	paths := []string{data.ResoucePackPath, data.BehaviorPackPath}
	if err := ZipFolders(paths, exportPath); err != nil {
		panic(err)
	}
}
func (a *App) Notify(title string, icon string) {
	notification := toast.Notification{
        AppID: "Add-On Compiler",
        Title: title, 
		Icon: icon,
    }
    notification.Push()
}
func (a *App) NotifyText(title string) {
	notification := toast.Notification{
        AppID: "Add-On Compiler",
        Title: title,
    }
    notification.Push()
}
func (a *App) GetImage(path string) string {
	imagePath := path + "//pack_icon.png"
	file, err := os.ReadFile(imagePath)
	if err != nil {
		return "data:image/gif;base64,R0lGODlhAQABAIAAAAUEBAAAACwAAAAAAQABAAACAkQBADs="
	}
	encodedImage := base64.StdEncoding.EncodeToString(file)
	return "data:image/png;base64," + encodedImage
}
func ZipFolders(folderPaths []string, exportPath string) error {
	zipFile, err := os.Create(exportPath)
	if err != nil {
		return err
	}
	defer func() {
		zipFile.Close()
	}()

	zipWriter := zip.NewWriter(zipFile)
	defer func() {
		zipWriter.Close()
	}()

	for _, folderPath := range folderPaths {
		err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			relativePath, err := filepath.Rel(folderPath, path)
			if err != nil {
				return err
			}

			zipEntry, err := zipWriter.Create(filepath.Join(filepath.Base(folderPath), relativePath))
			if err != nil {
				return err
			}

			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(zipEntry, file)
			return err
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func cleanName(name string) string {
    stringSuffixes := []string{"-RP", "RP", "-rp", "rp", "resource", "Resource", "resource pack", "-resource-pack", "resource-pack", "_resource_pack", "resource_pack",
							   "-BP", "BP", "-bp", "bp", "behavior", "Behavior", "behavior pack", "-behavior-pack", "behavior-pack", "_behavior_pack", "behavior_pack"}
    for _, suffix := range stringSuffixes {
        name = strings.ReplaceAll(name, suffix, "")
    }
    return strings.TrimSpace(name)
}
func ListResoucePack(dir string) []ResoucePack {
    dirHandle, _ := os.Open(dir)
    defer dirHandle.Close()
    entries, _ := dirHandle.Readdir(0)
	var resourcePackList []ResoucePack
	for _, entry := range entries {
		if entry.IsDir() {
			resourcePackList = append(resourcePackList, ResoucePack{Name: entry.Name(), CleanName: cleanName(entry.Name()), Path: dir + "\\" + entry.Name()})
		}
	}
	return resourcePackList
}
func ListBehaviorPack(dir string) []BehaviorPack {
    dirHandle, _ := os.Open(dir)
    defer dirHandle.Close()
    entries, _ := dirHandle.Readdir(0)
	var behaviorPackList []BehaviorPack
	for _, entry := range entries {
		if entry.IsDir() {
			file, _ := os.Open(dir + "\\" + entry.Name() + "\\" + "manifest.json")
			defer file.Close()
			var manifest Manifest
			json.NewDecoder(file).Decode(&manifest)
			var minecraftServerVersion string
			for _, dep := range manifest.Dependencies {
				moduleName, ok := dep["module_name"].(string)
				if !ok {
					continue
				}
				if moduleName == "@minecraft/server" {
					version, ok := dep["version"].(string)
					if ok {
						minecraftServerVersion = version
						break
					}
				}
			}
			if minecraftServerVersion == "" {
				minecraftServerVersion = "null-script"
			}
			behaviorPackList = append(behaviorPackList, BehaviorPack{Name: entry.Name(), CleanName: cleanName(entry.Name()), ScriptState: minecraftServerVersion, Path: dir + "\\" + entry.Name()})
		}
	}
	return behaviorPackList
}

func sortPacks(resourcePack []ResoucePack, behaviorPack []BehaviorPack ) ([]ResoucePack, []Addon, []BehaviorPack) {
	var exclusiveResourcePack []ResoucePack;
	var addonPack []Addon;
	var exclusiveBehaviorPack []BehaviorPack;
	behaviorPackMap := make(map[string]BehaviorPack)
	for _, bPack := range behaviorPack {
		behaviorPackMap[bPack.CleanName] = bPack
	}
	for _, rPack := range resourcePack {
		if bPack, ok := behaviorPackMap[rPack.CleanName]; ok {
			addon := Addon{
				CleanName:     rPack.CleanName,
				ScriptState:   bPack.ScriptState,
				ResourcePath:  rPack.Path,
				BehaviorPath:  bPack.Path,
			}
			addonPack = append(addonPack, addon)
		} else {
			exclusiveResourcePack = append(exclusiveResourcePack, rPack)
		}
	}
	for _, bPack := range behaviorPack {
		found := false
		for _, addon := range addonPack {
			if addon.CleanName == bPack.CleanName {
				found = true
				break
			}
		}
		if !found {
			exclusiveBehaviorPack = append(exclusiveBehaviorPack, bPack)
		}
	}
	return exclusiveResourcePack, addonPack, exclusiveBehaviorPack
}