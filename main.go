package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/cesarFuhr/gocrypto/keys"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
)

func main() {
	run()
}

func run() {

	cfgSource := getCfgSource()
	cfg, err := loadConfigs(cfgSource)
	if err != nil {
		panic(err)
	}

	keySource := keys.NewPoolKeySource(cfg.App.KeySource.RSAKeySize, cfg.App.KeySource.PoolSize)
	keySource.WarmUp()

	sqlKeyRepo := keys.SQLKeyRepository{Cfg: keys.SQLConfigs{
		Host:     cfg.Db.Host,
		Port:     cfg.Db.Port,
		User:     cfg.Db.User,
		Password: cfg.Db.Password,
		Dbname:   cfg.Db.Dbname,
		Driver:   cfg.Db.Driver,
	}}
	sqlKeyRepo.Connect()

	keyStore := keys.KeyStore{Source: &keySource, Repo: &sqlKeyRepo}
	crypto := JWECrypto{}
	server := &KeyServer{&keyStore, &crypto}
	if err := http.ListenAndServe(":"+cfg.Server.Port, server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}

func getCfgSource() string {
	var cfgFromEnv bool
	flag.BoolVar(&cfgFromEnv, "e", false, "load config from environment")
	flag.Parse()
	if cfgFromEnv == true {
		return "env"
	}
	return "yaml"
}

func loadConfigs(t string) (config, error) {
	switch t {
	case "env":
		return loadFromENV()
	default:
		f, err := os.Open("config.yaml")
		if err != nil {
			return config{}, err
		}
		defer f.Close()
		return loadFromYAML(f)
	}
}

func loadFromENV() (config, error) {
	var cfg config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return config{}, err
	}
	return cfg, nil
}

func loadFromYAML(r io.Reader) (config, error) {
	var cfg config
	err := yaml.NewDecoder(r).Decode(&cfg)
	if err != nil {
		return config{}, nil
	}
	return cfg, nil
}

type config struct {
	Server struct {
		Port string `yaml:"port" envconfig:"SERVER_PORT"`
	} `yaml:"server"`
	Db struct {
		Host     string `yaml:"host" envconfig:"DB_HOST"`
		Port     int    `yaml:"port" envconfig:"DB_PORT"`
		User     string `yaml:"user" envconfig:"DB_USER"`
		Password string `yaml:"password" envconfig:"DB_PASSWORD"`
		Dbname   string `yaml:"dbname" envconfig:"DB_NAME"`
		Driver   string `yaml:"driver" envconfig:"DB_DRIVER"`
	} `yaml:"database"`
	App struct {
		KeySource struct {
			PoolSize   int `yaml:"poolsize" envconfig:"APP_KEYSOURCE_POOL_SIZE"`
			RSAKeySize int `yaml:"rsakeysize" envconfig:"APP_KEYSOURCE_RSAKEY_SIZE"`
		} `yaml:"keysource"`
	} `yaml:"app"`
}
