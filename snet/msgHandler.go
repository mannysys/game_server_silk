package snet

import (
	"fmt"
	"game_server_silk/siface"
	"strconv"
)

/*
	消息管理处理模块的实现
 */
type MsgHandler struct {
	//存放每个MsgID 所对应的处理方法
	Apis map[uint32] siface.IRouter
}

//创建初始化MsgHandler方法
func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis: make(map[uint32] siface.IRouter),
	}
}

//调度执行对应的Router的消息处理方法
func (mh *MsgHandler) DoMsgHandler(request siface.IRequest) {
	//1 从Request找到msgID
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("调度执行的路由API无法找到", request.GetMsgID(), ", 你需要注册该路由")
	}

	//2 根据MsgID 调度对应router业务即可
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

//为消息添加具体的处理逻辑
func (mh *MsgHandler) AddRouter(msgID uint32, router siface.IRouter) {
	//1 判断 当前msg绑定的API处理方法是否已经存在（检查map集合中是否已经存在msgID对应的路由）
	if _, ok := mh.Apis[msgID]; ok {
		//msgID绑定的路由，已经注册了
		panic("repeat api, msgID = " + strconv.Itoa(int(msgID)))
	}

	//2 添加msg与API绑定关系（添加新msgID绑定的路由处理业务）
	mh.Apis[msgID] = router
	fmt.Println("路由添加到map集合中成功！", msgID)

}