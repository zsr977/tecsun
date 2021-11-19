package helper

func RegisterHelper() {
	// 注册配置服务
	var c ConfigCenter
	c.Init("runtime")
}
