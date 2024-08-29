package sqlitez_test

import (
	"context"
	"testing"
	"time"

	"github.com/bxcodec/go-clean-arch/domain"
	"github.com/bxcodec/go-clean-arch/internal/repository"
	"github.com/bxcodec/go-clean-arch/internal/repository/sqlitez"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestFetchBMI(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mockBMI := []domain.BMI{
		{
			ID: 1, UserName: "title 1", Weight: 55,
			Height: 1.75, BMI: 17.959183, CreatedAt: time.Now().Add(-time.Hour),
		},
		{
			ID: 2, UserName: "title 2", Weight: 70,
			Height: 1.75, BMI: 22.857142857142858, CreatedAt: time.Now(),
		},
		{
			ID: 3, UserName: "title 3", Weight: 70,
			Height: 1.75, BMI: 22.857142857142858, CreatedAt: time.Now(),
		},
	}

	rows := sqlmock.NewRows([]string{"id", "user_name", "weight", "height", "bmi", "created_at"}).
		AddRow(mockBMI[0].ID, mockBMI[0].UserName, mockBMI[0].Weight,
			mockBMI[0].Height, mockBMI[0].BMI, mockBMI[0].CreatedAt).
		AddRow(mockBMI[1].ID, mockBMI[1].UserName, mockBMI[1].Weight,
			mockBMI[1].Height, mockBMI[1].BMI, mockBMI[1].CreatedAt)

	query := `SELECT id,user_name, weight, height, bmi , created_at FROM bmi_records`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := sqlitez.NewBMIRepository(db)
	cursor := repository.EncodeCursor(mockBMI[0].CreatedAt)
	num := int64(2)
	list, nextCursor, err := a.Fetch(context.TODO(), cursor, num)
	assert.NotEmpty(t, nextCursor)
	assert.NoError(t, err)
	assert.Len(t, list, 2)
}

func TestGetBMIByName(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows := sqlmock.NewRows([]string{"id", "user_name", "weight", "height", "bmi", "created_at"}).
		AddRow(1, "title 1", 55, 1.75, 22, time.Now())

	query := `SELECT id,user_name, weight, height, bmi, created_at FROM bmi_records WHERE user_name = ?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := sqlitez.NewBMIRepository(db)

	anArticle, err := a.GetByName(context.TODO(), "title 1")
	assert.NoError(t, err)
	assert.NotNil(t, anArticle)
}

func TestStoreArticle(t *testing.T) {
	now := time.Now()
	ar := &domain.BMI{
		UserName:  "Judul",
		Height:    1.75,
		CreatedAt: now,
		Weight:    60,
	}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	query := `INSERT INTO bmi_records ( user_name, weight, height, bmi, created_at)
	VALUES ( ?, ?, ?, ?, ?);`
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(ar.UserName, ar.Weight, ar.Height, ar.BMI, ar.CreatedAt).WillReturnResult(sqlmock.NewResult(12, 1))

	a := sqlitez.NewBMIRepository(db)

	err = a.Store(context.TODO(), ar)
	assert.NoError(t, err)
	assert.Equal(t, int64(12), ar.ID)
}

func TestDeleteArticle(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	query := "DELETE FROM bmi_records WHERE id = \\?"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(12).WillReturnResult(sqlmock.NewResult(12, 1))

	a := sqlitez.NewBMIRepository(db)

	num := int64(12)
	err = a.Delete(context.TODO(), num)
	assert.NoError(t, err)
}
