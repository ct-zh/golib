package conf

import "log"

type ItemCfg struct {
	BaseConf
	Items []string
}

var Item ItemCfg

func init() {
	log.Println("ItemCfg init ENTER")
	Item.BaseConf.conf = &Item
	Item.Items = make([]string, 0)
}
