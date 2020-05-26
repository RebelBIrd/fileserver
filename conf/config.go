package conf

import (
	"fileserver/model"
	"io/ioutil"

	"github.com/qinyuanmao/go-utils/logutl"
	"github.com/qinyuanmao/go-utils/ormutl"
	"gopkg.in/yaml.v2"
)

type ApiConf struct {
	Url  string `yaml:"url"`
	Port int    `yaml:"port"`
}

type HttpKeyConfig struct {
	Key1 string `yaml:"key1"`
	Key2 string `yaml:"key2"`
}
type Config struct {
	ApiConf          `yaml:"api"`
	ormutl.MysqlConf `yaml:"mysql"`
	HttpKeyConfig    `yaml:"http_key"`
}

const yamlFilePath = "./conf/config.yaml"

var QsConfig *Config

func init() {
	if yamlFile, err := ioutil.ReadFile(yamlFilePath); err != nil {
		logutl.Error(err.Error())
	} else if err = yaml.Unmarshal(yamlFile, &QsConfig); err != nil {
		logutl.Error(err.Error())
	} else {
		ormutl.InitMysql(QsConfig.MysqlConf) // 初始化Mysql
		ormutl.GetEngine().InitTables(nil, model.FileNode{})
	}
}
