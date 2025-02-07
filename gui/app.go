package gui

import (
	"context"
	"fmt"
	"github.com/ergochat/irc-go/ircmsg"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"newirc/irc"
)

type App struct {
	Ctx         context.Context
	Connections []*irc.Client
}

func NewApp() *App {
	return &App{}
}

func (a *App) Startup(ctx context.Context) {
	a.Ctx = ctx
}

func (a *App) applicationMenu() *menu.Menu {
	AppMenu := menu.NewMenu()
	FileMenu := AppMenu.AddSubmenu("File")
	FileMenu.AddText("Quit", keys.CmdOrCtrl("q"), func(_ *menu.CallbackData) {
		runtime.Quit(a.Ctx)
	})
	return AppMenu
}

func (a *App) Connect(server irc.Server, profile irc.Profile) {
	fmt.Println(server)
	fmt.Println(profile)
	client := &irc.Client{}
	client.Server = fmt.Sprintf("%s:%d", server.Hostname, server.Port)
	client.UseTLS = server.TLS
	client.SASLLogin = server.Saslusername
	client.SASLPassword = server.Saslpassword
	client.Nick = profile.Nick
	client.Debug = true
	err := client.Connect()
	fmt.Println(err)
	go func() {
		client.Loop()
	}()
	a.Connections = append(a.Connections, client)
}

func (a *App) ExportTypesToWailsRuntime(ircmsg.Message) {}
