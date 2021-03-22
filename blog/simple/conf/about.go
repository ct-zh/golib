package conf

import "log"

type AbtCfg struct {
	BaseConf
	Name string
	Jobs string
	WX string
	QQ string
	Email string
}

var Abt AbtCfg

func init()  {
	log.Println("AbtCfg init ENTER")
	Abt.BaseConf.conf = &Abt
}