package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/cmd/fyne_demo/data"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/flopp/go-findfont"
	"net/url"
	"os"
	"strings"
	"ztun/dao"
	"ztun/dao/sqlite"
	"ztun/service"
)

const preferenceCurrentTab = "currentTab"

func init() {
	fontPaths := findfont.List()
	for _, path := range fontPaths {
		//楷体:simkai.ttf
		//黑体:simhei.ttf
		if strings.Contains(path, "STHeiti Light") {
			os.Setenv("FYNE_FONT", path)
			break
		}
	}

}

func parseURL(urlStr string) *url.URL {
	link, err := url.Parse(urlStr)
	if err != nil {
		fyne.LogError("Could not parse URL", err)
	}

	return link
}
func welcomeScreen(a fyne.App) fyne.CanvasObject {
	logo := canvas.NewImageFromResource(data.FyneScene)
	if fyne.CurrentDevice().IsMobile() {
		logo.SetMinSize(fyne.NewSize(171, 125))
	} else {
		logo.SetMinSize(fyne.NewSize(528, 367))
	}

	return widget.NewVBox(
		layout.NewSpacer(),
		widget.NewLabelWithStyle("欢迎使用WebSocket隧道工具", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewHBox(layout.NewSpacer(), logo, layout.NewSpacer()),

		widget.NewHBox(layout.NewSpacer(),
			widget.NewHyperlink("ztun", parseURL("https://github.com/zzpu/ztun")),
			widget.NewLabel("-"),
			widget.NewHyperlink("zserver", parseURL("https://github.com/zzpu/zserver")),
			layout.NewSpacer(),
		),
		layout.NewSpacer(),

		widget.NewGroup("Theme",
			fyne.NewContainerWithLayout(layout.NewGridLayout(2),
				widget.NewButton("Dark", func() {
					a.Settings().SetTheme(theme.DarkTheme())
				}),
				widget.NewButton("Light", func() {
					a.Settings().SetTheme(theme.LightTheme())
				}),
			),
		),
	)
}

func main() {
	a := app.NewWithID("io.fyne.demo")
	a.SetIcon(theme.FyneLogo())

	w := a.NewWindow("WebSocket隧道")

	db := sqlite.NewDB()
	d := dao.NewDao(db)
	pad := service.NewPad(w)
	svc := service.NewService(pad, d)

	tabs := widget.NewTabContainer(
		widget.NewTabItemWithIcon("欢迎", theme.HomeIcon(), welcomeScreen(a)),
		widget.NewTabItemWithIcon("隧道列表", theme.CheckButtonCheckedIcon(), pad.MakeSplitTab(svc)),
	)

	tabs.SetTabLocation(widget.TabLocationLeading)
	tabs.SelectTabIndex(a.Preferences().Int(preferenceCurrentTab))
	w.SetContent(tabs)

	w.ShowAndRun()
	a.Preferences().SetInt(preferenceCurrentTab, tabs.CurrentTabIndex())
}
