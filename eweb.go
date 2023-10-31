package eweb

import (
	stdContext "context"
	"crypto/tls"
	"errors"
	"fmt"
	stdLog "log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"sync"
	"time"

	"github.com/gwaylib/eweb/jsonp"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/color"
)

var (
	// Global instance
	defaultE     *Eweb
	defaultELock = sync.Mutex{}
)

func DebugMode() bool {
	return defaultE.Debug
}

// Struct for rendering map
type H jsonp.Params

type Eweb struct {
	*echo.Echo
}

func New() *Eweb {
	e := &Eweb{
		Echo: echo.New(),
	}

	// fix log to stderr
	//	defaultE.Logger.SetOutput(os.Stderr)
	//	defaultE.Logger.SetPrefix("eweb")
	e.Server.Handler = defaultE
	e.TLSServer.Handler = defaultE

	// monitor middleware
	e.Pre(defaultE.FilterHandle)

	// e.Use(middleware.Recover())
	return e
}

// Using global instance to manager router packages
func Default() *Eweb {
	defaultELock.Lock()
	defer defaultELock.Unlock()
	if defaultE == nil {
		defaultE = New()
	}
	return defaultE
}

func (e *Eweb) colorForStatus(code string) string {
	switch {
	case code >= "200" && code < "300":
		return color.Green(code)
	case code >= "300" && code < "400":
		return color.White(code)
	case code >= "400" && code < "500":
		return color.Yellow(code)
	default:
		return color.Red(code)
	}
}

func (e *Eweb) colorForMethod(method string) string {
	switch method {
	case "GET":
		return color.Blue(method)
	case "POST":
		return color.Cyan(method)
	case "PUT":
		return color.Yellow(method)
	case "DELETE":
		return color.Red(method)
	case "PATCH":
		return color.Green(method)
	case "HEAD":
		return color.Magenta(method)
	case "OPTIONS":
		return color.White(method)
	default:
		return color.Reset(method)
	}
}

// Support http.FileSystem
// create http.FileSystem with http.FS(fs) for for embed.FS, os.File
func (e *Eweb) StaticFS(uriPrefix, filePrefix string, root http.FileSystem) *echo.Route {
	if root == nil {
		panic("unsupport nil root")
	}

	h := func(c echo.Context) error {
		p, err := url.PathUnescape(c.Param("*"))
		if err != nil {
			return err
		}
		file := filepath.Join(filePrefix, path.Clean("/"+p))

		f, err := root.Open(file)
		if err != nil {
			return echo.NotFoundHandler(c)
		}
		defer f.Close()

		fi, _ := f.Stat()
		if fi.IsDir() {
			file = filepath.Join(file, "index.html")
			f, err = os.Open(file)
			if err != nil {
				return echo.NotFoundHandler(c)
			}
			defer f.Close()
			if fi, err = f.Stat(); err != nil {
				return err
			}
		}
		http.ServeContent(c.Response(), c.Request(), fi.Name(), fi.ModTime(), f)
		return nil
	}

	e.GET(uriPrefix, h)
	if uriPrefix == "/" {
		return e.GET(uriPrefix+"*", h)
	}

	return e.GET(uriPrefix+"/*", h)
}

func (e *Eweb) FilterHandle(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		defer func() {
			resp := c.Response()
			req := c.Request()

			// 回传数据
			req.Header.Add("resp-size", strconv.FormatInt(resp.Size, 10))
			req.Header.Add("resp-status", strconv.Itoa(resp.Status))
		}()

		if err := next(c); err != nil {
			// deal default error
			var (
				code = http.StatusInternalServerError
				msg  string
			)

			if he, ok := err.(*echo.HTTPError); ok {
				code = he.Code
				msg = fmt.Sprintf("%v", he.Message)
			} else if e != nil && e.Debug {
				msg = err.Error()
			} else {
				msg = http.StatusText(code)
				if e == nil {
					fmt.Println(err)
				}
			}
			// Send response
			if !c.Response().Committed {
				if c.Request().Method == echo.HEAD {
					return c.NoContent(code)
				}
				return c.String(code, msg)
			}
			return err
		}
		return nil
	}
}

func (e *Eweb) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// request performance
	defer func(start time.Time) {
		pErr := recover()

		if pErr != nil {
			r.Header.Set("resp-status", "500")
		}

		req := r
		n := r.Header.Get("resp-status")
		if !e.Debug && n < "400" {
			return
		}

		ip := ReadIp(r)
		stop := time.Now()
		color.Printf(
			"[eweb] %s | %s | %-4s | USE:%12s | FROM:%15s | URI:%s \n",
			start.Format("2006-01-02 15:04:05"),
			e.colorForStatus(n), req.Method,
			// fmt.Sprintf("%.3fs", float64(stop.Sub(start)/1e6)/1000),
			stop.Sub(start).String(),
			ip, req.RequestURI,
		)

		if pErr != nil {
			println(fmt.Sprintf("panic %+v\n", pErr))
			debug.PrintStack()
		}
	}(time.Now())

	e.Echo.ServeHTTP(w, r)
}

func (e *Eweb) Close() error {
	if err := e.Echo.Close(); err != nil {
		return err
	}
	return nil
}
func (e *Eweb) Shutdown(ctx stdContext.Context) error {
	if err := e.Echo.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}

// Start starts an HTTP server.
func (e *Eweb) Start(address string) error {
	e.Server.Addr = address
	return e.StartServer(e.Server)
}

// StartTLS starts an HTTPS server.
func (e *Eweb) StartTLS(address string, certFile, keyFile string) (err error) {
	if certFile == "" || keyFile == "" {
		return errors.New("invalid tls configuration")
	}
	c := new(tls.Config)
	c.Certificates = make([]tls.Certificate, 1)
	c.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return
	}
	return e.startTLS(address, c)
}

// StartAutoTLS starts an HTTPS server using certificates automatically installed from https://letsencrypt.org.
func (e *Eweb) StartAutoTLS(address string) error {
	go http.ListenAndServe(":http", e.AutoTLSManager.HTTPHandler(nil))

	c := new(tls.Config)
	c.GetCertificate = e.AutoTLSManager.GetCertificate
	return e.startTLS(address, c)
}

func (e Eweb) StartTLSConfig(addr string, c *tls.Config) error {
	return e.startTLS(addr, c)
}

func (e *Eweb) startTLS(addr string, c *tls.Config) error {
	s := e.TLSServer
	s.Addr = addr
	s.TLSConfig = c

	if !e.DisableHTTP2 {
		s.TLSConfig.NextProtos = append(s.TLSConfig.NextProtos, "h2")
	}
	return e.StartServer(e.TLSServer)
}

// StartServer starts a custom http server.
func (e *Eweb) StartServer(s *http.Server) (err error) {
	s.ErrorLog = stdLog.New(os.Stderr, "", stdLog.LstdFlags)
	s.Handler = e

	kl, err := newListener(s.Addr)
	if err != nil {
		return err
	}

	// for http server
	if s.TLSConfig != nil {
		// https
		return s.Serve(tls.NewListener(kl, s.TLSConfig))
	}
	return s.Serve(kl)

}
