package service

import (
	"log"
	"ztun/dao"
	"ztun/model"
)

// Service struct
type Service struct {
	dao     dao.Dao
	sel     *model.Tunnel
	pad     *Pad
	clients map[int64]*Client
}

// New init
func NewService(pad *Pad, d dao.Dao) (s *Service) {
	s = &Service{
		pad:     pad,
		dao:     d,
		sel:     new(model.Tunnel),
		clients: make(map[int64]*Client),
	}
	return s
}

func (svc *Service) SetPad(p *Pad) {

	svc.pad = p

}

func (svc *Service) GetTunList() (tuns []*model.Tunnel, err error) {
	return svc.dao.GetTunList()
}

func (svc *Service) SaveTun() (err error) {
	svc.sel.Name = svc.pad.GetNameText()
	svc.sel.Port = svc.pad.GetPortText()
	svc.sel.Remote = svc.pad.GetAddrText()

	err = svc.dao.SaveTun(svc.sel)
	if err != nil {
		return
	}
	svc.pad.RefreshTunBox(svc)
	return
}

func (svc *Service) DestroyTun() (err error) {
	svc.sel.Name = svc.pad.GetNameText()
	svc.sel.Port = svc.pad.GetPortText()
	svc.sel.Remote = svc.pad.GetAddrText()

	err = svc.dao.DestroyTun(svc.sel)
	if err != nil {
		return
	}
	svc.SelectTun(-1)
	svc.pad.RefreshTunBox(svc)
	return
}

func (svc *Service) SelectTun(id int64) (err error) {
	if svc.sel.ID == id {
		return
	}
	var tun *model.Tunnel
	//添加隧道
	if id == -1 {
		tun = new(model.Tunnel)
	} else {
		tun, err = svc.dao.GetTun(id)
		if err != nil {
			return
		}
	}
	//关闭原来客户端的日志
	client := svc.clients[svc.sel.ID]
	if client != nil && client.IsRunning() {
		client.DisableLog()
	}

	svc.sel = tun

	svc.pad.SetNameText(tun.Name)
	svc.pad.SetAddrText(tun.Remote)
	svc.pad.SetPortText(tun.Port)
	svc.pad.EnableConnect()

	svc.pad.SetLogText("")

	client = svc.clients[svc.sel.ID]
	if client != nil && client.IsRunning() {
		client.EnableLog(svc.GeLogFunc())
		svc.pad.SetConnectChecked(true)
	} else {
		svc.pad.SetConnectChecked(false)
	}

	return
}

func (svc *Service) Connect() (err error) {
	client := svc.clients[svc.sel.ID]
	if client != nil && client.IsRunning() {
		return
	}
	log.Printf("客户端>>>服务器,发起连接请求...")
	client, err = NewClient(svc.sel.Port, svc.sel.Remote)
	if err != nil {
		return
	}

	svc.clients[svc.sel.ID] = client

	go client.Start()
	client.EnableLog(svc.GeLogFunc())
	svc.pad.SetLogText("开始监听连接请求")
	return
}

func (svc *Service) Disconnect() (err error) {
	client := svc.clients[svc.sel.ID]
	if client == nil {
		return
	}
	client.Close()

	return
}
func (svc *Service) GeLogFunc() (info func(text string)) {
	info = func(text string) {
		svc.pad.SetLogText(text + "\n" + svc.pad.GetLogText())
	}
	return
}
