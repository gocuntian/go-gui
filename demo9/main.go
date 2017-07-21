package main

import (
	"fmt"
	"log"

	"github.com/sciter-sdk/go-sciter"
	"github.com/sciter-sdk/go-sciter/window"
	"github.com/xingcuntian/go-gui/demo9/controllers"
)

func main() {
	w, err := window.New(sciter.DefaultWindowCreateFlag, &sciter.Rect{448, 184, 1000, 540})
	if err != nil {
		log.Fatal(err)
	}
	w.LoadFile("views/login.html")
	w.SetTitle("用户登录")
	controllers.SetEventHandler(w)
	w.Show()
	fmt.Println("start run")
	w.Run()
	fmt.Println("end run !!!")
}
