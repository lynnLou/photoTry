/*
   @Time : 2020/5/6 13:59 下午
   @Author : Rebeta
   @Email : master@rebeta.cn
   @File : middleware
   @Software: GoLand
*/

package middleware

import (
	"github.com/offcn-jl/cscf"
	"github.com/offcn-jl/cscf/fake-http"
	"github.com/offcn-jl/go-common"
	"github.com/offcn-jl/go-common/codes"
	"github.com/offcn-jl/go-common/configer"
	"serverless/common/config"
	"strings"
)

// 向响应头添加版本信息
func AddVersions(apiVersion string) chaos.HandlerFunc {
	return func(c *chaos.Context) {
		c.Header("X-CSCF-Version", chaos.Version)
		c.Header("X-Common-Version", common.Version)
		c.Header("X-"+config.Project+"-Version", config.Version)
		c.Header("X-"+config.Project+"-Api-Version", apiVersion)
		c.Next()
	}
}

// 跨域检查
func CheckOrigin() chaos.HandlerFunc {
	return func(c *chaos.Context) {
		// 跨域校验
		allowOrigins := configer.GetString("AllowOrigins", "")
		allowOriginsArray := strings.Split(allowOrigins, ",")
		pass := false
		for _, origin := range allowOriginsArray {
			// 遍历配置中的跨域头，寻找匹配项
			if c.GetHeader("origin") == origin {
				c.Header("Access-Control-Allow-Origin", origin)
				pass = true
				// 只要有一个跨域头匹配就跳出循环
				break
			}
		}

		if !pass {
			c.JSON(http.StatusForbidden, chaos.H{"Code": codes.NotCertifiedCORS, "Error": codes.ErrorText(codes.NotCertifiedCORS)})
			c.Abort() // 出错后结束请求
		}

		// 通过跨域校验后，放行所有 OPTIONS 方法，并添加按照客户端的请求添加 Allow Headers
		//if c.Request.Method == "OPTIONS" {
		//	// 请求首部  Access-Control-Request-Headers 出现于 preflight request （预检请求）中，用于通知服务器在真正的请求中会采用哪些请求首部。
		//	c.Header("Access-Control-Allow-Headers", c.GetHeader("Access-Control-Request-Headers")) // 放行预检请求通知的请求首部。
		//	// https://cloud.tencent.com/developer/section/1189896
		//	c.Header("Access-Control-Allow-Methods", c.GetHeader("Access-Control-Request-Method")) // 放行预检请求通知的请求首部。
		//	c.AbortWithStatus(http.StatusNoContent)
		//}
	}
}
