package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/sqweek/dialog"
)

func main() {

	script, err := dialog.File().Filter("Scripts Asset Bundle", "*").Load()
	ErrorHandler(err)
	pathsplit := strings.Split(script, "\\")
	filename := pathsplit[len(pathsplit)-1]
	script = WorkspaceMover(script, filename)
	lualist := []string{"BuildShip.lua.txt"}

	progress := spinner.New(spinner.CharSets[26], 500*time.Millisecond)
	progress.Prefix = "Please wait"

	log.Println("Decrypting Scripts")
	progress.Start()
	Decrypt(script)
	progress.Stop()
	log.Println("Done!")

	log.Println("Unpacking Scripts")
	progress.Start()
	Unpack(script)
	progress.Stop()
	log.Println("Done!")

	log.Println("Decrypting Scripts")
	progress.Start()
	Unlock(lualist)
	progress.Stop()
	log.Println("Done!")

	log.Println("Decompiling Scripts")
	progress.Start()
	Decompile(lualist)
	progress.Stop()
	log.Println("Done!")

	log.Println("Recompiling Scripts")
	progress.Start()
	Repack()
	progress.Stop()
	log.Println("Done!")

	log.Println("Repacking Scripts")
	progress.Start()
	Repack()
	progress.Stop()
	log.Println("Done!")

	log.Println("Finished!")

	log.Print("\n")
	log.Print("\n")
	log.Print("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	CleanWorkSpace()
}

func Decrypt(path string) {
	args := []string{"--dev", "--decrypt", path}
	cmd := exec.Command("Azcli.exe", args...)
	err := cmd.Run()
	ErrorHandler(err)
}

func Unpack(path string) {
	args := []string{"--dev", "--unpack", path}
	cmd := exec.Command("Azcli.exe", args...)
	err := cmd.Run()
	ErrorHandler(err)
}

func Unlock(lualist []string) {
	for _, lua := range lualist {
		log.Println("Decrypting : " + lua)
		args := []string{"--dev", "--unlock", "Workspace\\Unity_Assets_Files\\scripts32i\\CAB-android32\\" + lua}
		cmd := exec.Command("Azcli.exe", args...)
		err := cmd.Run()
		ErrorHandler(err)
	}
}

func Decompile(lualist []string) {
	for _, lua := range lualist {
		log.Println("Decompiling : " + lua)
		args := []string{"--dev", "--decompile", "Workspace\\Unity_Assets_Files\\scripts32i\\CAB-android32\\" + lua}

		cmd := exec.Command("Azcli.exe", args...)
		err := cmd.Run()
		ErrorHandler(err)
	}
}

func Lock() {

}

func Repack() {

}

func Encrypt() {

}

func WorkspaceMover(path string, name string) string {
	ShowMessage("Copying asset bundle to workspace")
	err := os.Mkdir("Workspace", 777)
	if err != nil {
		decision := dialog.Message("Working space is not empty!\nWould you like to delete it?").Title("Attention!").YesNo()
		if decision {
			ShowMessage("Cleaning...")
			err := os.RemoveAll("Workspace")
			ErrorHandler(err)
		} else {
			ShowMessage("Please clean working directory first!")
		}
	}
	// ErrorHandler(err)
	srcFile, err := os.Open(name)
	ErrorHandler(err)
	defer srcFile.Close()

	destFile, err := os.Create("Workspace\\" + name) // creates if file doesn't exist
	ErrorHandler(err)
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile) // check first var for number of bytes copied
	ErrorHandler(err)

	err = destFile.Sync()
	ErrorHandler(err)
	return "Workspace\\" + name
}

func CleanWorkSpace() {
	_, err := os.Open("Workspace\\scripts32i")
	if os.IsExist(err) {
		ShowMessage("Move the file")
	} else {
		err := os.RemoveAll("Workspace")
		ErrorHandler(err)
	}
}

func ShowMessage(message string) {
	log.Println(message)
}

func ErrorHandler(err error) {
	if err != nil {
		log.Println(err)
	}
}
