package conf

import "log"

type LvMsg struct {
	User  string
	Email string
	Text  string
	Times string
	URL   string
}

type statsCfg struct {
	BaseConf
	ToCnt int
	IPs   []string
	Msgs  []LvMsg
}

var Stat statsCfg

func init() {
	log.Println("statsCfg init ENTER")
	Stat.BaseConf.conf = &Stat
	Stat.Msgs = make([]LvMsg, 0)
	Stat.IPs = make([]string, 0)
}
