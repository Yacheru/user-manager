package logger

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

const (
	Red    = "41"
	Yellow = "43"
	Green  = "42"
)

func HTTPLogger(param gin.LogFormatterParams) string {
	var color string
	switch {
	case param.StatusCode >= 500:
		color = Red
	case param.StatusCode >= 400:
		color = Yellow
	default:
		color = Green
	}

	return fmt.Sprintf("[HTTP] [%s] \033[%sm %d \033[0m %s | %s | %s | %s | %s\n",
		param.TimeStamp.Format("2006-01-02 15:04:05"),
		color,
		param.StatusCode,
		param.Method,
		param.Path,
		param.Latency.String(),
		param.ClientIP,
		param.Request.UserAgent(),
	)
}
