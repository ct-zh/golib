package conf

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
)

// config接口
type IConf interface {
	Save()
	Load()
}

type BaseConf struct {
	conf IConf
}

func (cfg *BaseConf) Save() {
	log.Print("Save ENTER")
	t := reflect.TypeOf(cfg.conf).Elem().Name()
	fName := fmt.Sprintf("%s.json", t)
	log.Println(fName)

	jsonD, err := json.Marshal(cfg.conf)
	if err != nil {
		log.Println("生成错误")
	}

	SaveFile(fName, jsonD)
	log.Println(cfg.conf)
}

func (cfg *BaseConf) Load() {
	log.Println("Load ENTER")
	t := reflect.TypeOf(cfg.conf).Elem().Name()
	fName := fmt.Sprintf("%s.json", t)
	log.Println(fName)
	data, err := ReadFile(fName)
	if err != nil {
		log.Println(err)
		return
	}
	err = json.Unmarshal(data, cfg.conf)
	if err != nil {
		log.Panicln(err)
	}
	log.Println(cfg.conf)
}

func SaveFile(fileName string, data []byte) error {
	fi, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer fi.Close()
	w := bufio.NewWriter(fi)
	w.WriteString(string(data))
	w.Flush()
	return nil
}

func ReadFile(fileName string) ([]byte, error) {
	fi, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer fi.Close()

	r := bufio.NewReader(fi)
	var data []byte
	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			log.Println(err)
			return nil, err
		}

		if 0 == n {
			break
		}
		data = append(data, buf[:n]...) // ?
	}
	return data, err
}
