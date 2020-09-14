package handlers

import (
	"consequences-api/api/models"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// RunConsequences lists alerts for a single instrument
func RunConsequences() echo.HandlerFunc {
	return func(c echo.Context) error {
		var i models.ConsequencesInputCollection
		fmt.Println(c)
		if err := c.Bind(&i); err != nil {
			return c.String(http.StatusBadRequest, "Invalid Input")
		}
		s, err := models.RunConsequences(i)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, s)
	}
}
