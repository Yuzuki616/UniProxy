package conf

import (
	"encoding/json"
	"os"
)

var C *Conf

type Conf struct {
	path string
	Log  Log `json:"Log"`
	Api  Api `json:"Api"`
}

type Log struct {
	Level string `json:"Level"`
}

type Api struct {
	Balance string   `json:"Balance"`
	Baseurl []string `json:"Baseurl"`
}

func Init(path string) error {
	C = New(path)
	if path == "" {
		return nil
	}
	err := C.Load()
	if err != nil {
		return err
	}
	return nil
}

func New(path string) *Conf {
	return &Conf{
		path: path,
		Log: Log{
			Level: "info",
		},
	}
}

func (c *Conf) Load() error {
	f, err := os.Open(c.path)
	if err != nil {
		return err
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(c)
	if err != nil {
		return err
	}
	return nil
}
