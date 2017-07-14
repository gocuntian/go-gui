package main

import (
	"log"

	"github.com/sciter-sdk/go-sciter"
	"github.com/sciter-sdk/go-sciter/window"
)

// GO语言使用GO-SCITER创建桌面应用(四) 固定窗口大小

// 有些时候我们需要创建的应用窗口大小不可改变。
func main() {
	//创建窗口并设置大小
	w, err := window.New(sciter.DefaultWindowCreateFlag, &sciter.Rect{200, 200, 500, 500})
	if err != nil {
		log.Fatal(err)
	}
	//加载文件
	w.LoadFile("index.html")
	//设置标题
	w.SetTitle("固定大小窗口")
	//显示窗口
	w.Show()
	//运行窗口，进入消息循环
	w.Run()
}
