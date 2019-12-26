package network

//主要是为了调用cli的用户输入方法
type Clier interface {
	ReceiveCMD()
}
