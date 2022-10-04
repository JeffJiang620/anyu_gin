package recover

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/JeffJiang620/anyu_gin/response"
	"github.com/JeffJiang620/anyu_logger/loggers"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net"
	"net/http/httputil"
	"os"
	"runtime"
	"strings"
	"time"
)

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

const (
	reset = "\033[0m"
)

func source(lines [][]byte, n int) []byte {
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())

	if lastSlash := bytes.LastIndex(name, slash); lastSlash >= 0 {
		name = name[lastSlash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}

func stack(skip int) []byte {
	buf := new(bytes.Buffer)
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		_, _ = fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := os.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		_, _ = fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}
	return buf.Bytes()
}

func timeFormat(t time.Time) string {
	var timeString = t.Format("2006-01-02 15:04:05")
	return timeString
}

func checkBrokenPipe(err interface{}) bool {
	var brokenPipe bool
	if ne, ok := err.(*net.OpError); ok {
		if se, ok := ne.Err.(*os.SyscallError); ok {
			if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
				brokenPipe = true
			}
		}
	}
	return brokenPipe
}

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var brokenPipe = checkBrokenPipe(err)
				var errLog string

				if brokenPipe {
					httpRequest, _ := httputil.DumpRequest(c.Request, false)
					errLog = fmt.Sprintf("%s\n%s%s", err, string(httpRequest), reset)
				} else if gin.IsDebugging() {
					httpRequest, _ := httputil.DumpRequest(c.Request, false)
					headers := strings.Split(string(httpRequest), "\r\n")
					for idx, header := range headers {
						current := strings.Split(header, ":")
						if current[0] == "Authorization" {
							headers[idx] = current[0] + ": *"
						}
					}
					stack := stack(3)
					errLog = fmt.Sprintf("[Recovery] %s panic recovered:\n%s\n%s\n%s%s",
						timeFormat(time.Now()), strings.Join(headers, "\r\n"), err, stack, reset)
				} else {
					stack := stack(3)
					errLog = fmt.Sprintf("[Recovery] %s panic recovered:\n%s\n%s%s",
						timeFormat(time.Now()), err, stack, reset)
				}

				if brokenPipe {
					_ = c.Error(err.(error))
				} else if ve, ok := err.(error); ok {
					er := response.UnknownError.WithErr(ve)
					er.Render(c)
				} else if s, ok := err.(string); ok {
					er := response.UnknownError.WithErr(errors.New(s))
					er.Render(c)
				} else {
					response.UnknownError.Render(c)
				}

				path := c.Request.URL.Path
				ip := c.ClientIP()
				method := c.Request.Method
				loggers.Logger.Name("request_panic").Error("",
					zap.String("ip", ip),
					zap.String("method", method),
					zap.String("path", path),
					zap.String("errMsg", errLog))
				c.Abort()
			}
		}()
		c.Next()
	}
}
