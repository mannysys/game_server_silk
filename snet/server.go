package snet

import (
	"fmt"
	"game_server_silk/siface"
	"game_server_silk/utils"
	"net"
)

/*
	IServer的接口方法实现，定义一个Server的服务模块
*/
type Server struct {
	//服务器的名称
	Name string
	//服务器绑定的tcp网络协议版本
	IPVersion string
	//服务器监听的IP
	IP string
	//服务器监听的端口
	Port int

	//当前的server添加一个router属性，server注册的链接对应的处理业务
	Router siface.IRouter
}


//启动服务器
func (s *Server) Start() {
	fmt.Printf("[Silk] Server Name:%s, listenner at IP:%s, Port:%d is starting\n",
		utils.GlobalObject.Name,utils.GlobalObject.Host,utils.GlobalObject.TcpPort)
	fmt.Printf("[Silk] Version %s, MaxConn:%d, MaxPackageSize:%d\n",
		utils.GlobalObject.Version,utils.GlobalObject.MaxConn,utils.GlobalObject.MaxPackageSize)

	//开启一个协程去处理阻塞等待客户端连接业务（异步形式去处理客户端连接）
	go func() {

		//1 获取一个TCP的Addr（实例一个tcp协议自定义的地址和端口）
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d",s.IP,s.Port))
		if err != nil {
			fmt.Println("[实例化TCP服务端网络连接地址和端口<失败>]", err)
			return
		}
		//2 监听服务器的地址（在自定义实例的地址和端口上进行监听）
		listenner, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("[监听服务端的TCP地址和端口<失败>]：", s.IPVersion, err)
			return
		}
		fmt.Println("start server success", s.Name, " success Listenning...")

		//生成一个连接ID
		var cid uint32
		cid = 0

		//3 阻塞的等待客户端连接，处理客户端连接业务（读写）
		for {
			//如果有客户端连接过来，阻塞会返回一个tcp连接对象
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("[当前阻塞客户端连接<失败>]", err)
				continue  //只跳出当前循环
			}

			//将原生连接对象conn和自定义方法，交给我们封装的新连接模块去处理
			dealConn := NewConnection(conn, cid, s.Router)
			cid ++

			//启动当前连接业务处理
			go dealConn.Start()

		}

	}()

}
//停止服务器
func (s *Server) Stop() {
	//TODO 将一些服务器的资源、状态或者一些已经开辟的连接信息 进行停止或者回收

}
//运行服务器
func (s *Server) Serve() {
	//启动server服务功能
	s.Start()

	//TODO 做一些启动服务器之后的额外业务

	//阻塞状态（阻塞主线程为了处理执行异步）
	select {}
}


//路由功能：给当前的服务注册一个路由方法，供客户端的链接处理使用
func (s *Server) AddRouter(router siface.IRouter) {
	s.Router = router
	fmt.Println("路由注册添加成功！！")
}

/*
	初始化Server模块的方法
*/
func NewServer(name string) siface.IServer {
	s := &Server{
		Name:      utils.GlobalObject.Name,
		IPVersion: "tcp4",
		IP:        utils.GlobalObject.Host,
		Port:      utils.GlobalObject.TcpPort,
		Router:    nil,
	}
	return s
}


