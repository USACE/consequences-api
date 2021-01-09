package handlers

import (
	"net/http"
	"time"

	"github.com/USACE/consequences-api/models"
	"github.com/USACE/go-simple-asyncer/asyncer"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/labstack/echo/v4"
)

// GetCompute gets a single compute
func GetCompute(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		computeID, err := uuid.Parse(c.Param("compute_id"))
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		compute, err := models.GetCompute(db, &computeID)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, &compute)
	}
}

// RunConsequencesByFips runs a FIPS Compute
func RunConsequencesByFips(db *sqlx.DB, ae asyncer.Asyncer) echo.HandlerFunc {
	return func(c echo.Context) error {
		fips := c.Param("fips_code")
		if fips == "" {
			return c.String(http.StatusBadRequest, "Invalid Input")
		}
		eventID, err := uuid.Parse(c.Param("event_id"))
		if err != nil {
			return c.String(http.StatusBadRequest, "Invalid Event ID")
		}
		actor, ok := c.Get("actor").(int)
		if !ok {
			return c.String(http.StatusBadRequest, "Something happened with type assertion")
		}
		// Build Valid Compute
		compute := models.Compute{
			EventID: eventID,
			FIPS:    fips,
		}
		compute.Creator = actor
		compute.CreateDate = time.Now()

		computeNew, err := models.RunCompute(db, ae, &compute)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, &computeNew)
	}
}
