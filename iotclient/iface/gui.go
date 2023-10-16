package iface

import (
	"context"
	_ "embed"
	"main/core"

	webview "github.com/webview/webview_go"
)

//go:embed gui.html
var html string

type GUIDisplay struct {
	view      webview.WebView
	responses chan core.Response
}

func NewGUIDisplay() *GUIDisplay {
	return &GUIDisplay{
		view:      webview.New(false),
		responses: make(chan core.Response, 256),
	}
}

func (d *GUIDisplay) Close() error {
	d.view.Destroy()
	return nil
}

func (d *GUIDisplay) Display(ctx context.Context, res core.Response) error {
	go d.view.Dispatch(func() {
		d.view.Eval("document.dispatchEvent(new Event('response'))")
	})
	d.responses <- res
	return nil
}

func (d *GUIDisplay) Start(ctx context.Context) {
	d.view.SetTitle("HAL 9000")
	d.view.SetSize(400, 800, webview.HintNone)
	d.view.Bind("getResponses", func() core.Response {
		return <-d.responses
	})
	d.view.SetHtml(html)
	d.view.Run()
}
