package tool

import "log"

// ANSI颜色代码
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
)

// 初始化日志输出配置
func init() {
	// 设置日志输出格式：包含时间、文件名和行号
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

// Success 记录成功日志（绿色）
func Success(message string) {
	log.SetPrefix(colorGreen + "[成功] " + colorReset)
	log.Println(colorGreen + message + colorReset)
}

// Error 记录错误日志（红色）
func Error(message string) {
	log.SetPrefix(colorRed + "[错误] " + colorReset)
	log.Println(colorRed + message + colorReset)
}

// Warn 记录警告日志（黄色）
func Warn(message string) {
	log.SetPrefix(colorYellow + "[警告] " + colorReset)
	log.Println(colorYellow + message + colorReset)
}

// Info 记录信息日志（蓝色）
func Info(message string) {
	log.SetPrefix(colorBlue + "[运行] " + colorReset)
	log.Println(colorBlue + message + colorReset)
}
