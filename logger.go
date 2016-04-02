package lessgo

import (
	"net"
	"time"

	"github.com/lessgo/lessgo/logs"
	"github.com/lessgo/lessgo/logs/color"
)

// Logger returns a middleware that logs HTTP requests.
func Logger() MiddlewareFunc {
	return func(next Handler) Handler {
		return HandlerFunc(func(c Context) (err error) {
			logs.Warn("进入Logger")
			req := c.Request()
			res := c.Response()

			remoteAddr := req.RemoteAddress()
			if ip := req.Header().Get(XRealIP); ip != "" {
				remoteAddr = ip
			} else if ip = req.Header().Get(XForwardedFor); ip != "" {
				remoteAddr = ip
			} else {
				remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
			}

			start := time.Now()
			if err := next.Handle(c); err != nil {
				c.Error(err)
			}
			stop := time.Now()
			method := req.Method()
			path := req.URL().Path()
			if path == "" {
				path = "/"
			}
			size := res.Size()

			n := res.Status()
			code := color.Green(n)
			switch {
			case n >= 500:
				code = color.Red(n)
			case n >= 400:
				code = color.Yellow(n)
			case n >= 300:
				code = color.Cyan(n)
			}

			logs.Debug("%s | %s | %s | %s | %s | %d", remoteAddr, method, path, code, stop.Sub(start), size)
			return nil
		})
	}
}