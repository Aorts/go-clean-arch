package rest_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/bxcodec/go-clean-arch/domain"
	"github.com/bxcodec/go-clean-arch/internal/rest"
	"github.com/bxcodec/go-clean-arch/internal/rest/mocks"
	faker "github.com/go-faker/faker/v4"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestFetchBMI(t *testing.T) {
	var mockBMI domain.BMI
	err := faker.FakeData(&mockBMI)
	assert.NoError(t, err)
	mockUCase := new(mocks.BMIService)
	mockListBMI := make([]domain.BMI, 0)
	mockListBMI = append(mockListBMI, mockBMI)
	num := 1
	cursor := "2"
	mockUCase.On("Fetch", mock.Anything, cursor, int64(num)).Return(mockListBMI, "10", nil)

	e := echo.New()
	req, err := http.NewRequestWithContext(context.TODO(),
		echo.GET, "/api/v1/bmis?num=1&cursor="+cursor, strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	handler := rest.BMIHandler{
		Service: mockUCase,
	}
	err = handler.FetchBMI(c)
	require.NoError(t, err)

	responseCursor := rec.Header().Get("X-Cursor")
	assert.Equal(t, "10", responseCursor)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestGetByName(t *testing.T) {
	var mockBMI []domain.BMI
	err := faker.FakeData(&mockBMI)
	assert.NoError(t, err)

	mockUCase := new(mocks.BMIService)

	username := mockBMI[0].UserName

	mockUCase.On("GetByName", context.TODO(), username).Return(mockBMI, nil)

	e := echo.New()
	req, err := http.NewRequestWithContext(context.TODO(), echo.GET, "/api/v1/bmi?userName="+username, strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/bmi")
	c.SetParamNames("userName")
	c.SetParamValues(username)
	handler := rest.BMIHandler{
		Service: mockUCase,
	}
	err = handler.GetByName(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestStoreBMI(t *testing.T) {
	mockBMI := domain.BMI{
		UserName: "Title",
		Height:   1.75,
		Weight:   55,
	}

	tempmockBMI := mockBMI
	tempmockBMI.ID = 0
	mockUCase := new(mocks.BMIService)

	j, err := json.Marshal(tempmockBMI)
	assert.NoError(t, err)

	mockUCase.On("Store", mock.Anything, mock.AnythingOfType("*domain.BMI")).Return(nil)

	e := echo.New()
	req, err := http.NewRequestWithContext(context.TODO(), echo.POST, "/api/v1/bmi", strings.NewReader(string(j)))
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/bmi")

	handler := rest.BMIHandler{
		Service: mockUCase,
	}
	err = handler.Store(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusCreated, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestDeleteBMI(t *testing.T) {
	var mockBMI domain.BMI
	err := faker.FakeData(&mockBMI)
	assert.NoError(t, err)

	mockUCase := new(mocks.BMIService)

	num := int(mockBMI.ID)

	mockUCase.On("Delete", mock.Anything, int64(num)).Return(nil)

	e := echo.New()
	req, err := http.NewRequestWithContext(context.TODO(), echo.DELETE, "/api/v1/bmi/"+strconv.Itoa(num), strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/bmi/")
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(num))
	handler := rest.BMIHandler{
		Service: mockUCase,
	}
	err = handler.Delete(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusNoContent, rec.Code)
	mockUCase.AssertExpectations(t)
}
