package main

import (
	"fmt"
	"log"

	"github.com/sciter-sdk/go-sciter"
	"github.com/sciter-sdk/go-sciter/window"
)

// GO语言使用GO-SCITER创建桌面应用(六) ELEMENT元素操作和EVENT事件响应
// 详细的文档请看下面两个链接:
// https://sciter.com/docs/content/sciter/Element.htm
// https://sciter.com/docs/content/sciter/Event.htm

func defFunc(w *window.Window) {
	//注册dump函数方便在tis脚本中打印数据
	w.DefineFunction("dump", func(args ...*sciter.Value) *sciter.Value {
		for _, v := range args {
			fmt.Print(v.String() + " ")
		}
		fmt.Println()
		return sciter.NullValue()
	})
	//注册reg函数，用于处理注册逻辑，这里只是简单的把数据打印出来
	w.DefineFunction("reg", func(args ...*sciter.Value) *sciter.Value {
		fmt.Println(args)
		for _, v := range args {
			fmt.Print(v.String() + " ")
		}
		fmt.Println()
		return sciter.NullValue()
	})
}

func main() {
	w, err := window.New(sciter.DefaultWindowCreateFlag, &sciter.Rect{300, 300, 700, 700})
	if err != nil {
		log.Fatal(err)
	}
	w.LoadFile("index.html")
	w.SetTitle("表单js提交")
	defFunc(w)
	w.Show()
	w.Run()
}
