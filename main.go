package main

import (
	"flag"
	"fmt"
	"gohub/bootstrap"

	btsConfig "gohub/config"

	"gohub/pkg/config"

	"github.com/gin-gonic/gin"
)

func init() {
	btsConfig.Initialize()
}

func main() {

	var env string
	flag.StringVar(&env, "env", "", "加载 .env 文件, 如 --env=testing 加载的是 .env.testing 文件")
	flag.Parse()
	config.InitConfig(env)

	config.InitConfig(env)

	// 初始化 Logger
	bootstrap.SetupLogger()

	// 设置 gin 的运行模式，支持 debug, release, test
	// release 会屏蔽调试信息，官方建议生产环境中使用
	// 非 release 模式 gin 终端打印太多信息，干扰到我们程序中的 Log
	// 故此设置为 release，有特殊情况手动改为 debug 即可
	gin.SetMode(gin.ReleaseMode)

	// new 一个 Gin Engine 实例
	router := gin.New()

	bootstrap.SetupDB()

	bootstrap.SetupRedis()
	// 初始化路由绑定r
	bootstrap.SetupRoute(router)

	// sms.NewSMS().Send("13552033797", sms.Message{
	// 	Template: config.GetString("sms.aliyun.template_code"),
	// 	Data:     map[string]string{"code": "888888"},
	// })

	// verifycode.NewVerifyCode().SendSMS("13552033797")

	// logger.Dump(captcha.NewCaptcha().VerifyCaptcha("J91gZOjvac0ZVMGdsGdy", "296872"))
	// logger.Dump(captcha.NewCaptcha().VerifyCaptcha("J91gZOjvac0ZVMGdsGdy", "000000"))
	// 运行服务
	err := router.Run(":" + config.Get("app.port"))
	if err != nil {
		// 错误处理，端口被占用了或者其他错误
		fmt.Println(err.Error())
	}
}
