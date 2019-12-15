package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/sqweek/dialog"
)

type Script struct {
	Name string
	Path string
}

type Lua struct {
	Name string
	Path string
}

var lualist []Lua
var metascript Script

func main() {
	Init()

	progress := spinner.New(spinner.CharSets[26], 500*time.Millisecond)

	progress = ReportProgress(progress, "Decrypting "+metascript.Name)
	progress.Start()
	Decrypt()
	progress.Stop()
	log.Println("Decrypting " + metascript.Name + " Completed!")

	progress = ReportProgress(progress, "Unpacking "+metascript.Name)
	progress.Start()
	Unpack()
	progress.Stop()
	log.Println("Unpacking " + metascript.Name + " Completed!")

	Unlock(progress)
	Decompile(progress)
	Clone(progress)
	Encrypt(progress)

	// log.Println("Repacking Scripts")
	// progress.Start()
	// Repack()
	// progress.Stop()
	// log.Println("Done!")

	log.Println("Finished!")
	// CleanWorkSpace()

	log.Print("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func Init() {

	fmt.Println("")
	fmt.Println(" Azur Lane Command Line Tool")
	fmt.Println(" Version v1.0.0")
	fmt.Println(" Credit & Thanks to github.com/k0np4ku")
	fmt.Println(" Continued Distribution github.com/AthanatiusC")
	fmt.Println("")

	luanames := []string{"BuildShip.lua.txt"}
	var lua Lua
	log.SetPrefix(" [ SYSTEM ] ")

	scripts, err := dialog.File().Filter("Scripts Asset Bundle", "*").Load()
	ErrorHandler(err)

	metascript.Path = scripts
	pathsplit := strings.Split(metascript.Path, "\\")
	metascript.Name = pathsplit[len(pathsplit)-1]
	// log.Println(metascript.Path)

	if metascript.Name == "" {
		return
	} else {
		WorkspaceMover()
	}

	for _, name := range luanames {
		lua.Name = name
		lua.Path = "Unity_Assets_Files\\" + metascript.Name + "\\CAB-android32\\" + lua.Name
		lualist = append(lualist, lua)
	}
}

func Decrypt() {
	args := []string{"--dev", "--decrypt", metascript.Path}
	cmd := exec.Command("Azcli.exe", args...)
	err := cmd.Run()
	ErrorHandler(err)
}

func Unpack() {
	args := []string{"--dev", "--unpack", metascript.Path}
	cmd := exec.Command("Azcli.exe", args...)
	err := cmd.Run()
	ErrorHandler(err)
}

func Unlock(progress *spinner.Spinner) {
	for _, lua := range lualist {
		progress = ReportProgress(progress, "Decrypting : "+lua.Name)
		args := []string{"--dev", "--unlock", lua.Path}
		cmd := exec.Command("Azcli.exe", args...)
		err := cmd.Run()
		ErrorHandler(err)
	}
	log.Println("Decrypting Lua Completed!")
}

func Decompile(progress *spinner.Spinner) {
	for _, lua := range lualist {
		progress = ReportProgress(progress, "Decompiling : "+lua.Name)
		args := []string{"--dev", "--decompile", lua.Path}

		cmd := exec.Command("Azcli.exe", args...)
		err := cmd.Run()
		ErrorHandler(err)
	}
	log.Println("Decompiling Lua Completed!")
}

func Clone(progress *spinner.Spinner) {
	os.Mkdir("Mod", 777)
	for _.lua := range lualist{
		
		srcFile, err := os.Open(lua.Path)
		ErrorHandler(err)
		defer srcFile.Close()

		lua.Path = "Workspace\\Mod\\" + lua.Name
		destFile, err := os.Create(metascript.Path) // creates if file doesn't exist
		ErrorHandler(err)
		defer destFile.Close()

		_, err = io.Copy(destFile, srcFile) // check first var for number of bytes copied
		ErrorHandler(err)

		err = destFile.Sync()
		ErrorHandler(err)
	}
}

func Encrypt(progress *spinner.Spinner) {
	for _, lua := range lualist {
		log.Println("Encrypting : " + lua.Name)
		args := []string{"--dev", "--lock", lua.Path}

		cmd := exec.Command("Azcli.exe", args...)
		err := cmd.Run()
		ErrorHandler(err)
	}
}

func Repack() {

}

func WorkspaceMover() {
	progress := spinner.New(spinner.CharSets[26], 500*time.Millisecond)

	progress = ReportProgress(progress, "Copying asset bundle to workspace")
	err := os.Mkdir("Workspace", 777)
	if err != nil {
		decision := dialog.Message("Working space is not empty!\nWould you like to delete it?").Title("Attention!").YesNo()
		if decision {

			progress = ReportProgress(progress, "Cleaning Existing Workspace")
			err := os.RemoveAll("Workspace")
			ErrorHandler(err)
			log.Println("Cleaning Completed!")

			progress = ReportProgress(progress, "Creating New Workspace")
			err = os.Mkdir("Workspace", 777)
			ErrorHandler(err)
			log.Println("Workspace Completed!")

		} else {
			log.Println("Please clean working directory first!")
			return
		}
	}
	// ErrorHandler(err)
	srcFile, err := os.Open(metascript.Path)
	ErrorHandler(err)
	defer srcFile.Close()

	metascript.Path = "Workspace\\" + metascript.Name
	destFile, err := os.Create(metascript.Path) // creates if file doesn't exist
	ErrorHandler(err)
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile) // check first var for number of bytes copied
	ErrorHandler(err)

	err = destFile.Sync()
	ErrorHandler(err)
}

func CleanWorkSpace() {
	_, err := os.Open("Workspace\\scripts32i")
	if os.IsExist(err) {
		log.Println("Move the file")
	} else {
		err := os.RemoveAll("Workspace")
		ErrorHandler(err)

		err = os.RemoveAll("Unity_Assets_Files")
		ErrorHandler(err)
	}
}

func ReportProgress(progress *spinner.Spinner, report string) *spinner.Spinner {
	progress.Prefix = " [ SYSTEM ] " + report
	return progress
}

func ErrorHandler(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
