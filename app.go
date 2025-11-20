package main

import (
	"archive/zip"
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gopkg.in/toast.v1"
)

type App struct {
	ctx context.Context
}
type PackData struct {
	CleanName        string `json:"CleanName"`
	PackType         string `json:"PackType"`
	ResoucePackPath  string `json:"ResoucePackPath"`
	BehaviorPackPath string `json:"BehaviorPackPath"`
	ExportPath       string `json:"ExportPath"`
	Format           string `json:"Format"`
}
type ResourcePack struct {
	CleanName    string
	Path         string
	UUID         string
	IsSignatures bool
}
type BehaviorPack struct {
	CleanName        string
	ScriptState      string
	Path             string
	UUID             string
	DependenciesUUID []string
	IsSignatures     bool
}
type Addon struct {
	CleanName    string
	ResourcePack ResourcePack
	BehaviorPack BehaviorPack
}
type Manifest struct {
	Header struct {
		Name string `json:"name"`
		UUID string `json:"uuid"`
	} `json:"header"`
	Dependencies []map[string]interface{} `json:"dependencies,omitempty"`
}

func NewApp() *App {
	return &App{}
}
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}
func (a *App) SelfVersion() string {
	ico := "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"
	data, _ := base64.StdEncoding.DecodeString(ico)
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	path := filepath.Join(dir, "icon.png")
	_ = os.WriteFile(path, data, 0644)
	return "3.0.0"
}
func (a *App) GetUserSetting() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	path := filepath.Join(dir, "userSettings.json")

	if _, err := os.Stat(path); err != nil {
		content := `{"mode": "development","exportLocation": "desktop", "theme": "default", "format": "mcaddon"}`
		_ = os.WriteFile(path, []byte(content), 0644)
		return content
	}
	content, _ := os.ReadFile(path)
	return string(content)
}
func (a *App) SaveUserSetting(data string) {
	err := os.WriteFile(filepath.Join(filepath.Dir(os.Args[0]), "userSettings.json"), []byte(data), 0644)
	fmt.Println(err)
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
	}
*/
func (a *App) OpenDirectoryDialog() (string, error) {
	directoryPath, err := runtime.OpenDirectoryDialog(
		a.ctx,
		runtime.OpenDialogOptions{
			Title: "Select Export Directory",
		},
	)
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
		os.Exit(1)
	}
}
func (a *App) GetData() string {
	//C:\Users\An\AppData\Roaming\Minecraft Bedrock\Users\Shared\games\com.mojang\development_resource_packs
	var gamePath string = "C:\\Users\\" + os.Getenv("USERNAME") + "\\AppData\\Roaming\\Minecraft Bedrock\\Users\\Shared\\games\\com.mojang"
	var resourcePacksPath string = gamePath + "resource_packs"
	var behaviorPacksPath string = gamePath + "behavior_packs"
	var developmentResourcePacksPath string = gamePath + "\\development_resource_packs"
	var developmentBehaviorPacksPath string = gamePath + "\\development_behavior_packs"

	listResourcePacks := ListResoucePack(resourcePacksPath)
	listResourcePacks = append(listResourcePacks, ListResoucePack(developmentResourcePacksPath)...)

	listBehaviorPacks := ListBehaviorPack(behaviorPacksPath)
	listBehaviorPacks = append(listBehaviorPacks, ListBehaviorPack(developmentBehaviorPacksPath)...)

	exclusiveResourcePack, addonPack, exclusiveBehaviorPack := sortPacks(listResourcePacks, listBehaviorPacks)

	data := []any{exclusiveResourcePack, addonPack, exclusiveBehaviorPack}
	jsonData, _ := json.Marshal(data)
	return string(jsonData)
}
func (a *App) UpdateScriptVersion(behaviorPath string, oldVersion string) {
	content, _ := os.ReadFile(behaviorPath + "\\manifest.json")
	newContent := strings.ReplaceAll(string(content), oldVersion, "beta")
	_ = os.WriteFile(behaviorPath+"\\manifest.json", []byte(newContent), 0644)
}
func (a *App) CompilePack(packData string) {
	data := &PackData{}
	if err := json.Unmarshal([]byte(packData), data); err != nil {
		fmt.Println("Error parsing JSON data:", err)
	}
	fmt.Println("Parsed Pack Data:", data.PackType)
	var suffixPackType string = " Resource Pack." //behaviorPack addOnPack resourcePack
	if data.PackType == "addOnPack" {
		suffixPackType = " Addon."
	} else if data.PackType == "behaviorPack" {
		suffixPackType = " Behavior Pack."
	}
	var exportPath string = "C:\\Users\\" + os.Getenv("USERNAME") + "\\Desktop\\" + data.CleanName + suffixPackType + data.Format
	if data.ExportPath != "desktop" {
		exportPath = data.ExportPath + "\\" + data.CleanName + suffixPackType + data.Format
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
		Icon:  icon,
	}
	notification.Push()
}
func (a *App) NotifyText(title string) {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	path := filepath.Join(dir, "icon.png")
	notification := toast.Notification{
		AppID: "Add-On Compiler",
		Title: title,
		Icon:  path,
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
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
func ListResoucePack(dir string) []ResourcePack {
	dirHandle, _ := os.Open(dir)
	defer dirHandle.Close()
	entries, _ := dirHandle.Readdir(0)
	var resourcePackList []ResourcePack
	for _, entry := range entries {
		if entry.IsDir() {
			file, err := os.Open(dir + "\\" + entry.Name() + "\\" + "manifest.json")
			if err != nil {
				os.Exit(1)
			}
			defer file.Close()
			var manifest Manifest
			json.NewDecoder(file).Decode(&manifest)
			packName := entry.Name()
			if manifest.Header.Name == "pack.name" {
				// pack name is hidden in lang file
				// get lang file to get the name
				langFile, err := os.Open(dir + "\\" + entry.Name() + "\\" + "texts" + "\\" + "en_US.lang")
				if err != nil {
					const targetKey = "pack.name"
					// Use bufio.Scanner to read the file line by line
					scanner := bufio.NewScanner(langFile)
					for scanner.Scan() {
						line := scanner.Text()
						if strings.HasPrefix(line, targetKey+"=") {
							parts := strings.SplitN(line, "=", 2)
							if len(parts) == 2 {
								packName = parts[1]
								break
							}
						}
					}
				}
				defer langFile.Close()
			}
			resourcePackList = append(resourcePackList, ResourcePack{CleanName: packName, UUID: manifest.Header.UUID, Path: dir + "\\" + entry.Name(), IsSignatures: fileExists(dir + "\\" + entry.Name() + "\\signatures.json")})
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
			packName := entry.Name()
			if manifest.Header.Name == "pack.name" {
				// pack name is hidden in lang file
				// get lang file to get the name
				langFile, err := os.Open(dir + "\\" + entry.Name() + "\\" + "texts" + "\\" + "en_US.lang")
				if err != nil {
					const targetKey = "pack.name"
					// Use bufio.Scanner to read the file line by line
					scanner := bufio.NewScanner(langFile)
					for scanner.Scan() {
						line := scanner.Text()
						if strings.HasPrefix(line, targetKey+"=") {
							parts := strings.SplitN(line, "=", 2)
							if len(parts) == 2 {
								packName = parts[1]
								break
							}
						}
					}
				}
				defer langFile.Close()
			}
			var minecraftServerVersion string
			var dependenciesUUID []string
			for _, dep := range manifest.Dependencies {
				moduleName, ok := dep["module_name"].(string)
				if ok {
					if moduleName == "@minecraft/server" {
						version, ok := dep["version"].(string)
						if ok {
							minecraftServerVersion = version
						}
					}
				}
				depedencyUUID, ok := dep["uuid"].(string)
				if ok {
					dependenciesUUID = append(dependenciesUUID, depedencyUUID)
				}
			}
			if minecraftServerVersion == "" {
				minecraftServerVersion = "null-script"
			}
			behaviorPackList = append(behaviorPackList, BehaviorPack{CleanName: packName, DependenciesUUID: dependenciesUUID, ScriptState: minecraftServerVersion, Path: dir + "\\" + entry.Name(), IsSignatures: fileExists(dir + "\\" + entry.Name() + "\\signatures.json")})
		}
	}
	return behaviorPackList
}

func sortPacks(resourcePack []ResourcePack, behaviorPack []BehaviorPack) ([]ResourcePack, []Addon, []BehaviorPack) {
	// Create a map for O(1) lookup of resource packs by UUID
	rpMap := make(map[string]ResourcePack, len(resourcePack))
	for _, rp := range resourcePack {
		rpMap[rp.UUID] = rp
	}
	usedRP := make(map[string]bool)
	usedBP := make(map[string]bool)

	var addonPack []Addon
	for _, bp := range behaviorPack {
		matchedRP := []ResourcePack{}
		for _, depUUID := range bp.DependenciesUUID {
			if rp, exists := rpMap[depUUID]; exists {
				matchedRP = append(matchedRP, rp)
				usedRP[depUUID] = true
			}
		}
		if len(matchedRP) == 1 {
			addonPack = append(addonPack, Addon{
				CleanName:    cleanName(bp.CleanName),
				ResourcePack: matchedRP[0],
				BehaviorPack: bp,
			})
			usedBP[bp.UUID] = true
		}
	}
	var exclusiveResourcePack []ResourcePack
	for _, rp := range resourcePack {
		if !usedRP[rp.UUID] {
			exclusiveResourcePack = append(exclusiveResourcePack, rp)
		}
	}
	var exclusiveBehaviorPack []BehaviorPack
	for _, bp := range behaviorPack {
		if !usedBP[bp.UUID] {
			exclusiveBehaviorPack = append(exclusiveBehaviorPack, bp)
		}
	}
	return exclusiveResourcePack, addonPack, exclusiveBehaviorPack
}

// Normalize system
func ExtractPreTileString(input string) string {
	if nameIndex := strings.Index(input, ".name"); nameIndex != -1 {
		return input[5:nameIndex]
	}
	return ""
}
func ExtractPreEntityString(input string) string {
	if nameIndex := strings.Index(input, ".name"); nameIndex != -1 {
		return input[7:nameIndex]
	}
	return ""
}
func ExtractPostString(input string) string {
	if equalIndex := strings.Index(input, "="); equalIndex != -1 {
		value := strings.TrimSpace(strings.TrimSuffix(input[equalIndex+1:], "	#"))
		return value
	}
	return ""
}
func (a *App) EmitMessageToNormalizePanel(message string) {
	runtime.EventsEmit(a.ctx, "file:rename", message)
}

func (a *App) Normalize(rpPath string, bpPath string) string {
	file, _ := os.Open(rpPath + "\\texts\\en_US.lang")
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var entityArray [][]string
	var itemArray [][]string
	var tileArray [][]string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "entity.") {
			entityArray = append(entityArray, []string{ExtractPreEntityString(line), ExtractPostString(line)})
		} else if strings.HasPrefix(line, "item.") {
			itemArray = append(itemArray, []string{ExtractPreTileString(line), ExtractPostString(line)})
		} else if strings.HasPrefix(line, "tile.") {
			tileArray = append(tileArray, []string{ExtractPreTileString(line), ExtractPostString(line)})
		}
	}
	fmt.Println("Entity List: ", len(entityArray))
	fmt.Println(entityArray[0])
	fmt.Println("Item List: ", len(itemArray))
	fmt.Println(itemArray[0])
	fmt.Println("Tile List: ", len(tileArray))
	fmt.Println(tileArray[0])

	// Tile Section
	runtime.EventsEmit(a.ctx, "stage:name", "Normalizing Tiles")
	tilePath := bpPath + "\\blocks"
	filesStageTile, _ := os.ReadDir(tilePath)
	for _, file := range filesStageTile {
		fileName := file.Name()
		if !file.IsDir() && strings.HasSuffix(fileName, ".json") {
			filePath := filepath.Join(tilePath, fileName)
			fileData, _ := os.ReadFile(filePath)
			var blockInfo struct {
				MinecraftBlock struct {
					Components struct {
						MinecraftGeometry struct {
							Identifier string `json:"identifier"`
						} `json:"minecraft:geometry"`
					} `json:"components"`
					Description struct {
						Identifier string `json:"identifier"`
					} `json:"description"`
				} `json:"minecraft:block"`
			}
			json.Unmarshal(fileData, &blockInfo)
			for index, tileSubArray := range tileArray {
				if blockInfo.MinecraftBlock.Description.Identifier == tileSubArray[0] {
					if blockInfo.MinecraftBlock.Components.MinecraftGeometry.Identifier != "" && blockInfo.MinecraftBlock.Components.MinecraftGeometry.Identifier != "minecraft:geometry.full_block" {
						tileArray[index] = append(tileArray[index], blockInfo.MinecraftBlock.Components.MinecraftGeometry.Identifier)
					}
					a.EmitMessageToNormalizePanel("Renaming Tile: " + tileSubArray[0] + " to " + tileSubArray[1])
					newFilePath := filepath.Join(tilePath, tileSubArray[1]+".json")
					if err := os.Rename(filePath, newFilePath); err != nil {
						fmt.Println("Error renaming file:", err)
						continue
					}
				}
			}
		}
	}
	// Tile Model
	runtime.EventsEmit(a.ctx, "stage:name", "Normalizing Tiles Model (RP)")
	tileModelPath := rpPath + "\\models\\blocks"
	var jsonFiles []string
	filepath.Walk(tileModelPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(path, ".json") {
			jsonFiles = append(jsonFiles, path)
		}
		return nil
	})
	for _, filePath := range jsonFiles {
		a.EmitMessageToNormalizePanel("Reading Tile Model: " + filePath)
		fileData, _ := os.ReadFile(filePath)
		var blockInfo struct {
			MinecraftGeometry []struct {
				Description struct {
					Identifier string `json:"identifier"`
				} `json:"description"`
			} `json:"minecraft:geometry"`
		}
		json.Unmarshal(fileData, &blockInfo)
		for index, tileSubArray := range tileArray {
			if len(tileSubArray) == 2 {
				continue
			}
			if blockInfo.MinecraftGeometry[0].Description.Identifier == tileSubArray[2] {
				tileArray = append(tileArray[:index], tileArray[index+1:]...)
				newFilePath := filepath.Join(filepath.Dir(filePath), tileSubArray[1]+".json")
				if err := os.Rename(filePath, newFilePath); err != nil {
					fmt.Println("Error renaming file:", err)
					continue
				}
			}
		}
	}
	// Behavior Entity Section
	runtime.EventsEmit(a.ctx, "stage:name", "Normalizing Behavior Entity files")
	entityPath := bpPath + "\\entities"
	filesStageEntity, _ := os.ReadDir(entityPath)
	for _, file := range filesStageEntity {
		fileName := file.Name()
		if !file.IsDir() && strings.HasSuffix(fileName, ".json") {
			filePath := filepath.Join(entityPath, fileName)
			fileData, _ := os.ReadFile(filePath)

			var entityInfo struct {
				MinecraftEntity struct {
					Description struct {
						Identifier string `json:"identifier"`
					} `json:"description"`
				} `json:"minecraft:entity"`
			}
			json.Unmarshal(fileData, &entityInfo)
			for index, entitySubArray := range entityArray {
				if entityInfo.MinecraftEntity.Description.Identifier == entitySubArray[0] {
					if entityInfo.MinecraftEntity.Description.Identifier != "" {
						entityArray[index] = append(entityArray[index], entityInfo.MinecraftEntity.Description.Identifier)
					}
					a.EmitMessageToNormalizePanel("Renaming Tile: " + entitySubArray[0] + " to " + entitySubArray[1])
					newFilePath := filepath.Join(entityPath, entitySubArray[1]+".json")
					if err := os.Rename(filePath, newFilePath); err != nil {
						fmt.Println("Error renaming file:", err)
						continue
					}
				}
			}
		}
	}
	// Resource Entity Section
	runtime.EventsEmit(a.ctx, "stage:name", "Normalizing Resource Entity files")
	clientEntityPath := rpPath + "\\entity"
	filesStageClientEntity, _ := os.ReadDir(clientEntityPath)
	for _, file := range filesStageClientEntity {
		fileName := file.Name()
		if !file.IsDir() && strings.HasSuffix(fileName, ".json") {
			filePath := filepath.Join(clientEntityPath, fileName)
			fileData, _ := os.ReadFile(filePath)

			var clientEntityInfo struct {
				MinecraftEntity struct {
					Description struct {
						Identifier string `json:"identifier"`
					} `json:"description"`
				} `json:"minecraft:client_entity"`
			}
			json.Unmarshal(fileData, &clientEntityInfo)
			for _, clientEntitySubArray := range entityArray {
				if clientEntityInfo.MinecraftEntity.Description.Identifier == clientEntitySubArray[0] {
					a.EmitMessageToNormalizePanel("Renaming Tile: " + clientEntitySubArray[0] + " to " + clientEntitySubArray[1])
					newFilePath := filepath.Join(clientEntityPath, clientEntitySubArray[1]+".json")
					if err := os.Rename(filePath, newFilePath); err != nil {
						fmt.Println("Error renaming file:", err)
						continue
					}
				}
			}
		}
	}
	// Item Section
	runtime.EventsEmit(a.ctx, "stage:name", "Normalizing Items")
	itemPath := bpPath + "\\items"
	filesStageItem, _ := os.ReadDir(itemPath)
	for _, file := range filesStageItem {
		fileName := file.Name()
		if !file.IsDir() && strings.HasSuffix(fileName, ".json") {
			filePath := filepath.Join(itemPath, fileName)
			fileData, _ := os.ReadFile(filePath)
			var itemInfo struct {
				MinecraftItem struct {
					Description struct {
						Identifier string `json:"identifier"`
					} `json:"description"`
				} `json:"minecraft:item"`
			}
			json.Unmarshal(fileData, &itemInfo)
			for _, itemSubArray := range itemArray {
				if itemInfo.MinecraftItem.Description.Identifier == itemSubArray[0] {
					a.EmitMessageToNormalizePanel("Renaming Tile: " + itemSubArray[0] + " to " + itemSubArray[1])
					newFilePath := filepath.Join(itemPath, itemSubArray[1]+".json")
					if err := os.Rename(filePath, newFilePath); err != nil {
						fmt.Println("Error renaming file:", err)
						continue
					}
				}
			}
		}
	}
	return "Done"
}