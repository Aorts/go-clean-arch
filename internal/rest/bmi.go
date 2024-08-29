package rest

import (
	"context"
	"net/http"
	"strconv"

	"github.com/bxcodec/go-clean-arch/domain"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

const defaultNum = 10

type BMIHandler struct {
	Service BMIService
}

type BMIService interface {
	Fetch(ctx context.Context, cursor string, num int64) ([]domain.BMI, string, error)
	Store(context.Context, *domain.BMI) error
	Delete(ctx context.Context, id int64) error
	GetByName(ctx context.Context, name string) ([]domain.BMI, error)
}

func NewBMIHandler(e *echo.Echo, svc BMIService) {
	handler := &BMIHandler{
		Service: svc,
	}
	e.GET("/api/v1/bmi", handler.GetByName)
	e.POST("/api/v1/bmi", handler.Store)
	e.GET("/api/v1/bmis", handler.FetchBMI)
	e.DELETE("/api/v1/bmi/:id", handler.Delete)
}

func (b *BMIHandler) FetchBMI(c echo.Context) error {

	numS := c.QueryParam("num")
	num, err := strconv.Atoi(numS)
	if err != nil || num == 0 {
		num = defaultNum
	}

	cursor := c.QueryParam("cursor")

	listAr, nextCursor, err := b.Service.Fetch(c.Request().Context(), cursor, int64(num))
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	c.Response().Header().Set(`X-Cursor`, nextCursor)
	return c.JSON(http.StatusOK, listAr)
}

func (b *BMIHandler) GetByName(c echo.Context) error {
	name := c.QueryParam("userName")
	res, err := b.Service.GetByName(c.Request().Context(), name)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}

func (b *BMIHandler) Store(c echo.Context) error {
	var bmi domain.BMI
	err := c.Bind(&bmi)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	bmi.CalculateBMI()

	err = b.Service.Store(c.Request().Context(), &bmi)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, bmi)
}

func (b *BMIHandler) Delete(c echo.Context) error {
	idP, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, domain.ErrNotFound.Error())
	}

	id := int64(idP)
	err = b.Service.Delete(c.Request().Context(), id)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	logrus.Error(err)
	switch err {
	case domain.ErrInternalServerError:
		return http.StatusInternalServerError
	case domain.ErrNotFound:
		return http.StatusNotFound
	case domain.ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

type ResponseError struct {
	Message string `json:"message"`
}
