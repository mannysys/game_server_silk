package utils

import (
	"encoding/json"
	"game_server_silk/siface"
	"io/ioutil"
)

/*
	存储一切有关silk框架的全局参数，供其它模块使用
	一些参数是可以通过框架服务端用户进行配置
 */
type GlobalObj struct {
	/*
		Server
	 */
	TcpServer siface.IServer	//当前框架全局的server模块对象
	Host string					//当前服务器主机监听的IP
	TcpPort int					//当前服务器主机监听的端口号
	Name string					//当前自定义的服务器名称

	/*
		silk框架
	 */
	Version string			//当前框架的版本号
	MaxConn int 			//当前服务器主机允许的最大连接数
	MaxPackageSize uint32	//当前框架限制的数据包的最大值
}
/*
	定义一个全局的对外Globalobj
 */
var GlobalObject *GlobalObj

/*
	从 silk.json 去加载用于自定义的参数
 */
func (g *GlobalObj) Reload()  {
	data, err := ioutil.ReadFile("conf/silk.json")
	if err != nil {
		panic(err)
	}
	//将读取json文件数据映射解析到struct 属性中
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

/*
	提供一个init方法，初始化当前的GlobalObject
 */
func init() {
	//如果没有加载json配置文件，就是要以下初始化默认值
	GlobalObject = &GlobalObj{
		Name:           "SilkServerApp",
		Version:        "V0.5",
		TcpPort:		8999,
		Host:			"0.0.0.0",
		MaxConn:		1000,
		MaxPackageSize:	4096,
	}

	//应该尝试从conf/silk.json 去加载一些用户自定义的参数
	GlobalObject.Reload()
}