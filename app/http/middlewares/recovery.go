package middlewares

import (
	"gohub/pkg/logger"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				//获取用户请求
				httpRequest, _ := httputil.DumpRequest(c.Request, true)

				//链路终端，客户端终端连接为正常行为，不需要记录堆栈信息
				var brokenPipe bool

				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						errStr := strings.ToLower(se.Error())
						if strings.Contains(errStr, "broken pipe") || strings.Contains(errStr, "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				// 链接中断的情况
				if brokenPipe {
					logger.Error(c.Request.URL.Path,
						zap.Time("time", time.Now()),
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					c.Error(err.(error))
					c.Abort()
					// 链接已断开，无法写状态码
					return
				}

				// 如果不是链接中断，就开始记录堆栈信息
				logger.Error("recovery from panic",
					zap.Time("time", time.Now()),               // 记录时间
					zap.Any("error", err),                      // 记录错误信息
					zap.String("request", string(httpRequest)), // 请求信息
					zap.Stack("stacktrace"),                    // 调用堆栈信息
				)

				// 返回 500 状态码
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"message": "服务器内部错误，请稍后再试",
				})
			}
		}()
		c.Next()
	}
}
