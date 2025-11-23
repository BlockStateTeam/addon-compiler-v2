package main

import (
	"archive/zip"
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"syscall"
	"unsafe"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gopkg.in/toast.v1"
)

type App struct {
	ctx context.Context
}
type ZipAsset struct {
	SourcePath string
	ZipPath    string
}
type PackData struct {
	CleanName  string `json:"CleanName"`
	PackType   string `json:"PackType"`
	PackPath   string `json:"PackPath"`
	ExportPath string `json:"ExportPath"`
	Format     string `json:"Format"`
	RequiredRP string `json:"RequiredRP"`
	RequiredBP string `json:"RequiredBP"`
}
type ResourcePack struct {
	CleanName    string
	Icon         string
	Path         string
	UUID         string
	IsSignatures bool
}
type BehaviorPack struct {
	CleanName        string
	Icon             string
	ScriptState      string
	Path             string
	UUID             string
	DependenciesUUID []string
	IsSignatures     bool
}
type RequiredPacks struct {
	CleanName string
	Path      string
	Icon      string
}
type WorldData struct {
	CleanName  string
	Path       string
	Icon       string
	RequiredRP string
	RequiredBP string
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
type Nameable interface {
	GetCleanName() string
}

func SortByCleanName[T Nameable](items []T) {
	slices.SortFunc(items, func(a, b T) int {
		nameA := strings.ToLower(a.GetCleanName())
		nameB := strings.ToLower(b.GetCleanName())
		return strings.Compare(nameA, nameB)
	})
}

type GenericPack interface {
	GetCleanName() string
	GetPath() string
	GetUUID() string
}

// Implement the interface for all 4 types (this enables the Generic function to work)
func (r ResourcePack) GetCleanName() string { return r.CleanName }
func (b BehaviorPack) GetCleanName() string { return b.CleanName }
func (a Addon) GetCleanName() string        { return a.CleanName }
func (w WorldData) GetCleanName() string    { return w.CleanName }

// Also implement GetPath for ResourcePack and BehaviorPack
func (r ResourcePack) GetPath() string { return r.Path }
func (b BehaviorPack) GetPath() string { return b.Path }
func (r ResourcePack) GetUUID() string { return r.UUID }
func (b BehaviorPack) GetUUID() string { return b.UUID }

type WarningType string

const (
	WarningConflict      WarningType = "DEPENDENCY_CONFLICT" // Multiple BPs want same RP
	WarningMismatch      WarningType = "NAME_MISMATCH"       // UUID match but names differ
	WarningDuplicateUUID WarningType = "DUPLICATE_UUID"      // Two packs have the same UUID
)

type PackWarning struct {
	Type    WarningType
	Message string

	// For Dependency Conflicts (BP -> RP)
	ResourcePack ResourcePack
	InvolvedBPs  []BehaviorPack
	WinnerBP     BehaviorPack

	// For Duplicate UUIDs (RP vs RP or BP vs BP)
	ConflictingRPs []ResourcePack
	ConflictingBPs []BehaviorPack
}

func NewApp() *App {
	return &App{}
}
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}
func (a *App) SelfVersion() string {
	ico := "iVBORw0KGgoAAAANSUhEUgAAAQAAAAEACAYAAABccqhmAAAAAXNSR0IArs4c6QAAAARnQU1BAACxjwv8YQUAAAAJcEhZcwAADsMAAA7DAcdvqGQAAAZ/SURBVHhe7d2/ilVXGMbhSZwhIAgJjk4UFIKFY5CkCEGLkNIitzAXMbWtV5ApcwFzCyksrRRJJzoWElDwzzhiQFDCKEkOfBew9oK92OR9HhC/zrPPlh+rOB/rs7W1tX/++wME+rz+BgIJAAQTAAgmABBMACCYAEAwAYBgAgDBBACCCQAEEwAIJgAQTAAgmABAMAGAYAIAwQQAggkABBMACCYAEEwAIJgAQDABgGACAMEEAIINuxlo7+xeTcuy+cVmTe22t8/X1O74wXFNbZ6sv66p3dHfRzXN69qJKzXNZ9Tzv3n3vqZ5nT51sqY2u4e7Nc3LCQCCCQAEEwAIJgAQTAAgmABAMAGAYAIAwQQAggkABBMACDZsF2D/wn5N8+n5LXjPLsClj2dqms/G1Y2a2h0cPK+pXc931rMLMPV5pu5OrPTsDzw+elpTu6m/6+9hFwCYnQBAMAGAYAIAwQQAggkABBMACCYAEEwAIJgAQDABgGACAMGGLQPdPXe7pjZLXuwYcZnIkpdhLm9erKndiAWqUZeJ9Lz/qXae7dQ0LycACCYAEEwAIJgAQDABgGACAMEEAIIJAAQTAAgmABBMACCYAECwxd4M1LOk8fDtnzW121rfqqldzwLR1Nt0Rt0MdOfl/ZraffvVNzW1G7FA02PEktLK1EUly0DA7AQAggkABBMACCYAEEwAIJgAQDABgGACAMEEAIIJAAQTAAg2bBlo7+xeTW2WujyystTbZHo+15t372tq17MMNdWozzV1SavXvU+Pamqze7hb07ycACCYAEAwAYBgAgDBBACCCQAEEwAIJgAQTAAgmABAMAGAYAIAwf5Xy0BLveVl5fHR05ra9Cy2jFqg6Xk3UxeVRvwbK6OWzqZ+NstAwOwEAIIJAAQTAAgmABBMACCYAEAwAYBgAgDBBACCCQAEG7YLsH9hv6b59OwCbFzdqKnd8YPjmtpNvRiix6j9iZ5nmfrZRu11jHgvK1P3NG59uFnTvJwAIJgAQDABgGACAMEEAIIJAAQTAAgmABBMACCYAEAwAYBgAgDBFrsMNOqSh1ELRAcHz2tqM/UikZVfvvyhpnY9l5yMWCC6duJKTe1GLXaN+M6uv7hR07ycACCYAEAwAYBgAgDBBACCCQAEEwAIJgAQTAAgmABAMAGAYAIAwYYtA+2d3aupzdSbVFZOnzpZU7ufvv+upnYjFkjuvLxfU7ufv/6xpnY9iz09yzBTl7tGLXaNMvU723m2U9O8nAAgmABAMAGAYAIAwQQAggkABBMACCYAEEwAIJgAQDABgGACAMEWuwz08O2fNbXbWt+qqd3lzYs1zWvqMkzPMlTPzUCjbtP5/a8/amrTs9jVs0A0ytT3v3u4W9O8nAAgmABAMAGAYAIAwQQAggkABBMACCYAEEwAIJgAQDABgGACAMEWuww0arFj6pJGr6nLPT1LStvb52tq17PY0+Pep0c1tVny+x+xqHbrw82a5uUEAMEEAIIJAAQTAAgmABBMACCYAEAwAYBgAgDBBACCCQAEG7YLsH9hv6Y2lz6eqWleU3+jvjLid+qjnr/nYpCDg+c1zWfU8z9Zf11Tu8dHT2tq9+rjq5ra/Hb8a03zcgKAYAIAwQQAggkABBMACCYAEEwAIJgAQDABgGACAMEEAIIJAAQbtgx099ztmtr0LKn0XHLRswzUY+oC0ahLPnqWYXou05j6/D3LQEv+PzP1YhgXgwCzEwAIJgAQTAAgmABAMAGAYAIAwQQAggkABBMACCYAEEwAIFj8zUCjlmGmGnH7UK87L+/X1G5rfaumNqdPnayp3bUTV2panqkLRLuHuzXNywkAggkABBMACCYAEEwAIJgAQDABgGACAMEEAIIJAAQTAAgmABBs2DLQ3tm9mtqMWoYZccvNyojlpp7Fph6Pj57W1O7y5sWa2iz5vfR8z1P/nesvbtQ0LycACCYAEEwAIJgAQDABgGACAMEEAIIJAAQTAAgmABBMACCYAECwxS4DLVnP0sn29vma2hw/OK6p3VJvOVoZsdw16ll6biCa+m52nu3UNC8nAAgmABBMACCYAEAwAYBgAgDBBACCCQAEEwAIJgAQTAAg2LBdAGB5nAAgmABAMAGAYAIAwQQAggkABBMACCYAEEwAIJgAQDABgGACAMEEAIIJAAQTAAgmABBMACCYAEAwAYBgAgDBBACCCQAEEwAIJgAQTAAg1trav956WuxJBDjbAAAAAElFTkSuQmCC"
	data, _ := base64.StdEncoding.DecodeString(ico)
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	path := filepath.Join(dir, "icon.png")
	os.WriteFile(path, data, 0644)
	return "3.0.0"
}
func (a *App) GetUserSetting() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	path := filepath.Join(dir, "userSettings.json")

	if _, err := os.Stat(path); err != nil {
		content := `{"exportLocation": "desktop", "theme": "default", "format": "mcaddon"}`
		os.WriteFile(path, []byte(content), 0644)
		return content
	}
	content, _ := os.ReadFile(path)
	return string(content)
}
func (a *App) UninstallApp() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	path := filepath.Join(dir, "uninstall.exe")

	// 1. Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println("Uninstaller not found")
		return "Something went wrong! Uninstaller not found"
	}

	// 2. Prepare Windows API call
	// We use ShellExecute instead of exec.Command to support "runas" (Admin prompt)
	shell32 := syscall.NewLazyDLL("shell32.dll")
	procShellExecute := shell32.NewProc("ShellExecuteW")

	// Convert Go strings to Windows UTF-16 pointers
	ptrOperation, _ := syscall.UTF16PtrFromString("runas") // Triggers UAC
	ptrFile, _ := syscall.UTF16PtrFromString(path)
	ptrDir, _ := syscall.UTF16PtrFromString(dir)

	// Window visibility: 1 = SW_SHOWNORMAL, 0 = SW_HIDE
	var showCmd int32 = 1

	// 3. Call the API
	// ShellExecuteW(hwnd, lpOperation, lpFile, lpParameters, lpDirectory, nShowCmd)
	ret, _, _ := procShellExecute.Call(
		0, // hwnd (Parent window, 0 for none)
		uintptr(unsafe.Pointer(ptrOperation)),
		uintptr(unsafe.Pointer(ptrFile)),
		0, // Parameters (0 or nil if none)
		uintptr(unsafe.Pointer(ptrDir)),
		uintptr(showCmd),
	)

	// ShellExecute returns a value > 32 on success
	if ret <= 32 {
		fmt.Printf("Failed to launch uninstaller. Error code: %d\n", ret)
		return fmt.Sprintf("Something went wrong! Failed to launch uninstaller: %v", ret)
	}

	// 4. CRITICAL: Exit immediately
	// You must close this app so the uninstaller can delete the .exe file
	os.Exit(0)
	return ""
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
func (a *App) OpenDirectoryDialogCall() (string, error) {
	directoryPath, err := runtime.OpenDirectoryDialog(
		a.ctx,
		runtime.OpenDialogOptions{
			Title: "Select Folder", // Give it a title
		},
	)

	// Handle the cancellation error specifically if needed
	if err != nil {
		return "", err
	}
	return directoryPath, nil
}
func (a *App) GetData() string {
	//C:\Users\An\AppData\Roaming\Minecraft Bedrock\Users\Shared\games\com.mojang\development_resource_packs
	var gamePath string = "C:\\Users\\" + os.Getenv("USERNAME") + "\\AppData\\Roaming\\Minecraft Bedrock\\Users\\Shared\\games\\com.mojang\\"
	listOfResourcePacks := ListResourcePack(gamePath + "resource_packs")
	listOfResourcePacks = append(listOfResourcePacks, ListResourcePack(gamePath+"development_resource_packs")...)
	listOfBehaviorPacks := ListBehaviorPack(gamePath + "behavior_packs")
	listOfBehaviorPacks = append(listOfBehaviorPacks, ListBehaviorPack(gamePath+"development_behavior_packs")...)
	worldData := ListWorldData(listOfBehaviorPacks, listOfResourcePacks)
	// fmt.Println("Total Resource Packs: ", len(listOfResourcePacks))
	// fmt.Println("Total Behavior Packs: ", len(listOfBehaviorPacks))
	// fmt.Println("List Resource Packs: ")
	// index := 1
	// for _, rp := range listOfResourcePacks {
	// 	fmt.Println(index, " - ", rp.CleanName, " | UUDI: ", rp.UUID, " | Signatures:", rp.IsSignatures)
	// 	index++
	// }
	// fmt.Println("List Behavior Packs: ")
	// index = 1
	// for _, rp := range listOfBehaviorPacks {
	// 	fmt.Println(index, " - ", rp.CleanName, " | UUDI: ", rp.UUID, " | Signatures:", rp.IsSignatures)
	// 	index++
	// }
	exclusiveResourcePack, addonPack, exclusiveBehaviorPack, warning := sortPacks(listOfResourcePacks, listOfBehaviorPacks)
	// fmt.Println("List of Exclusive Resource Packs: ")
	// for _, rp := range exclusiveResourcePack {
	// 	fmt.Println(" - ", rp.CleanName, " | UUDI: ", rp.UUID, " | Signatures:", rp.IsSignatures)
	// }
	// fmt.Println("List of Addon Packs: ")
	// for _, ap := range addonPack {
	// 	fmt.Println(" - ", ap.CleanName, " | RP UUDI: ", ap.ResourcePack.UUID, " | BP UUDI:", ap.BehaviorPack.UUID)
	// }
	// fmt.Println("List of Exclusive Behavior Packs: ")
	// for _, bp := range exclusiveBehaviorPack {
	// 	fmt.Println(" - ", bp.CleanName, " | UUDI: ", bp.UUID, " | Signatures:", bp.IsSignatures)
	// }
	// for _, warn := range warning {
	// 	fmt.Println(warn.Type, ": ", warn.Message)
	// }
	SortByCleanName(exclusiveResourcePack)
	SortByCleanName(addonPack)
	SortByCleanName(exclusiveBehaviorPack)
	SortByCleanName(worldData)
	jsonData, _ := json.Marshal([]interface{}{
		exclusiveResourcePack,
		addonPack,
		exclusiveBehaviorPack,
		worldData,
		warning,
	})
	return string(jsonData)
}
func (a *App) UpdateScriptVersion(behaviorPathEncode string, oldVersion string) string {
	fmt.Print("behaviorPathEncode: ", behaviorPathEncode)
	behaviorPathDecoded, err := base64.StdEncoding.DecodeString(behaviorPathEncode)
	if err != nil {
		log.Println("Error decoding behavior pack path:", err)
		return "Error: Failed to decode behavior pack path"
	}
	var paths []string
	err = json.Unmarshal(behaviorPathDecoded, &paths)
	if err != nil {
		return "Error: Failed to parse pack's path"
	}
	var path string
	switch len(paths) {
	case 0:
		return "Error: Something went wrong"
	case 1: // Behavior Pack
		path = paths[0]
	case 2: // Addon Pack
		path = paths[1]
	}
	path = path + "\\manifest.json"
	content, err := os.ReadFile(path)
	if err != nil {
		log.Println("Error reading manifest.json:", err)
		return "Error: Failed to read manifest.json"
	}
	newContent := strings.ReplaceAll(string(content), oldVersion, "beta")
	_ = os.WriteFile(path, []byte(newContent), 0644)
	return "Success: Updated Script API version"
}
func (a *App) CompilePack(packData string) string {
	log.Println("!!!!Compiling Pack with Data:", packData)
	data := &PackData{}
	if err := json.Unmarshal([]byte(packData), data); err != nil {
		fmt.Println("Error parsing JSON data:", err)
	}
	log.Println("Parsed Pack Data:", data.PackType)
	var suffixPackType string // behaviorPack addOnPack resourcePack
	format := data.Format
	switch data.PackType {
	case "addOnPack":
		suffixPackType = " Addon."
	case "behaviorPack":
		suffixPackType = " Behavior Pack."
		format = "mcpack"
	case "resourcePack":
		suffixPackType = " Resource Pack."
		format = "mcpack"
	case "world":
		suffixPackType = " World."
		format = "mcworld"
	}
	var exportPath string = "C:\\Users\\" + os.Getenv("USERNAME") + "\\Desktop\\" + cleanName(data.CleanName) + suffixPackType + format
	if data.ExportPath != "desktop" {
		exportPath = data.ExportPath + "\\" + data.CleanName + suffixPackType + format
	}
	log.Println("Export Path:", exportPath)
	decodedPackPathArrayString, err := base64.StdEncoding.DecodeString(data.PackPath)
	if err != nil {
		return "Error: Failed to decoded pack's path"
	}
	log.Println("!!!Decoded Pack Path Array String:", string(decodedPackPathArrayString))
	var paths []string
	// decodedPackPathArrayString is a string of JSON array
	// We need to unmarshal it and assign to paths
	err = json.Unmarshal(decodedPackPathArrayString, &paths)
	if err != nil {
		return "Error: Failed to parse pack's path array"
	}
	if len(paths) == 0 {
		return "Error: No valid pack path found"
	}
	assets := []ZipAsset{}
	switch data.PackType {
	case "addOnPack":
		assets = append(assets, ZipAsset{
			SourcePath: paths[0],
			ZipPath:    filepath.Base(paths[0]),
		})

		assets = append(assets, ZipAsset{
			SourcePath: paths[1],
			ZipPath:    filepath.Base(paths[1]),
		})
	case "world":
		assets = append(assets, ZipAsset{
			SourcePath: paths[0],
			ZipPath:    "",
		})
		// Handle behaviorPack and resourcePack if available
		encodedRequiredRP := data.RequiredRP
		fmt.Println("!!!!encodedRequiredRP: ", encodedRequiredRP)
		encodedRequiredBP := data.RequiredBP
		fmt.Println("!!!!encodedRequiredBP: ", encodedRequiredBP)
		// Decode base64
		decodedRP, err := base64.StdEncoding.DecodeString(encodedRequiredRP)
		if err == nil { // No error. Success
			var requiredRPs []RequiredPacks
			if err := json.Unmarshal(decodedRP, &requiredRPs); err == nil {
				for _, rp := range requiredRPs {
					assets = append(assets, ZipAsset{
						SourcePath: rp.Path,
						ZipPath:    filepath.Join("resource_packs", filepath.Base(rp.Path)),
					})
				}
			}
		}
		decodedBP, err := base64.StdEncoding.DecodeString(encodedRequiredBP)
		if err == nil {
			var requiredBPs []RequiredPacks
			if err := json.Unmarshal(decodedBP, &requiredBPs); err == nil {
				for _, bp := range requiredBPs {
					assets = append(assets, ZipAsset{
						SourcePath: bp.Path,
						ZipPath:    filepath.Join("behavior_packs", filepath.Base(bp.Path)),
					})
				}
			}
		}
	default: // behaviorPack or resourcePack
		assets = append(assets, ZipAsset{
			SourcePath: paths[0],
			ZipPath:    "",
		})
	}
	if err := ZipAssets(assets, exportPath); err != nil {
		return "Error: Failed to create file"
	}
	return fmt.Sprintf("Success: Export pack %s", data.CleanName)
}
func (a *App) Notify(title string, icon string) {
	// concat icon path + "\\
	var iconPath string
	iconPath = filepath.Join(icon, "pack_icon.png")
	// Check if iconPath exists
	if _, err := os.Stat(iconPath); os.IsNotExist(err) {
		// Its not exist, use other name
		iconPath = filepath.Join(icon, "world_icon.jpeg")
		if _, err := os.Stat(iconPath); os.IsNotExist(err) {
			// Its not exist, use default icon
			iconPath = filepath.Join(filepath.Dir(os.Args[0]), "icon.png")
		}
	}
	notification := toast.Notification{
		AppID: "Add-On Compiler",
		Title: title,
		Icon:  iconPath,
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
	err := notification.Push()
	if err != nil {
		log.Println("Error showing notification:", err)
	}
}
func addFileToZip(zipWriter *zip.Writer, sourcePath, zipPath string) error {
	// 1. Open the source file
	fileToZip, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	// 2. Clean the path for Zip specification (Always forward slashes)
	zipPath = filepath.ToSlash(zipPath)

	// 3. Create the header manually so we can set Compression
	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	// IMPORTANT: Set the name to our custom path, not the original filename
	header.Name = zipPath
	// Enable Compression (Deflate is the standard zip compression)
	header.Method = zip.Deflate

	// 4. Create the writer for this specific file
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	// 5. Copy content
	_, err = io.Copy(writer, fileToZip)
	return err
}
func ZipAssets(assets []ZipAsset, exportPath string) error {
	// 1. Create the output file
	zipFile, err := os.Create(exportPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	// 2. Initialize the Zip Writer
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// 3. Process every asset in the list
	for _, asset := range assets {
		// Stat the file to ensure it exists and check if it's a folder
		info, err := os.Stat(asset.SourcePath)
		if err != nil {
			return fmt.Errorf("could not find source %s: %w", asset.SourcePath, err)
		}

		// If it's a single file, just add it
		if !info.IsDir() {
			if err := addFileToZip(zipWriter, asset.SourcePath, asset.ZipPath); err != nil {
				return err
			}
			continue
		}

		// If it's a directory, walk it recursively
		// We use the 'asset.SourcePath' as the base to calculate relative paths
		baseDir := asset.SourcePath
		err = filepath.Walk(asset.SourcePath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil // We don't need to create explicit directory entries, files create them implicitly
			}

			// relativePath becomes the path *inside* the source folder
			// e.g. if walking "myFolder", and we find "myFolder/sub/file.txt", relative is "sub/file.txt"
			relativePath, err := filepath.Rel(baseDir, path)
			if err != nil {
				return err
			}

			// Join the relative path with the user's desired ZipPath
			// e.g. User wants "custom/location/" + "sub/file.txt"
			finalZipPath := filepath.Join(asset.ZipPath, relativePath)

			return addFileToZip(zipWriter, path, finalZipPath)
		})

		if err != nil {
			return fmt.Errorf("error walking directory %s: %w", asset.SourcePath, err)
		}
	}

	return nil
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

// resolvePackName handles the logic of checking "pack.name" vs the .lang file
func resolvePackName(packPath, entryName, manifestName string) string {
	if manifestName != "pack.name" {
		return entryName // Or manifestName if you prefer the raw name
	}

	// Use filepath.Join for cross-platform safety
	langFilePath := filepath.Join(packPath, "texts", "en_US.lang")
	langFile, err := os.Open(langFilePath)
	if err != nil {
		return entryName // Fallback to directory name if lang file fails
	}
	defer langFile.Close()

	scanner := bufio.NewScanner(langFile)
	const targetKey = "pack.name"

	for scanner.Scan() {
		line := scanner.Text()
		// TrimSpace is safer to handle accidental leading spaces
		if strings.HasPrefix(strings.TrimSpace(line), targetKey+"=") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				return strings.TrimSpace(parts[1])
			}
		}
	}
	return entryName // Fallback
}

// scanPacks iterates a directory and processes valid Minecraft packs
func scanPacks(rootDir string, processor func(path string, name string, manifest Manifest, icon string)) {
	entries, err := os.ReadDir(rootDir) // os.ReadDir is newer and faster than os.Open + Readdir
	if err != nil {
		log.Println("ERROR: Failed to read directory:", rootDir, err)
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		packPath := filepath.Join(rootDir, entry.Name())
		//log.Println("Scanning: ", packPath)
		manifestPath := filepath.Join(packPath, "manifest.json")

		file, err := os.Open(manifestPath)
		if err != nil {
			continue // Not a pack or permission denied
		}

		var manifest Manifest
		err = json.NewDecoder(file).Decode(&manifest)
		file.Close() // Close immediately!

		if err != nil {
			log.Println("WARNING: Invalid JSON in:", packPath)
			continue
		}

		// Resolve the name using our helper
		cleanName := resolvePackName(packPath, entry.Name(), manifest.Header.Name)
		icon := filepath.Join(packPath, "pack_icon.png")

		// Hand off the data to the specific logic
		processor(packPath, cleanName, manifest, icon)
	}
}
func getDependency[T GenericPack](filePath string, list []T) string {
	var returnData []RequiredPacks
	file, err := os.Open(filePath)
	if err != nil { // file doesn't exist
		return "W10=" // Return empty list [] on error Note that base64 of [] is W10=
	}
	defer file.Close()
	type PackManifest struct {
		PackID string `json:"pack_id"`
	}
	var entries []PackManifest // Using the struct from Tip #2
	if err := json.NewDecoder(file).Decode(&entries); err == nil {
		for _, entry := range entries {
			// The entry is UUID, use it as lookup by listOfBehaviorPacks and listOfResourcePacks
			for _, pack := range list {
				if pack.GetUUID() == entry.PackID {
					returnData = append(returnData, RequiredPacks{
						CleanName: pack.GetCleanName(),
						Path:      pack.GetPath(),
						Icon:      pack.GetPath() + "\\pack_icon.png",
					})
					continue
				}
			}
		}
	}
	jsonData, err := json.Marshal(returnData)
	if err != nil {
		return "W10=" // Return empty list [] on error Note that base64 of [] is W10=
	}
	if string(jsonData) == "null" { // JSON marshal can return null on empty slice which is invalid for our case
		return "W10=" // Return empty list [] on error Note that base64 of [] is W10=
	}
	// convert jsonData to base64
	encoded := base64.StdEncoding.EncodeToString(jsonData)
	return encoded
}
func ListWorldData(listOfBehaviorPacks []BehaviorPack, listOfResourcePacks []ResourcePack) []WorldData {
	//C:\Users\An\AppData\Roaming\Minecraft Bedrock\Users\63352952236107600\games\com.mojang\minecraftWorlds\R3MByq0SLb0=
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// handle error
	}
	// filepath.Join handles slashes automatically
	worldsDir := filepath.Join(homeDir, "AppData", "Roaming", "Minecraft Bedrock", "Users")
	// Here's we first encounter the first problem
	// Users have different numeric IDs based on Xbox Live accounts
	// We need to handle multiple user folders assume that each valid user folder name is numeric only and ignore others
	// Then we scan each user's minecraftWorlds folder then normalize the data
	// Finally we combine all worlds into a single list
	var allWorlds []WorldData
	userEntries, err := os.ReadDir(worldsDir)
	if err != nil {
		log.Println("ERROR: Failed to read worlds root directory:", worldsDir, err)
		return allWorlds
	}
	for _, userEntry := range userEntries {
		if !userEntry.IsDir() {
			log.Println("Skipping non-numeric user folder:", userEntry.Name())
			continue
		}
		// Check if the folder name is numeric only
		if _, err := strconv.ParseInt(userEntry.Name(), 10, 64); err != nil {
			log.Println("Skipping non-numeric user folder:", userEntry.Name())
			continue
		}
		userWorldsDir := filepath.Join(worldsDir, userEntry.Name(), "games", "com.mojang", "minecraftWorlds")
		worldEntries, err := os.ReadDir(userWorldsDir)
		if err != nil {
			log.Println("ERROR: Failed to read user worlds directory:", userWorldsDir, err)
			continue
		}
		for _, worldEntry := range worldEntries {
			if !worldEntry.IsDir() {
				log.Println("Skipping non-directory world entry:", worldEntry.Name())
				continue
			}
			worldPath := filepath.Join(userWorldsDir, worldEntry.Name())
			// World don't have manifest.json, so we read levelname.txt for the name
			// And icon is world_icon.jpeg
			// For resource and behavior packs, we read from world_resource_packs.json and world_behavior_packs.json respectively
			levelNamePath := filepath.Join(worldPath, "levelname.txt")
			levelNameBytes, err := os.ReadFile(levelNamePath)
			var worldName string
			if err != nil {
				log.Println("WARNING: Failed to read levelname.txt in:", worldPath, err)
				worldName = worldEntry.Name() // Fallback to folder name
			} else {
				worldName = strings.TrimSpace(string(levelNameBytes))
			}
			icon := filepath.Join(worldPath, "world_icon.jpeg")
			resourcePackUUIDs := getDependency(filepath.Join(worldPath, "world_resource_packs.json"), listOfResourcePacks)
			behaviorPackUUIDs := getDependency(filepath.Join(worldPath, "world_behavior_packs.json"), listOfBehaviorPacks)

			allWorlds = append(allWorlds, WorldData{
				CleanName:  worldName,
				Path:       worldPath,
				Icon:       icon,
				RequiredRP: resourcePackUUIDs,
				RequiredBP: behaviorPackUUIDs,
			})
		}
	}
	return allWorlds
}
func ListResourcePack(dir string) []ResourcePack {
	var list []ResourcePack

	scanPacks(dir, func(path string, name string, manifest Manifest, icon string) {
		// Specific Logic for Resource Packs
		list = append(list, ResourcePack{
			CleanName:    name,
			Icon:         icon,
			UUID:         manifest.Header.UUID,
			Path:         path,
			IsSignatures: fileExists(filepath.Join(path, "signatures.json")),
		})
	})
	return list
}
func ListBehaviorPack(dir string) []BehaviorPack {
	var list []BehaviorPack

	scanPacks(dir, func(path string, name string, manifest Manifest, icon string) {
		// Specific Logic for Behavior Packs
		var serverVersion string = "null-script"
		var dependenciesUUID []string

		for _, dep := range manifest.Dependencies {
			// Extract UUIDs
			if uuid, ok := dep["uuid"].(string); ok {
				dependenciesUUID = append(dependenciesUUID, uuid)
			}
			// Check for Script API version
			if module, ok := dep["module_name"].(string); ok && module == "@minecraft/server" {
				if ver, ok := dep["version"].(string); ok {
					serverVersion = ver
				}
			}
		}

		list = append(list, BehaviorPack{
			CleanName:        name,
			Icon:             icon,
			UUID:             manifest.Header.UUID,
			DependenciesUUID: dependenciesUUID,
			ScriptState:      serverVersion,
			Path:             path,
			IsSignatures:     fileExists(filepath.Join(path, "signatures.json")),
		})
	})
	return list
}
func calculateNameSimilarity(name1, name2 string) int {
	tokens1 := tokenize(name1)
	tokens2 := tokenize(name2)

	score := 0
	for _, t1 := range tokens1 {
		if len(t1) < 3 {
			continue
		}
		for _, t2 := range tokens2 {
			if len(t2) < 3 {
				continue
			}
			if t1 == t2 {
				score++
			}
		}
	}
	return score
}

func tokenize(s string) []string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, "-", " ")
	s = strings.ReplaceAll(s, "_", " ")
	// You can add more cleanup here (e.g., remove "v1", "mcpack", etc.)
	return strings.Fields(s)
}
func sanitizeRPs(allRPs []ResourcePack) ([]ResourcePack, []PackWarning) {
	var valid []ResourcePack
	var warnings []PackWarning

	// Map to track UUIDs: UUID -> Slice of RPs found with that UUID
	occurrenceMap := make(map[string][]ResourcePack)

	// Group by UUID
	for _, rp := range allRPs {
		occurrenceMap[rp.UUID] = append(occurrenceMap[rp.UUID], rp)
	}

	for uuid, group := range occurrenceMap {
		if len(group) == 1 {
			valid = append(valid, group[0])
		} else {
			// Found duplicates!
			// 1. Keep the first one as "valid" (arbitrary choice to prevent crash)
			valid = append(valid, group[0])

			// 2. Create Warning
			names := []string{}
			for _, rp := range group {
				names = append(names, rp.CleanName)
			}

			warnings = append(warnings, PackWarning{
				Type:           WarningDuplicateUUID,
				Message:        fmt.Sprintf("Duplicate Resource Pack UUID (%s) found in: %s. Using first one.", uuid, strings.Join(names, ", ")),
				ConflictingRPs: group,
			})
		}
	}
	return valid, warnings
}

func sanitizeBPs(allBPs []BehaviorPack) ([]BehaviorPack, []PackWarning) {
	var valid []BehaviorPack
	var warnings []PackWarning

	// Map: UUID -> List of packs with that UUID
	occurrenceMap := make(map[string][]BehaviorPack)

	for _, bp := range allBPs {
		occurrenceMap[bp.UUID] = append(occurrenceMap[bp.UUID], bp)
	}

	for uuid, group := range occurrenceMap {
		if len(group) == 1 {
			// Unique UUID, add to valid list
			valid = append(valid, group[0])
		} else {
			// Duplicate UUIDs found

			// 1. Keep the first one found (arbitrary selection)
			valid = append(valid, group[0])

			// 2. Collect names for logging/warning
			names := []string{}
			for _, bp := range group {
				names = append(names, bp.CleanName)
			}

			// --- LOGGING START ---
			// This prints the specific UUID and the names of the packs fighting for it
			log.Printf("⚠️ [DUPLICATE DETECTED] UUID: %s\n\t→ Found in %d packs: [%s]",
				uuid, len(group), strings.Join(names, ", "))
			// --- LOGGING END ---

			// 3. Add to Warning list
			warnings = append(warnings, PackWarning{
				Type:           WarningDuplicateUUID,
				Message:        fmt.Sprintf("Duplicate Behavior Pack UUID (%s) found in: %s. Using first one.", uuid, strings.Join(names, ", ")),
				ConflictingBPs: group,
			})
		}
	}
	return valid, warnings
}
func sortPacks(resourcePacks []ResourcePack, behaviorPacks []BehaviorPack) ([]ResourcePack, []Addon, []BehaviorPack, []PackWarning) {
	var warnings []PackWarning

	// --- STEP 0: Sanitize & Check for Duplicates ---
	// We clean the lists first so the rest of the logic doesn't bug out
	validRPs, rpWarnings := sanitizeRPs(resourcePacks)
	validBPs, bpWarnings := sanitizeBPs(behaviorPacks)

	warnings = append(warnings, rpWarnings...)
	warnings = append(warnings, bpWarnings...)

	// --- STEP 1: Indexing (Using Valid RPs) ---
	rpMap := make(map[string]ResourcePack, len(validRPs))
	rpOrder := make([]string, 0, len(validRPs))
	for _, rp := range validRPs {
		rpMap[rp.UUID] = rp
		rpOrder = append(rpOrder, rp.UUID)
	}

	// --- STEP 2: Bidding ---
	claims := make(map[string][]BehaviorPack)
	orphanBPCandidates := make(map[string]BehaviorPack)

	for _, bp := range validBPs {
		foundDependency := false
		for _, depUUID := range bp.DependenciesUUID {
			if _, exists := rpMap[depUUID]; exists {
				claims[depUUID] = append(claims[depUUID], bp)
				foundDependency = true
				break
			}
		}
		if !foundDependency {
			orphanBPCandidates[bp.UUID] = bp
		}
	}

	// --- STEP 3: Resolution ---
	var addonPacksLoop []Addon
	usedRP := make(map[string]bool)
	usedBP := make(map[string]bool)

	for rpUUID, bidders := range claims {
		rp := rpMap[rpUUID]
		var winnerBP BehaviorPack
		winnerScore := -1

		for _, bp := range bidders {
			score := calculateNameSimilarity(bp.CleanName, rp.CleanName)
			if score > winnerScore {
				winnerScore = score
				winnerBP = bp
			}
		}

		// Dependency Warnings
		if len(bidders) > 1 {
			bpNames := []string{}
			for _, b := range bidders {
				bpNames = append(bpNames, b.CleanName)
			}
			warnings = append(warnings, PackWarning{
				Type:         WarningConflict,
				Message:      fmt.Sprintf("Multiple BPs (%s) claimed RP '%s'. Winner: '%s'", strings.Join(bpNames, ", "), rp.CleanName, winnerBP.CleanName),
				ResourcePack: rp,
				InvolvedBPs:  bidders,
				WinnerBP:     winnerBP,
			})
		} else if winnerScore == 0 {
			warnings = append(warnings, PackWarning{
				Type:         WarningMismatch,
				Message:      fmt.Sprintf("BP '%s' claimed RP '%s' via UUID, but names look unrelated.", winnerBP.CleanName, rp.CleanName),
				ResourcePack: rp,
				InvolvedBPs:  bidders,
				WinnerBP:     winnerBP,
			})
		}

		addonPacksLoop = append(addonPacksLoop, Addon{
			CleanName:    cleanName(winnerBP.CleanName),
			ResourcePack: rp,
			BehaviorPack: winnerBP,
		})
		usedRP[rp.UUID] = true
		usedBP[winnerBP.UUID] = true
	}

	// --- STEP 4: Final Compilation ---
	var exclusiveResourcePacksLoop []ResourcePack
	for _, uuid := range rpOrder {
		if !usedRP[uuid] {
			exclusiveResourcePacksLoop = append(exclusiveResourcePacksLoop, rpMap[uuid])
		}
	}

	var exclusiveBehaviorPacksLoop []BehaviorPack
	for _, bp := range orphanBPCandidates {
		exclusiveBehaviorPacksLoop = append(exclusiveBehaviorPacksLoop, bp)
	}
	for _, bp := range validBPs {
		_, isOrphan := orphanBPCandidates[bp.UUID]
		if !usedBP[bp.UUID] && !isOrphan {
			exclusiveBehaviorPacksLoop = append(exclusiveBehaviorPacksLoop, bp)
		}
	}

	return exclusiveResourcePacksLoop, addonPacksLoop, exclusiveBehaviorPacksLoop, warnings
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
