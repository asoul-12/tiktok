package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type Config struct {
	Mysql  Mysql  `json:"mysql"`
	Server Server `json:"server"`
	JWT    JWT    `json:"JWT"`
	Cos    Cos    `json:"cos"`
}

type Cos struct {
	Addr      string
	SecretID  string
	SecretKey string
	CdnAddr   string
}
type JWT struct {
	Issuer              string        `json:"issuer,omitempty"`
	TokenExpireDuration time.Duration `json:"tokenExpireDuration,omitempty"`
	Secrete             string        `json:"secrete,omitempty"`
	Subject             string        `json:"subject"`
}
type Mysql struct {
	Addr     string `json:"addr,omitempty"`
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
	Database string `json:"database,omitempty"`
}

type Server struct {
	Addr         string        `json:"addr,omitempty"`
	ReadTimeOut  time.Duration `json:"readTimeOut,omitempty"`
	WriteTimeOut time.Duration `json:"writeTimeOut,omitempty"`
}

func LoadConfigFromFile(file string) (config *Config) {
	bytes, err := os.ReadFile(file)
	if err != nil {
		log.Fatalf("从文件加载配置出错了,%s", err)
	}
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		log.Fatalf("文件反序列化失败,%s", err)
	}
	return config
}

func (mysql *Mysql) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", mysql.User, mysql.Password, mysql.Addr, mysql.Database)
}
