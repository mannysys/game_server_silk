package snet

import (
	"game_server_silk/siface"
)

/*
	IServer的接口方法实现，定义一个Server的服务模块
*/
type Server struct {
	//服务器的名称
	Name string
	//服务器绑定的ip版本
	IPVersion string
	//服务器监听的IP
	IP string
	//服务器监听的端口
	Port int
}

//启动服务器
func (s *Server) Start()  {

}
//停止服务器
func (s *Server) Stop() {

}
//运行服务器
func (s *Server) Serve() {

}

func NewServer(name string) siface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
	}
	return s
}


