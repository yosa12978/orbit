package config

import (
	"os"
	"sync"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
)

type Config struct {
	App struct {
		Name string `yaml:"name"`
	} `yaml:"app"`
	Server struct {
		InstanceID uint16 `yaml:"instance_id"`
		Host       string `yaml:"hostname"`
		Addr       string `yaml:"addr"`
	} `yaml:"server"`
	Redis struct {
		Addr     string `yaml:"addr"`
		DB       int    `yaml:"db"`
		Password string `yaml:"password"`
	} `yaml:"redis"`
}

var (
	conf     Config
	initConf sync.Once
)

func Get() Config {
	initConf.Do(func() {
		file, err := os.Open("config.yml")
		if err != nil {
			panic(err)
		}
		defer file.Close()
		if err := yaml.NewDecoder(file).Decode(&conf); err != nil {
			panic(err)
		}
		envconfig.MustProcess("orbit", &conf)
	})
	return conf
}
