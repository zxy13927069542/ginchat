package config

import "github.com/spf13/viper"
//import log "github.com/pion/ion-log"

//	Name: edge-cygateway-api
//
// Host: 0.0.0.0
// Port: 8801
//
// # mysql配置信息
// Mysql:
// Dns: root:066311@tcp(127.0.0.1:3306)/ginchat?charset=utf8mb4&parseTime=True&loc=Local
type Config struct {
	Name  string //	服务名称
	Host  string //	IP
	Port  int    //	端口

	//	Mysql配置
	Mysql struct {
		Dns string
	}

	//Redis:
	//  Host: 0.0.0.0
	//  Port: 6379
	//  Password:
	//  DB: 0
	//  Poolsize: 30
	//  MinIdleConn: 30
	Redis struct {
		Host string
		Port string
		Password string
		DB int
		PoolSize int
		MinIdleConn int
	}
}

var c = &Config{}

//	Init() 读取配置文件
func Init(path string) Config {
	viper.SetConfigName("app")
	viper.AddConfigPath(path)
	err := viper.ReadInConfig()
	if err != nil {
		//log.Errorf("Read config failed! Err: [%v]", err)
		panic(err)
	}

	if err = viper.Unmarshal(c); err != nil {
		//log.Errorf("Parse config file failed! Err: [%v]", err)
		panic(err)
	}

	return *c
}
