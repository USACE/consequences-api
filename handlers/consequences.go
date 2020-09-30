package handlers

import (
	"net/http"

	"github.com/USACE/consequences-api/models"

	"github.com/labstack/echo/v4"
)

// RunConsequences lists alerts for a single instrument
func RunConsequencesByBoundingBox() echo.HandlerFunc {
	return func(c echo.Context) error {
		var i models.ConsequencesBoundingBox
		if err := c.Bind(&i); err != nil {
			return c.String(http.StatusBadRequest, "Invalid Input")
		}
		s, err := models.RunConsequencesByBoundingBox(i)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, s)
	}
}

// RunConsequences lists alerts for a single instrument
func RunConsequencesByFips() echo.HandlerFunc {
	return func(c echo.Context) error {
		fips := c.Param("fips_code")
		if fips == "" {
			return c.String(http.StatusBadRequest, "Invalid Input")
		}
		s, err := models.RunConsequencesByFips(fips)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.String(http.StatusOK, s)
	}
}
