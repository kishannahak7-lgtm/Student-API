package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)


type HTTPServer struct{
	Address string `yaml:"address" env-required:"true"`
}

type Config struct{
	Env string `yaml:"env" env:"ENV" env-required:"true" env-default:"production"`
	Storagepath string `yaml:"storage_path" env-required:"true"`
	HTTPServer `yaml:"http_server"`
}


func MustLoad() *Config{
	var configpath string

	configpath = os.Getenv("CONFIG_PATH")

	if configpath == ""{
		flags:= flag.String("config","","path to the cofiguration file")

		flag.Parse()

		configpath= *flags

		if configpath == ""{
			log.Fatal("config pathis not set")

		}
	}
	if _,err:=os.Stat(configpath); os.IsNotExist(err){
		log.Fatal("config file is does not exist :%s", configpath)
	}

	var cnfg Config

	err := cleanenv.ReadConfig(configpath,&cnfg)
	if err != nil{
		log.Fatal("can not read config file :%s",err.Error())
	}
	return &cnfg

}