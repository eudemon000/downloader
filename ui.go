package main

import (
	"downloader/load"
	"downloader/log"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var l *log.Log = log.New()

func (d *DownApp) makeUI() fyne.CanvasObject {
	top := d.createTop()
	left := d.createLeft()
	middle := d.createMiddle()
	//contentLayout := container.New(layout.NewBorderLayout(top, nil, left, right), top, nil, left, right)
	contentLayout := container.NewBorder(top, nil, left, nil, top, middle)
	return contentLayout
}

func (d *DownApp) createTop() fyne.CanvasObject {
	newTaskItme := widget.NewToolbarAction(resourceResourceAddPng, d.newDialog())
	t := widget.NewToolbar(widget.NewToolbarSpacer(), newTaskItme)
	topBackground := canvas.NewRectangle(color.RGBA{51, 51, 51, 255})
	topBackground.SetMinSize(fyne.Size{Width: 1050, Height: 50})
	toolbarLayout := container.New(layout.NewPaddedLayout(), t)
	topContent := container.New(layout.NewMaxLayout(), topBackground, layout.NewSpacer(), toolbarLayout)
	return topContent
}

func (d *DownApp) createLeft() fyne.CanvasObject {
	leftBackground := canvas.NewRectangle(color.RGBA{244, 245, 247, 255})
	leftBackground.SetMinSize(fyne.Size{Width: 200, Height: 0})
	l := widget.NewList(func() int {
		return len(middleItemText)
	}, func() fyne.CanvasObject {
		return widget.NewLabel("")
	}, func(lii widget.ListItemID, co fyne.CanvasObject) {
		co.(*widget.Label).SetText(middleItemText[lii])
		l.PrintMulti("aaa")
	})
	abc := container.NewMax(l)
	lb := container.New(layout.NewMaxLayout(), leftBackground, abc)
	return lb
}

func (d *DownApp) createMiddle() fyne.CanvasObject {
	return widget.NewLabel("middle")
}

func (d *DownApp) newDialog() func() {
	pathEntry := widget.NewEntry()
	items := []*widget.FormItem{
		widget.NewFormItem("Path", pathEntry),
	}
	return func() {
		dl := dialog.NewForm("New",
			"Submit",
			"Cancel",
			items,
			func(b bool) {
				if b {
					url := pathEntry.Text
					load.Header(url)
				}
			},
			*d.mainWindow)
		dl.Resize(fyne.Size{Width: 600, Height: 100})
		d.newTaskDialog = &dl
		dl.Show()
	}
}
