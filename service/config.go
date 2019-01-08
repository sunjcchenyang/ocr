package service

// 配置文件模块

import (
	"io/ioutil"
	"log"

	"go.uber.org/zap"

	"github.com/naoina/toml"
	"github.com/naoina/toml/ast"
)

// CommonConfig Common
type CommonConfig struct {
	Version  string
	IsDebug  bool
	LogLevel string
	LogPath  string

}

// ServerConf echo config struct
type ServerConf struct {
	Addr string
	Dir string
}
// libpython
type PythonConf struct {
	LibDir string
}

// Config ...
type Config struct {
	Common  *CommonConfig
	ServerC *ServerConf
	//add for othre config
	PythonC *PythonConf
}

// Conf ...
var Conf = &Config{}

// LoadConfig ...
func LoadConfig() {
	// init the new config params
	initConf()

	contents, err := ioutil.ReadFile("ocr.toml")
	if err != nil {
		log.Fatal("[FATAL] load kline.toml: ", err)
	}
	tbl, err := toml.Parse(contents)
	if err != nil {
		log.Fatal("[FATAL] parse kline.toml: ", err)
	}
	// parse common config
	parseCommon(tbl)
	// init log
	InitLogger()

	// parse Echo config
	parseServer(tbl)

	//parse the other config
	parsePython(tbl)

	Logger.Info("LoadConfig", zap.Any("Config", Conf))
}

func initConf() {
	Conf = &Config{
		Common:  &CommonConfig{},
		ServerC: &ServerConf{},
		//add for other config
		PythonC: &PythonConf{},
	}
}

func parseCommon(tbl *ast.Table) {
	if val, ok := tbl.Fields["common"]; ok {
		subTbl, ok := val.(*ast.Table)
		if !ok {
			log.Fatalln("[FATAL] : ", subTbl)
		}

		err := toml.UnmarshalTable(subTbl, Conf.Common)
		if err != nil {
			log.Fatalln("[FATAL] parseCommon: ", err, subTbl)
		}
	}
}

func parseServer(tbl *ast.Table) {
	if val, ok := tbl.Fields["server"]; ok {
		subTbl, ok := val.(*ast.Table)
		if !ok {
			log.Fatalln("[FATAL] : ", subTbl)
		}

		err := toml.UnmarshalTable(subTbl, Conf.ServerC)
		if err != nil {
			log.Fatalln("[FATAL] parseServer: ", err, subTbl)
		}
	}
}
//parse python config
func parsePython(tbl *ast.Table) {
	if val ,ok := tbl.Fields["python"]; ok {
		subTbl,ok := val.(*ast.Table)
		if !ok {
			log.Fatalln("[FATAL]: ",subTbl)
		}
		err := toml.UnmarshalTable(subTbl,Conf.PythonC)
		if err !=nil {
			log.Fatalln("[FATAL] parsePython:", err,subTbl)
		}
	}

}
