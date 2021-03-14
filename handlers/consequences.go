package handlers

import (
	"net/http"

	"github.com/USACE/consequences-api/models"
	"github.com/jmoiron/sqlx"

	"github.com/labstack/echo/v4"
)

// ByStructureFromFile handles the arguments from file
func ByStructureFromFile() echo.HandlerFunc {
	return func(c echo.Context) error {
		var i models.Compute
		if err := c.Bind(&i); err != nil {
			return c.String(http.StatusBadRequest, "Invalid Input")
		}
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		c.Response().WriteHeader(http.StatusOk)
		//return json.NewEncoder(c.Response())
		s, err := models.ComputeByStructureFromFile(i,c.Response())
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, s)
	}
}


