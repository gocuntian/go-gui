package main

import (
	"log"

	"github.com/sciter-sdk/go-sciter"
	"github.com/sciter-sdk/go-sciter/window"
	"github.com/xingcuntian/go-gui/Toolv1.1/controllers"
)

func main() {
	w, err := window.New(sciter.DefaultWindowCreateFlag, &sciter.Rect{448, 184, 1000, 540})
	if err != nil {
		log.Fatal(err)
	}
	w.LoadFile("views/login_en.html")
	w.SetTitle("用户登录")
	controllers.SetEventHandler(w)
	w.Show()
	w.Run()
}
