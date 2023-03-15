package main

import (
	"downloader/load"
	"flag"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"html/template"
	"image/color"
	"net/http"
	"os"
	"sync"
	"time"
)

var url = flag.String("u", "", "The download url")

func main() {
	//webMode()
	//consoleMode()
	guiMode1()
}

/*
Web方式
*/
func webMode() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		//writer.Write([]byte("aaa"))
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Header().Add("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		writer.Header().Add("Access-Control-Allow-Credentials", "true")
		writer.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		writer.Header().Set("content-type", "application/json;charset=UTF-8")
		writer.WriteHeader(http.StatusOK)

		t, err := template.ParseFiles("web/index.html")
		if err != nil {
			writer.Write([]byte(err.Error()))
			return
		}
		//t.Execute(writer, nil)
		t.ExecuteTemplate(writer, "layout", nil)
		//t, _ := template.ParseFiles("web/index.html", "web/left.html", "web/right.html")
		//t.ExecuteTemplate(writer, "layout", nil)
	})
	http.ListenAndServe(":8080", nil)
}

/*
控制台方式下载
*/
func consoleMode() {
	flag.Parse()
	if *url == "" {
		fmt.Fprintf(os.Stderr, "The url can not empty")
		return
	}
	a := sync.WaitGroup{}
	a.Add(1)
	go func() {
		//load.Run("https://dl.moapp.me/https://github.com/agalwood/Motrix/releases/download/v1.6.11/Motrix-Setup-1.6.11.exe")
		//load.Run("http://cdimage.deepin.com/releases/20.3/deepin-desktop-community-20.3-amd64.iso")
		load.Header(*url)
		//load.Header("https://download.jetbrains.com.cn/go/goland-2022.3.2.exe")
		a.Done()
	}()

	a.Wait()
}

func guiMode() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Border Layout")

	left := widget.NewLabel("left")
	right := widget.NewLabel("right")
	top := widget.NewLabel("top")
	bottom := widget.NewLabel("bottom")
	content := widget.NewLabel(`Lorem ipsum dolor, 
  sit amet consectetur adipisicing elit.
  Quidem consectetur ipsam nesciunt,
  quasi sint expedita minus aut,
  porro iusto magnam ducimus voluptates cum vitae.
  Vero adipisci earum iure consequatur quidem.`)
	_ = left

	leftBack := canvas.NewRectangle(color.Black)
	leftBack.SetMinSize(fyne.NewSize(200, 0))
	labelLeft := widget.NewLabel("left")
	leftContent := container.New(layout.NewVBoxLayout(), labelLeft)
	leftLayout := container.New(layout.NewMaxLayout(), leftBack, leftContent)

	bottomLayout := container.New(layout.NewHBoxLayout(), layout.NewSpacer(), bottom)

	container := container.New(
		layout.NewBorderLayout(top, bottomLayout, leftLayout, right),
		top, bottomLayout, leftLayout, right, content,
	)
	myWindow.Resize(fyne.Size{Width: 1000, Height: 800})
	go getCurrentTime(bottom)
	myWindow.SetContent(container)
	myWindow.ShowAndRun()
}

/* 获取当前时间 时间格式化参考 https://blog.51cto.com/u_15891990/5908904*/
func getCurrentTime(label *widget.Label) {
	for {
		currentTime := time.Now()
		currentTimeStr := currentTime.Format("2006-01-02 15:04:05")
		fmt.Println(currentTimeStr)
		label.SetText(currentTimeStr)
		time.Sleep(time.Second)
	}
}

func guiMode1() {
	app := app.New()
	window := app.NewWindow("下载器")
	window.Resize(fyne.NewSize(1050, 750))
	window.SetMaster()
	window.SetPadded(false)

	leftBack := canvas.NewRectangle(color.RGBA{51, 51, 51, 255})
	leftBack.SetMinSize(fyne.NewSize(80, 0))
	var icn []byte
	fyne.Resource()
	s := fyne.NewStaticResource("resource", icn)
	fmt.Println(s)
	_ = s
	iconAdd := widget.NewIcon(fyne.NewStaticResource("resource", icn))
	addAction := widget.NewToolbarAction(fyne.NewStaticResource("resource", icn), func() {

	})

	toolbar := widget.NewToolbar(widget.NewToolbarSpacer(), addAction, widget.NewToolbarSpacer())
	leftContent := container.New(layout.NewVBoxLayout(), toolbar, iconAdd)
	leftLayout := container.New(layout.NewMaxLayout(), leftBack, leftContent)

	middleBack := canvas.NewRectangle(color.RGBA{244, 245, 247, 255})
	middleBack.SetMinSize(fyne.NewSize(200, 0))
	middleContent := container.New(layout.NewVBoxLayout())
	middleLayout := container.New(layout.NewMaxLayout(), middleBack, middleContent)
	rootLayout := container.New(layout.NewHBoxLayout(), leftLayout, middleLayout)
	window.SetContent(rootLayout)

	window.ShowAndRun()
}

func guiMode2() {
	myApp := app.New()

	myWindow := myApp.NewWindow("Demo")

	myWindow.SetMaster()

	myWindow.SetPadded(false)

	myWindow.Resize(fyne.NewSize(1024, 600))

	//myWindow.SetFullScreen(true)

	r1 := canvas.NewRectangle(color.RGBA{255, 0, 0, 255})

	r1.SetMinSize(fyne.NewSize(80, 80))

	top := fyne.NewContainerWithLayout(layout.NewMaxLayout(), r1)

	content := container.New(layout.NewBorderLayout(nil, nil, top, nil), top)

	myWindow.SetContent(content)

	myWindow.ShowAndRun()
}
