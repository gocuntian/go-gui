package main

import (
	"log"

	sciter "github.com/sciter-sdk/go-sciter"
	"github.com/sciter-sdk/go-sciter/window"
	"github.com/xingcuntian/go-gui/Toolv1.1/controllers"
)

func main() {
	// NewRect(top, left, width, height int)
	w, err := window.New(sciter.DefaultWindowCreateFlag, sciter.NewRect(120, 400, 640, 460))
	if err != nil {
		log.Fatal(err)
	}
	w.LoadFile("views/login.html")
	w.SetTitle("用户登录")
	controllers.SetEventHandler(w)
	w.Show()
	w.Run()
}
