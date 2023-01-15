package utils

import (
	"encoding/json"
	"os"
	"zinx/ziface"
)

/*
@author: shg
@since: 2023/1/16 2:14 AM
@mail: shgang97@163.com
*/

var Config *GlobalConfig

/*
用来初始化当前的全局配置对象
*/
func init() {
	// 全局默认配置
	Config = &GlobalConfig{
		Host:           "0.0.0.0",
		TcpPort:        8999,
		Name:           "ZinxServerApp",
		Version:        "V1.0",
		MaxConn:        1000,
		MaxPackageSize: 4096,
	}
	// 自定义全局配置
	Config.ReLoad()
}

// ReLoad 通过读取配置文件加载全局配置
func (config *GlobalConfig) ReLoad() {
	data, err := os.ReadFile("/home/zinx/test/conf/zinx.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &Config)
	if err != nil {
		panic(err)
	}
}

type GlobalConfig struct {
	/*
		Server
	*/
	TcpServer ziface.IServer // 当前 Zinx 全局的 Server 对象
	Host      string         // 当前服务器主机监听的IP
	TcpPort   int            // 当前服务器主机监听的端口号
	Name      string         // 当前服务器的名称

	/*
		Zinx
	*/
	Version        string // 当前 Zinx 版本号
	MaxConn        int    // 当前服务主机允许的最大连接数
	MaxPackageSize uint32 // 当前 Zinx 框架数据包的最大值
}
