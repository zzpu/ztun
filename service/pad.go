package service

import (
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/widget"
)

type Pad struct {
	win     fyne.Window
	tunBox  *widget.Box
	name    *widget.Entry
	addr    *widget.Entry
	port    *widget.Entry
	addTun  *widget.Button
	save    *widget.Button
	delete  *widget.Button
	log     *widget.Entry
	connect *widget.Check
}

func NewPad(win fyne.Window) (p *Pad) {
	p = &Pad{

		win: win,
	}
	return
}
func (p *Pad) SetConnectChecked(b bool) {
	p.connect.SetChecked(b)
}
func (p *Pad) EnableConnect() {
	p.connect.Enable()
}

func (p *Pad) DisableConnect() {
	p.connect.Disable()
}

func (p *Pad) SetNameText(text string) {
	p.name.SetText(text)
}
func (p *Pad) GetNameText() (text string) {
	return p.name.Text
}

func (p *Pad) SetAddrText(text string) {
	p.addr.SetText(text)
}
func (p *Pad) GetAddrText() (text string) {
	return p.addr.Text
}

func (p *Pad) SetPortText(text string) {
	p.port.SetText(text)
}
func (p *Pad) GetPortText() (text string) {
	return p.port.Text
}

func (p *Pad) SetLogText(text string) {
	p.log.SetText(text)
}
func (p *Pad) GetLogText() (text string) {
	return p.log.Text
}

func (p *Pad) MakeSplitTab(svc *Service) fyne.CanvasObject {
	left := p.MakeLeft(svc)

	right := p.makeRight(svc)
	tab := widget.NewHSplitContainer(widget.NewVScrollContainer(left), right)
	tab.Offset = 0.2

	return tab
}
func (p *Pad) makeTunList(svc *Service) (items []fyne.CanvasObject) {
	tuns, err := svc.GetTunList()
	if err != nil {
		return
	}
	for _, tun := range tuns {
		t := *tun
		items = append(items, widget.NewButton(fmt.Sprintf(t.Name), func() {
			svc.SelectTun(t.ID)
		}))
	}

	p.addTun = widget.NewButton("添加隧道", func() {
		svc.SelectTun(-1)
	})

	items = append(items, p.addTun)
	return items
}
func (p *Pad) MakeLeft(svc *Service) fyne.CanvasObject {

	vlist := p.makeTunList(svc)
	p.tunBox = widget.NewVBox(vlist...)

	return widget.NewVScrollContainer(p.tunBox)

}

func (p *Pad) RefreshTunBox(svc *Service) {

	vlist := p.makeTunList(svc)
	p.tunBox.Children = vlist
	p.tunBox.Refresh()
}

func (p *Pad) makeRight(svc *Service) fyne.Widget {
	p.name = widget.NewEntry()
	p.name.SetPlaceHolder("名称")
	//p.name.OnChanged = func(text string){
	//	p.svc.
	//}

	p.addr = widget.NewEntry()
	p.addr.SetPlaceHolder("远程WebSocket地址")

	p.port = widget.NewEntry()
	p.port.SetPlaceHolder("本地监听端口")

	p.log = widget.NewMultiLineEntry()
	p.log.Wrapping = fyne.TextWrapWord

	entryLoremIpsumScroller := widget.NewVScrollContainer(p.log)
	p.save = widget.NewButton("保存", func() {
		svc.SaveTun()
	})

	p.delete = widget.NewButton("删除", func() {
		svc.DestroyTun()
	})
	p.connect = widget.NewCheck("连接", func(on bool) {
		if on {
			err := svc.Connect()
			if err != nil {
				dialog.ShowError(err, p.win)
				p.connect.Checked = false
				return
			}
			p.save.Disable()

		} else {
			svc.Disconnect()
			p.save.Enable()

		}

	})
	p.connect.Disable()

	box := widget.NewVBox(
		p.name,
		p.addr,
		p.port,
		p.delete,
		p.save,
		p.connect,
	)

	right := widget.NewVSplitContainer(box, entryLoremIpsumScroller)
	right.Offset = 0.2

	return right
}
