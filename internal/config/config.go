package config

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/ghodss/yaml"
	"github.com/zekroTJA/kirby/internal/database"
)

type Main struct {
	Discord *Discord              `json:"discord"`
	MongoDB *database.MongoConfig `json:"mongodb"`
}

type Discord struct {
	Token   string `json:"token"`
	Prefix  string `json:"prefix"`
	OwnerID string `json:"owner_id"`
}

func Open(loc string) (*Main, error) {
	data, err := ioutil.ReadFile(loc)
	if os.IsNotExist(err) {
		err = cretaeDefault(loc)
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	cfg := new(Main)
	err = yaml.Unmarshal(data, cfg)
	return cfg, err
}

func cretaeDefault(loc string) error {
	def := &Main{
		Discord: new(Discord),
		MongoDB: &database.MongoConfig{
			Host:     "localhost",
			Port:     "27017",
			Username: "kirby",
			AuthDB:   "kirby",
			DataDB:   "kirby",
		},
	}

	data, err := yaml.Marshal(def)

	basePath := path.Base(loc)
	if _, err = os.Stat(basePath); os.IsNotExist(err) {
		err = os.MkdirAll(basePath, 0750)
		if err != nil {
			return err
		}
	}
	err = ioutil.WriteFile(loc, data, 0750)
	return err
}
