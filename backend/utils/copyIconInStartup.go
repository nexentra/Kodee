package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func CopyIconInStartup() {
    dirname, err := os.UserHomeDir()
    CheckErr(err)

    desktopPath := filepath.Join(dirname, "Desktop", "kodee.lnk")
    if _, err := os.Stat(desktopPath); os.IsNotExist(err) {
        fmt.Println("kodee.lnk not found on desktop")
        return
    }

    in, err := os.Open(desktopPath)
    CheckErr(err)
    defer in.Close()

    out, err := os.Create(filepath.Join(dirname, "AppData", "Roaming", "Microsoft", "Windows", "Start Menu", "Programs", "Startup", "kodee.lnk"))
    CheckErr(err)
    defer out.Close()

    _, err = io.Copy(out, in)
    CheckErr(err)
}


// func CopyIconInStartup() {
// 	dirname, err := os.UserHomeDir()
// 	CheckErr(err)
// 	in, err := os.Open(dirname + `\Desktop\kodee.lnk`)
// 	CheckErr(err)
// 	defer in.Close()

// 	out, err := os.Create(dirname + `\AppData\Roaming\Microsoft\Windows\Start Menu\Programs\Startup\kodee.lnk`)
// 	CheckErr(err)
// 	defer out.Close()

// 	_, err = io.Copy(out, in)
// 	CheckErr(err)
// }