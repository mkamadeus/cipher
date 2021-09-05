package hill

import "github.com/labstack/echo/v4"

func Register(e *echo.Echo) {
	group := e.Group("/hill")
	group.POST("/encrypt", Encrypt)
	group.POST("/decrypt", Decrypt)
}
