package main

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"image/color"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/kercre123/caprice/pkg/shared"
	"github.com/kercre123/vector-gobot/pkg/vscreen"
	"github.com/schollz/progressbar/v3"
)

const (
	ConfigPath = "/etc/caprice"
)

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
	var pers shared.Personality
	if FExist(fjoin(ConfigPath, "current.json")) {
		json.Unmarshal(GetFile(currentFile), &pers)
		return false, pers.Version, pers.ID
	}
	return true, GetFileAsString(ankiCurrentFile), "stock"
}

func runCmd(cmds ...string) {
	for _, cmd := range cmds {
		exec.Command("/bin/bash", "-c", cmd).Run()
	}
}

func sysKill(processes ...string) {
	for _, p := range processes {
		exec.Command("/bin/bash", "-c", "systemctl kill -s SIGKILL "+p).Run()
	}
}

func StopAnki() {
	runCmd("killall -9 vic-*", "pkill vic-")
	sysKill("vic-robot", "vic-engine", "vic-cloud", "vic-anim", "vic-gateway", "vic-neuralnets", "vic-crashuploader")
	time.Sleep(time.Second / 3)
}

func StartAnki() {
	runCmd("systemctl start anki-robot.target")
}

func ShowInstalling(display string, stop chan bool, progress chan int, stopped chan bool) {
	vscreen.InitLCD()
	vscreen.BlackOut()
	var prog int
	var doStop bool
	//.oO@*
	var progSpinner []string = []string{".", "o", "O", "@", "*"}
	var progSpinnerIn int
	go func() {
		for range stop {
			doStop = true
			break
		}
	}()
	go func() {
		for i := range progress {
			if i == 999 {
				break
			}
			prog = i
		}
	}()
	var progBuffer int
	for {
		if progSpinnerIn == 4 {
			if progBuffer != 4 {
				progBuffer++
			} else {
				progBuffer = 0
				progSpinnerIn = 0
			}
		} else {
			if progBuffer != 4 {
				progBuffer++
			} else {
				progBuffer = 0
				progSpinnerIn++
			}
		}
		if doStop {
			break
		}
		scrnImg := vscreen.CreateTextImageFromLines([]vscreen.Line{
			{
				Color: color.RGBA{0, 255, 0, 255},
				Text:  display,
			},
			{
				Color: color.RGBA{0, 255, 0, 255},
				Text:  "",
			},
			{
				Color: color.RGBA{255, 255, 255, 255},
				Text:  progSpinner[progSpinnerIn] + " " + fmt.Sprint(prog) + "%",
			},
		},
		)
		vscreen.SetScreen(scrnImg)
		time.Sleep(time.Millisecond * 33)
	}
	vscreen.SetScreen(vscreen.CreateTextImage("Installation is complete."))
	progress <- 999
	stopped <- true
}

func TestInstalling() {
	stop := make(chan bool)
	progress := make(chan int)
	stopped := make(chan bool)
	go ShowInstalling("Installing <version>...", stop, progress, stopped)
	bar := progressbar.Default(100)
	for i := 0; i < 100; i++ {
		bar.Add(1)
		progress <- i
		time.Sleep(40 * time.Millisecond)
	}
	stop <- true
	for range stopped {
		return
	}
}

func downloadAndExtract(url, root string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}
	stop := make(chan bool)
	progressChan := make(chan int)
	stopped := make(chan bool)
	go ShowInstalling("Installing <version>...", stop, progressChan, stopped)

	bar := pb.Full.Start64(resp.ContentLength)
	barReader := bar.NewProxyReader(resp.Body)

	gz, err := gzip.NewReader(barReader)
	if err != nil {
		return err
	}
	defer gz.Close()

	tarReader := tar.NewReader(gz)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(filepath.Join(root, header.Name), 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			outFile, err := os.Create(filepath.Join(root, header.Name))
			if err != nil {
				return err
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return err
			}
			outFile.Close()
		default:
			return fmt.Errorf("unknown type: %v in %s", header.Typeflag, header.Name)
		}

		progress := int(float64(bar.Current()) / float64(bar.Total()) * 100)
		progressChan <- progress
	}

	stop <- true
	for range stopped {
		break
	}

	bar.Finish()
	return nil
}

func DealWithBlobs(p shared.Personality) {
	runCmd("systemctl stop mm-anki-camera mm-qcamera-daemon")
}

func RemoveOldAnki() {
	runCmd("rm -rf /anki")
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
	StopAnki()
	fmt.Println("Removing old Anki... (this will take a while)")
	RemoveOldAnki()
	fmt.Println("Downloading new Anki...")
	downloadAndExtract("http://192.168.1.105:8000/test.tar.gz", "/")
	//StartAnki()
	fmt.Println("done!")
}

func fjoin(in ...string) string {
	return filepath.Join(in...)
}
