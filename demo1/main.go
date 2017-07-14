package main

import (
	"log"

	"github.com/sciter-sdk/go-sciter"
	"github.com/sciter-sdk/go-sciter/window"
)

// GO语言使用GO-SCITER创建桌面应用(一) 简单的通过HTML,CSS写UI

// 我们使用go-sciter，就不得不提Sciter，Sciter 是一个嵌入式的 HTML/CSS/脚本引擎，旨在为桌面应用创建一个 UI 框架层。

// 说简单点就是我们通过它可以像写html,css那样写桌面UI。
// 一、环境准备
// 第一步：从https://sciter.com/download/地址下载sciter-sdk
// 解压，找到sciter-sdk\bin\64\sciter.dll复制到c:\windows\system32
// 注意上面的请根据你自已的系统选择相应文件
// 第二步：由于使用到cgo，所以window下需要安装mingw或tdm-gcc(建议安装tdm-gcc)
// 下载地址：
// https://sourceforge.net/projects/mingw/files/
// 下载地址：
// http://tdm-gcc.tdragon.net/download
// 下载：mingw-get-setup.exe或tdm64-gcc-5.1.0-2.exe
// 安装，然后把mingw\bin或tdm-gcc\bin加入到环境变量中
// 第三步：cmd进入gopath目录并运行
// go get -x github.com/sciter-sdk/go-sciter
//

// 二、通过html,css编写简单UI

func main() {
	//创建window窗口
	//参数一表示创建窗口的样式
	//SW_TITLEBAR 顶层窗口，有标题栏
	//SW_RESIZEABLE 可调整大小
	//SW_CONTROLS 有最小/最大按钮
	//SW_MAIN 应用程序主窗口，关闭后其他所有窗口也会关闭
	//SW_ENABLE_DEBUG 可以调试
	//参数二表示创建窗口的矩形
	w, err := window.New(sciter.SW_TITLEBAR|
		sciter.SW_RESIZEABLE|
		sciter.SW_CONTROLS|
		sciter.SW_MAIN|
		sciter.SW_ENABLE_DEBUG,
		nil)
	if err != nil {
		log.Fatal(err)
	}
	//加载文件
	w.LoadFile("index.html")
	//设置标题
	w.SetTitle("你好，世界")
	//显示窗口
	w.Show()
	//运行窗口，进入消息循环
	w.Run()
}
