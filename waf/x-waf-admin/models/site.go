package models

import (
	"log"
	"time"
)

// debuglevel: debug, info, notice, warn, error, crit, alert, emerg
// ssl: on, off
type Site struct {
	Id          int64
	SiteName    string `xorm:"unique"`
	Port        int
	BackendAddr []string
	UnrealAddr  []string
	Ssl         string    `xorm:"varchar(10) notnull default 'off'"`
	DebugLevel  string    `xorm:"varchar(10) notnull default 'error'"`
	LastChange  time.Time `xorm:"updated"`
	Version     int       `xorm:"version"` // 乐观锁
}

func ListSite() (sites []Site, err error) {
	sites = make([]Site, 0)
	err = Engine.Find(&sites)
	log.Println(err, sites)
	return sites, err
}

func ListSiteById(Id int64) (sites []Site, err error) {
	sites = make([]Site, 0)
	err = Engine.Id(Id).Find(&sites)
	log.Println(err, sites)
	return sites, err
}

func NewSite(siteName string, Port int, BackendAddr []string, UnrealAddr []string, SSL string, DebugLevel string) (err error) {
	if SSL == "" {
		SSL = "off"
	}
	if DebugLevel == "" {
		DebugLevel = "error"
	}

	_, err = Engine.Insert(&Site{SiteName: siteName, Port: Port, BackendAddr: BackendAddr, UnrealAddr: UnrealAddr, Ssl: SSL, DebugLevel: DebugLevel})
	return err
}

func UpdateSite(Id int64, SiteName string, Port int, BackendAddr []string, UnrealAddr []string, SSL string, DebugLevel string) (err error) {
	if SSL == "" {
		SSL = "off"
	}
	if DebugLevel == "" {
		DebugLevel = "error"
	}

	site := new(Site)
	Engine.Id(Id).Get(site)
	site.SiteName = SiteName
	site.Port = Port
	site.BackendAddr = BackendAddr
	site.UnrealAddr = UnrealAddr
	site.Ssl = SSL
	site.DebugLevel = DebugLevel
	_, err = Engine.Id(Id).Update(site)
	return err
}

func DelSite(id int64) (err error) {
	_, err = Engine.Delete(&Site{Id: id})
	return err
}
