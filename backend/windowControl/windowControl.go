package windowControl

import (
	"fmt"
	"runtime"
	"syscall"
	"time"

	"github.com/micmonay/keybd_event"
)

var (
	user32              = syscall.NewLazyDLL("user32.dll")
	getForegroundWindow = user32.NewProc("GetForegroundWindow")
	procMinimize        = user32.NewProc("ShowWindow")
	restoreWindow       = user32.NewProc("ShowWindow")
)

const (
	swMinimize    = 6
	SW_SHOWNORMAL = 1
)

// func main() {
// 	// NotMain()
// 	if err := restoreWindowToNormalSize(); err != nil {
// 		fmt.Printf("Failed to restore the window to its normal size: %s\n", err)
// 	}

// 	// Wait for 2 seconds before minimizing the window
// 	time.Sleep(2 * time.Second)

// 	if err := moveWindowToNextMonitor(); err != nil {
// 		fmt.Printf("Failed to restore the window to its normal size: %s\n", err)
// 	}

// 	time.Sleep(2 * time.Second)

// 	if err := Minimize(); err != nil {
// 		fmt.Printf("Failed to minimize the window: %s\n", err)
// 	}
// }

func Normalize() error {
	// Get the handle of the currently focused window
	hwnd, _, err := getForegroundWindow.Call()
	if hwnd == 0 {
		return fmt.Errorf("failed to get the handle of the currently focused window: %w", err)
	}

	// Restore the window to its normal size
	ret, _, err := restoreWindow.Call(hwnd, SW_SHOWNORMAL)
	if ret == 0 {
		return fmt.Errorf("failed to restore the window to its normal size: %w", err)
	}

	fmt.Println("Window restored to its normal size successfully.")
	return nil
}

func Minimize() error {
	// Get the handle of the foreground window
	hWnd, _, err := getForegroundWindow.Call()
	if hWnd == 0 {
		return fmt.Errorf("failed to get the handle of the foreground window: %w", err)
	}

	// Minimize the window after a delay
	ret, _, err := procMinimize.Call(hWnd, uintptr(swMinimize))
	if ret == 0 {
		return fmt.Errorf("failed to minimize the window: %w", err)
	}

	fmt.Println("Window minimized successfully.")
	return nil
}

func Move() error {
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		panic(err)
	}

	// For linux, it is very important to wait 2 seconds
	if runtime.GOOS == "linux" {
		time.Sleep(2 * time.Second)
	}

	// Select keys to be pressed
	kb.SetKeys(keybd_event.VK_RIGHT)

	// Set shift to be pressed
	kb.HasSuper(true)

	// Set shift to be pressed
	kb.HasSHIFT(true)

	// Press the selected keys
	err = kb.Launching()
	if err != nil {
		panic(err)
	}

	return nil
}
