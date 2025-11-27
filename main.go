package main

import (
	"embed"
	"log"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

// Wails uses Go's `embed` package to embed the frontend files into the binary.
// Any files in the frontend/dist folder will be embedded into the binary and
// made available to the frontend.
// See https://pkg.go.dev/embed for more information.

//go:embed all:frontend/dist
var assets embed.FS

// main function serves as the application's entry point. It initializes the application, creates a window,
// and starts a goroutine that emits a time-based event every second. It subsequently runs the application and
// logs any error that might occur.
func main() {

	// Create a new Wails application by providing the necessary options.
	// Variables 'Name' and 'Description' are for application metadata.
	// 'Assets' configures the asset server with the 'FS' variable pointing to the frontend files.
	// 'Bind' is a list of Go struct instances. The frontend has access to the methods of these instances.
	// 'Mac' options tailor the application when running an macOS.
	appService := &App{}
	app := application.New(application.Options{
		Name:        "main",
		Description: "A demo of using raw HTML & CSS",
		Services: []application.Service{
			application.NewService(appService),
			application.NewService(&GreetService{}),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	// Create a new window with the necessary options.
	// 'Title' is the title of the window.
	// 'Mac' options tailor the window when running on macOS.
	// 'BackgroundColour' is the background colour of the window.
	// 'URL' is the URL that will be loaded into the webview.
	mainWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: "Window 1",
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              "/",
		Width:            1240,
		Height:           850,
	})

	// 窗口关闭时隐藏到托盘而不是退出
	mainWindow.RegisterHook(events.Common.WindowClosing, func(event *application.WindowEvent) {
		log.Println("窗口关闭，隐藏到系统托盘")
		mainWindow.Hide()
		event.Cancel()
	})

	appService.app = app

	// 创建系统托盘
	systray := app.SystemTray.New()

	// 读取并设置托盘图标
	// iconData, err := os.ReadFile("build/windows/icon.ico")
	// if err != nil {
	// 	// 如果读取失败，尝试使用 favicon
	// 	iconData, _ = assets.ReadFile("frontend/dist/favicon.ico")
	// }
	// if iconData != nil {
	// 	systray.SetIcon(iconData)
	// }

	// 设置托盘提示文本
	// systray.SetTooltip(appConfig.AppName)

	// 创建托盘菜单
	trayMenu := app.NewMenu()
	trayMenu.Add("显示主窗口").OnClick(func(ctx *application.Context) {
		log.Println("托盘菜单：显示主窗口")
		mainWindow.Show()
		mainWindow.Focus()
	})
	trayMenu.AddSeparator()
	trayMenu.Add("退出").OnClick(func(ctx *application.Context) {
		log.Println("托盘菜单：退出应用")
		app.Quit()
	})

	systray.SetMenu(trayMenu)

	// 单击托盘图标显示主窗口
	systray.OnClick(func() {
		log.Println("托盘图标被点击")
		mainWindow.Show()
		mainWindow.Focus()
	})

	log.Println("系统托盘创建成功")

	// Create a goroutine that emits an event containing the current time every second.
	// The frontend can listen to this event and update the UI accordingly.
	go func() {
		for {
			now := time.Now().Format(time.RFC1123)
			app.Event.Emit("time", now)
			time.Sleep(time.Second)
		}
	}()

	// Run the application. This blocks until the application has been exited.
	err := app.Run()

	// If an error occurred while running the application, log it and exit.
	if err != nil {
		log.Fatal(err)
	}
}
