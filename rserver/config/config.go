package config

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

type Limit struct {
	Limit int `json:"limit"`
	Burst int `json:"burst"`
}

var (
	l     *Limit
	once  sync.Once
	cFile = "rserver/config/server.json"
)

func Load() *Limit {
	once.Do(func() {
		file, err := os.Open(cFile)
		if err != nil {
			log.Fatalln("打开配置文件出错", err)
		}

		defer func(file *os.File) {
			err := file.Close()
			if err != nil {

			}
		}(file)

		if err = json.NewDecoder(file).Decode(&l); err != nil {
			log.Fatalln("读取 json 文件出错", err)
		}
	})

	return l
}
