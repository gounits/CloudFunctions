package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gounits/CloudFunctions/tool"
)

// LoggingMiddleware 是一个HTTP中间件，用于记录请求日志
func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		info := fmt.Sprintf("方法=%s 路由=%s 耗时=%s", r.Method, r.URL.Path, duration)
		tool.Info(info)
	}
}
