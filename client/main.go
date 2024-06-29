package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/schollz/progressbar/v3"
)

const (
	ConfigPath = "/etc/caprice"
)

const (
	MMEra_Feb9 = 0
	MMEra_010  = 1
	MMERA_100  = 2
)

type Personality struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	Version      string `json:"version"`
	ID           string `json:"id"`
	MMEra        int    `json:"mmera"`
	BluetoothEra int    `json:"bluetoothera"`
	// needs an extra daemon
	CustomWireProgram bool `json:"customwireprogram"`
}

func FExist(file string) bool {
	_, err := os.ReadFile(file)
	return err == nil
}

func GetFile(file string) []byte {
	fileBytes, _ := os.ReadFile(file)
	return fileBytes
}

func GetFileAsString(file string) string {
	fileBytes, _ := os.ReadFile(file)
	return string(fileBytes)
}

func OSDetect() (isStock bool, version string, id string) {
	currentFile := fjoin(ConfigPath, "current.json")
	ankiCurrentFile := "/anki/etc/version"
	var pers Personality
	if FExist(fjoin(ConfigPath, "current.json")) {
		json.Unmarshal(GetFile(currentFile), &pers)
		return false, pers.Version, pers.ID
	}
	return true, GetFileAsString(ankiCurrentFile), "stock"
}

func main() {
	stock, ver, id := OSDetect()
	fmt.Println("Current Personality Version: " + ver)
	if stock {
		fmt.Println("Personality is stock")
	} else {
		fmt.Println("Personality is NOT stock")
		fmt.Println("Current Personality details:")
		fmt.Println("  - ID: ", id)
		fmt.Println("  - Version: ", ver)
	}
	fmt.Println("Loading bar test")
	bar := progressbar.Default(100)
	for i := 0; i < 100; i++ {
		bar.Add(1)
		time.Sleep(40 * time.Millisecond)
	}
}

func fjoin(in ...string) string {
	return filepath.Join(in...)
}
