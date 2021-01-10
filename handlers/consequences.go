package handlers

import (
	"net/http"

	"github.com/USACE/consequences-api/models"
	"github.com/USACE/go-consequences/compute"
	"github.com/jmoiron/sqlx"

	"github.com/labstack/echo/v4"
)

// RunConsequences lists alerts for a single instrument
func RunConsequencesByBoundingBox() echo.HandlerFunc {
	return func(c echo.Context) error {
		var i compute.Bbox
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
func RunAgConsequencesByXY() echo.HandlerFunc {
	return func(c echo.Context) error {
		year := c.Param("year")
		x := c.Param("x")
		y := c.Param("y")
		at := c.Param("arrivaltime")
		duration := c.Param("duration")
		if year == "" {
			return c.String(http.StatusBadRequest, "Please Specify a Year")
		}
		if x == "" {
			return c.String(http.StatusBadRequest, "Please Specify an X coordinate")
		}
		if y == "" {
			return c.String(http.StatusBadRequest, "Please Specify a Y coordinate")
		}
		if at == "" {
			return c.String(http.StatusBadRequest, "Please Specify an Arrival Time")
		}
		if duration == "" {
			return c.String(http.StatusBadRequest, "Please Specify a Duration")
		}
		var i = compute.FipsCode{FIPS: fips}
		s, err := models.RunAgConsequencesByXY(i, 5.67)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, s)
	}
}

func RunConsequencesByFips(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		fips := c.Param("fips_code")
		if fips == "" {
			return c.String(http.StatusBadRequest, "Invalid Input")
		}
		var i = compute.FipsCode{FIPS: fips}
		s, err := models.RunConsequencesByFips(i, 5.67)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, s)
	}
}
