package methods

import (
	"fmt"

	"github.com/spf13/viper"
)

type Infra struct {
	Method   string
	Region   string
	Network  NetworkConf
	Instance InstanceConf
	SSH      SSHConf
}

type InstanceConf struct {
	Hostname string
	Label    string
	OSId     int
	Plan     string
	Tag      string
}

type NetworkConf struct {
	Description string
}

type SSHConf struct {
	Name string
}

func ReadYaml(yamlFile string) (conf Infra, err error) {
	viper.SetConfigName(yamlFile)
	viper.AddConfigPath(".")
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		return Infra{}, fmt.Errorf("error reading config file, %s", err)
	}

	err = viper.Unmarshal(&conf)
	if err != nil {
		return Infra{}, fmt.Errorf("unable to decode into struct, %v", err)
	}
	return conf, nil

}
