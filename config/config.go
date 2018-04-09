package config

import (
	"crypto/tls"
	"flag"
	"os"
	"strings"
)

var Cfg *Config

type Config struct {
	Cli bool

	Http    string
	TlsHttp string
	Tcp     string

	DriverName     string
	DataSourceName string

	CertFile  string
	KeyFile   string
	TlsConfig *tls.Config

	StaticDir    string
	DataDir      string
	MpnConfigDir string
}

func New() *Config {
	c := &Config{}

	flag.BoolVar(&c.Cli, "cli", false, "cli >")

	/*
		flag.StringVar(&c.Http, "Http", ":8080", "http port")
		flag.StringVar(&c.TlsHttp, "TlsHttp", ":8000", "tls http port")
		flag.StringVar(&c.Tcp, "Tcp", ":6666", "tcp port")
		flag.StringVar(&c.DriverName, "DriverName", "mysql", "driver name")
		flag.StringVar(&c.DataSourceName, "DataSourceName", "chip:chip@tcp(127.0.0.1:3306)/chipdb", "data source name")
		flag.StringVar(&c.CertFile, "CertFile", "/home/iot/gopath/src/espressif.com/chip/factory/tls/factory.espressif.cn.crt", "cert file path")
		flag.StringVar(&c.KeyFile, "KeyFile", "/home/iot/gopath/src/espressif.com/chip/factory/tls/factory.espressif.cn.key", "key file path")
		flag.StringVar(&c.StaticDir, "StaticDir", "/home/iot/gopath/src/espressif.com/chip/factory/static", "static dir")
		flag.StringVar(&c.DataDir, "DataDir", "/home/iot/gopath/src/espressif.com/chip/factory/data", "data dir")
		flag.StringVar(&c.MpnConfigDir, "MpnConfigDir", "/home/iot/mpnCondif", "mpn config dir")
	*/

	flag.StringVar(&c.Http, "Http", ":8080", "http port")
	flag.StringVar(&c.TlsHttp, "TlsHttp", ":8000", "tls http port")
	flag.StringVar(&c.Tcp, "Tcp", ":6666", "tcp port")
	flag.StringVar(&c.DriverName, "DriverName", "mysql", "driver name")
	flag.StringVar(&c.DataSourceName, "DataSourceName", "root:hmj888888@tcp(127.0.0.1:3306)/chipdb", "data source name")
	path, _ := os.Getwd()
	path = path[:strings.Index(path, "factory")+8]
	flag.StringVar(&c.CertFile, "CertFile", path+"tls/factory.espressif.cn.crt", "cert file path")
	flag.StringVar(&c.KeyFile, "KeyFile", path+"tls/factory.espressif.cn.key", "key file path")
	flag.StringVar(&c.StaticDir, "StaticDir", path+"static", "static dir")
	flag.StringVar(&c.DataDir, "DataDir", path+"data", "data dir")
	flag.StringVar(&c.MpnConfigDir, "MpnConfigDir", "/home/hongmingjie/iot/mpnCondif", "mpn config dir")

	flag.Parse()

	if c.CertFile != "" && c.KeyFile != "" {
		var err error
		if c.TlsConfig, err = c.buildTLSConfig(); err != nil {
			panic("build TLS config error: " + err.Error())
		}
	}

	return c
}

func (c *Config) buildTLSConfig() (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(c.CertFile, c.KeyFile)
	if err != nil {
		return nil, err
	}
	TlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	return TlsConfig, nil
}

func init() {
	Cfg = New()
}
