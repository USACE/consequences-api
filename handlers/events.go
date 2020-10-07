package handlers

import (
	"net/http"
	"time"

	"github.com/USACE/consequences-api/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

func ListEvents(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ee, err := models.ListEvents(db)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, ee)
	}
}

func CreateEvent(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var e models.Event
		if err := c.Bind(&e); err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		actor, ok := c.Get("actor").(int)
		if !ok {
			return c.String(http.StatusBadRequest, "Something happened with type assertion")
		}
		e.Creator = actor
		e.CreateDate = time.Now()

		eCreated, err := models.CreateEvent(db, &e)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusCreated, eCreated)
	}
}

func DeleteEvent(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		eventID, err := uuid.Parse(c.Param("event_id"))
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		err = models.DeleteEvent(db, &eventID)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		// Empty response
		return c.JSON(http.StatusOK, make(map[string]interface{}))
	}
}
