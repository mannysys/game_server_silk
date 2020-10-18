package siface

/*
	IRequest 模块接口
	实际上是把客户端请求的连接信息模块，和请求的数据包封装到一个 Request模块中
 */
type IRequest interface {
	//得到当前连接模块
	GetConnection() IConnection

	//得到请求的消息数据
	GetData() []byte

	//获取到消息包类型ID
	GetMsgID() uint32
}