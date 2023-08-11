package utils

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

// Response in JSON format
func JsonResponse(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, data)
}

// Response in XML format
func XmlResponse(c echo.Context, data interface{}) error {
	return c.XML(http.StatusOK, data)
}
