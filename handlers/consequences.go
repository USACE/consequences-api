package handlers

import (
	"net/http"

	"github.com/USACE/consequences-api/models"

	"github.com/labstack/echo/v4"
)

// RunConsequences lists alerts for a single instrument
func RunConsequences() echo.HandlerFunc {
	return func(c echo.Context) error {
		var i models.ConsequencesInputCollection
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
func GetStructureLocations() echo.HandlerFunc {
	return func(c echo.Context) error {
		var i models.ConsequencesBoundingBox
		if err := c.Bind(&i); err != nil {
			return c.String(http.StatusBadRequest, "Invalid Input")
		}
		s, err := models.GetInventory(i)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, s)
	}
}
