package routes

import (
	"github.com/gwaylib/eweb"
	echo "github.com/labstack/echo/v4"
)

// register routes
func init() {
	e := eweb.Default()
	e.GET("/login", LoginView)
}

// Login Handler
func LoginView(c echo.Context) error {
	return c.HTML(200, "test/tpl/login")
}
