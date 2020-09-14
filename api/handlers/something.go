package handlers

import (
	"consequences-api/api/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

// GetSomething lists alerts for a single instrument
func GetSomething() echo.HandlerFunc {
	return func(c echo.Context) error {
		s, err := models.GetSomething()
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, s)
	}
}
