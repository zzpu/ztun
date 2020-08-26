package dao

import (
	"ztun/dao/sqlite"
	"ztun/model"
)

type (
	Dao interface {
		Close()
		GetTunList() (tuns []*model.Tunnel, err error)
		GetTun(id int64) (tun *model.Tunnel, err error)
		SaveTun(tun *model.Tunnel) (err error)
	}
)

// dao dao.
type daoEx struct {
	dbx *sqlite.DB //读库
}

func (d daoEx) GetTunList() (tuns []*model.Tunnel, err error) {
	err = d.dbx.DBWrite().Table("tunnels").Find(&tuns).Error
	return
}

func (d daoEx) GetTun(id int64) (tun *model.Tunnel, err error) {
	tun = new(model.Tunnel)
	err = d.dbx.DBWrite().Table("tunnels").Where("id=?", id).First(tun).Error
	return
}

func (d daoEx) SaveTun(tun *model.Tunnel) (err error) {
	err = d.dbx.DBWrite().Save(tun).Error
	return
}

func (d daoEx) Close() {
	panic("implement me")
}

// New new a dao and retu
func NewDao(dbx *sqlite.DB) (d Dao) {
	d = &daoEx{
		dbx: dbx,
	}
	return
}
