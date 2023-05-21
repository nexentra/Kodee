package main

import (
	"context"
	_ "embed"
	"fmt"
	"io/ioutil"

	"github.com/KnockOutEZ/Kodee/backend/healthReminder"
	"github.com/KnockOutEZ/Kodee/backend/server"
	"github.com/KnockOutEZ/Kodee/backend/systemUsage"
	"github.com/KnockOutEZ/Kodee/backend/utils"
	"github.com/KnockOutEZ/Kodee/backend/weatherApi"
	"github.com/KnockOutEZ/Kodee/backend/windowControl"
	// "github.com/getlantern/systray"
	"fyne.io/systray"
	"github.com/sadlil/go-avro-phonetic"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var myCtx context.Context

//go:embed frontend/src/assets/wails.ico
var logoICO []byte

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called at application startup
func (a *App) startup(ctx context.Context) {
	// Perform your setup here
	a.ctx = ctx
	systray.Run(onReady, onExit)
}

// domReady is called after front-end resources have been loaded
func (a *App) domReady(ctx context.Context) {
	// Add your action here
	myCtx = ctx
	utils.CopyIconInStartup()
	ConvertBanglishToBangla("ki obostha korechi")

	// server.TestAuth()
	server.TestOut()
	// call reminderFuncs here
	go healthReminder.LookAwayReminder()
	go healthReminder.StandUpReminder()
	go healthReminder.HydrateReminder()

}

func ConvertBanglishToBangla(text string) string {
	// Parse() tries to parse the given string
	// In case of failure it returns an erros
	text, err := avro.Parse("ami banglay gan gai")
	if err != nil {
		return err.Error()
	}

	fmt.Println(text) // আমি বাংলায় গান গাই
	return text
}

func onReady() {
	fmt.Println("onReady")
	systray.SetIcon(logoICO)
	systray.SetTitle("Kodee")
	systray.SetTooltip("Kodee")
	mOpen := systray.AddMenuItem("Open", "Open the app")
	go func() {
		for {
			<-mOpen.ClickedCh
			runtime.Show(myCtx)
		}
	}()
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")
	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
		runtime.Quit(myCtx)
	}()
	// Sets the icon of a menu item.
	// mQuit.SetIcon(logoICO)
}

func onExit() {
	// clean up here
}

func getIcon(s string) []byte {
	b, err := ioutil.ReadFile(s)
	utils.CheckErr(err)
	return b
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// beforeClose is called when the application is about to quit,
// either by clicking the window close button or calling runtime.Quit.
// Returning true will cause the application to continue, false will continue shutdown as normal.
func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	// dialog, err := runtime.MessageDialog(ctx, runtime.MessageDialogOptions{
    //     Type:          runtime.QuestionDialog,
    //     Title:         "Quit?",
    //     Message:       "Are you sure you want to quit?",
    // })

    // if err != nil {
        return false
    // }
    // return dialog != "Yes"
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
	// Perform your teardown here
}

// Meet returns a greeting for the given name
func (a *App) Notification(title, message string) {
	utils.NotificationFunc(title, message)
}

func (a *App) GetSystemUsage(usageType string) interface{} {
	switch usageType {
	case "cpu":
		return systemUsage.GetCpuUsage()
	case "ram":
		return systemUsage.GetRamUsage()
	case "bandwidth":
		return systemUsage.GetBandwithSpeed()
	default:
		return nil
	}
}

func (a *App) UseWindowControl(usageType string) interface{} {
	switch usageType {
	case "minimize":
		if err := windowControl.Minimize(); err != nil {
			return err
		}
	case "normalize":
		if err := windowControl.Normalize(); err != nil {
			return err
		}
	case "move":
		if err := windowControl.Move(); err != nil {
			return err
		}
	default:
		return nil
	}
	return nil
}

func (a *App) Auth(usageType string, user server.User, jwt string) (interface{}, error) {
	fmt.Println("Auth called", usageType, user, jwt)
	switch usageType {
	case "signup":
		userData := &server.User{Username: user.Username, Email: user.Email, Password: user.Password}
		res, err := server.SignUp(userData)
		if err != nil {
			return nil, err
		}
		return res, nil
	case "signin":
		userData := server.User{Username: user.Username, Password: user.Password}
		res, err := server.Login(userData.Username, userData.Password)
		if err != nil {
			return nil, err
		}
		return res, nil
	case "me":
		res, err := server.Me(jwt)
		if err != nil {
			return nil, err
		}
		return res, nil
	default:
		return nil, nil
	}
}

func (a *App) GetWeather() {
	weatherApi.GetWeather()
}
