package snet

import (
	"bytes"
	"encoding/binary"
	"game_server_silk/siface"
	"game_server_silk/utils"
	"github.com/pkg/errors"
)

/*
	封包、拆包的具体实现模块
 */
type DataPack struct {}

//封包、拆包实例的一个初始化方法
func NewDataPack() *DataPack {
	return &DataPack{}
}

/*
	获取包的头长度方法
	固定包的消息头长度 Datalen uint32（4字节）+ ID uint32（4字节）
 */
func (dp *DataPack) GetHeadLen() uint32 {
	return 8
}

/*
	将消息数据进行封包，写到包里的顺序 |datalen|msgID|data|
 */
func (dp *DataPack) Pack(msg siface.IMessage) ([]byte, error) {
	//创建一个存放的bytes字节流的缓冲
	//bytes.NewBuffer是一个缓冲byte类型的缓冲器（容器）存放着都是byte
	dataBuff := bytes.NewBuffer([]byte{})

	//将dataLen数据内容长度，写进缓冲器 databuff中
	//binary.Write是将数据转换成二进制形式 和 小端字节序形式的，然后写进到buf缓冲器中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}

	//将MsgId 写进databuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}

	//将data数据 写进databuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	//返回二进制的包
	return dataBuff.Bytes(), nil
}

/*
	拆包（将包的Head读出来）然后再根据head里的data长度，在进行一次读
 */
func (dp *DataPack) Unpack(binaryData []byte) (siface.IMessage, error) {
	//传入二进制字节数据的，返回一个创建好的ioReader流对象
	dataBuff := bytes.NewReader(binaryData)

	//只读取包的头部head信息，得到datalen和MsgID
	//实例一个消息模块对象
	msg := &Message{}

	//读消息数据长度datalen，将databBuff包中的数据长度读取到消息模块DataLen属性中
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	//读消息数据类型MsgID
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	//判断datalen是否已经超出了我们允许的最大包长度
	if (utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize) {
		return nil, errors.New("消息数据太大了，超出了设置的数据包大小")
	}

	return msg, nil
}