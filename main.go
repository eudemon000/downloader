package main

import (
	"downloader/load"
	"flag"
	"fmt"
	"html/template"
	"image/color"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type DownApp struct {
	mainWindow    *fyne.Window
	pathItem      *widget.FormItem
	pathEntry     *widget.Entry
	newTaskDialog *dialog.Dialog
	dataLog       *log.Logger
	errLog        *log.Logger
}

var d *DownApp

var url = flag.String("u", "", "The download url")

var middleItemText = []string{"Downloading", "Waiting", "Stop"}

func main() {
	//webMode()
	//consoleMode()
	d = &DownApp{
		dataLog: log.New(os.Stdout, "SUCCESS ", log.Ldate|log.Ltime|log.Lshortfile),
		errLog:  log.New(os.Stderr, "FAIL ", log.Ldate|log.Ltime|log.Lshortfile|log.Lmsgprefix),
	}
	//d.guiMode1()
	d.gui()
	//guiMode()
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

func (d *DownApp) guiMode1() {
	app := app.New()
	window := app.NewWindow("Downloader")
	window.Resize(fyne.NewSize(1050, 750))
	window.SetMaster()
	window.SetPadded(false)
	pathEntry := widget.NewEntry()
	pathItem := &widget.FormItem{Text: "Path", Widget: pathEntry}
	d.pathItem = pathItem
	d.pathEntry = pathEntry
	formItems := []*widget.FormItem{pathItem}
	leftBack := canvas.NewRectangle(color.RGBA{51, 51, 51, 255})
	leftBack.SetMinSize(fyne.NewSize(80, 0))
	addAction := widget.NewToolbarAction(resourceResourceAddPng, func() {
		newDialog := dialog.NewForm("New", "Submit", "Cancel", formItems, func(a bool) {
			url := d.pathEntry.Text
			load.Header(url)
		}, window)
		newDialog.Resize(fyne.NewSize(500, 200))
		newDialog.Show()
	})

	toolbar := widget.NewToolbar(widget.NewToolbarSpacer(), addAction, widget.NewToolbarSpacer())
	leftContent := container.New(layout.NewVBoxLayout(), toolbar)
	//leftContent.Resize(fyne.NewSize(80, 60))
	leftLayout := container.New(layout.NewMaxLayout(), leftBack, leftContent)

	middleBack := canvas.NewRectangle(color.RGBA{244, 245, 247, 255})
	middleBack.SetMinSize(fyne.NewSize(200, 0))
	middleContent := container.New(layout.NewMaxLayout())
	middleLayout := container.New(layout.NewMaxLayout(), middleBack, middleContent)
	rootLayout := container.New(layout.NewHBoxLayout(), leftLayout, middleLayout)
	window.SetContent(rootLayout)

	window.ShowAndRun()
}

func (d *DownApp) createMiddleItem() fyne.CanvasObject {
	emptyLabel := widget.NewLabel("")
	emptyLabel.Resize(fyne.NewSize(20, 50))
	titleLabel := widget.NewLabel("Tasks")
	titleLabel.Resize(fyne.NewSize(50, 30))

	list := widget.NewList(
		func() int {
			return len(middleItemText)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(middleItemText[i])

		})
	middleItemLayout := container.New(layout.NewVBoxLayout(), emptyLabel, titleLabel, list)
	middleItemLayout.Resize(fyne.NewSize(100, 900))
	return list
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

func (d *DownApp) gui() {
	a := app.New()
	window := a.NewWindow("Downloader")
	window.Resize(fyne.NewSize(1050, 750))
	window.SetMaster()
	window.SetPadded(false)
	contentLayout := d.makeUI()
	d.mainWindow = &window
	window.SetContent(contentLayout)
	window.ShowAndRun()
}

func (d *DownApp) log(str string) {
	d.dataLog.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	d.dataLog.Println(str)
}

func (d *DownApp) err(str string) {
	d.errLog.Println(str)
}
