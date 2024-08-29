package sqlitez

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/bxcodec/go-clean-arch/domain"
	"github.com/bxcodec/go-clean-arch/internal/repository"
	"github.com/sirupsen/logrus"
)

type BMIrepository struct {
	DB *sql.DB
}

// NewMysqlAuthorRepository will create an implementation of author.Repository
func NewBMIRepository(db *sql.DB) *BMIrepository {
	return &BMIrepository{
		DB: db,
	}
}

func (m *BMIrepository) Fetch(ctx context.Context, cursor string, num int64) (res []domain.BMI, nextCursor string, err error) {
	query := `SELECT id,user_name, weight, height, bmi , created_at
  						FROM bmi_records WHERE created_at > ? ORDER BY created_at LIMIT ? `

	decodedCursor, err := repository.DecodeCursor(cursor)
	if err != nil && cursor != "" {
		return nil, "", domain.ErrBadParamInput
	}

	res, err = m.fetch(ctx, query, decodedCursor, num)
	if err != nil {
		return nil, "", err
	}

	if len(res) == int(num) {
		nextCursor = repository.EncodeCursor(res[len(res)-1].CreatedAt)
	}

	return
}

func (m *BMIrepository) Store(ctx context.Context, a *domain.BMI) (err error) {
	query := `INSERT INTO bmi_records ( user_name, weight, height, bmi, created_at)
	VALUES ( ?, ?, ?, ?, ?);`
	stmt, err := m.DB.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, a.UserName, a.Weight, a.Height, a.BMI, time.Now().UTC())
	if err != nil {
		return
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		return
	}
	a.ID = lastID
	return nil
}

func (m *BMIrepository) Delete(ctx context.Context, id int64) (err error) {
	query := "DELETE FROM bmi_records WHERE id = ?"

	stmt, err := m.DB.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return
	}

	rowsAfected, err := res.RowsAffected()
	if err != nil {
		return
	}

	if rowsAfected != 1 {
		err = fmt.Errorf("weird  Behavior. Total Affected: %d", rowsAfected)
		return
	}

	return
}

func (m *BMIrepository) GetByName(ctx context.Context, name string) (res []domain.BMI, err error) {
	query := `SELECT id,user_name, weight, height, bmi, created_at
  						FROM bmi_records WHERE user_name = ?`

	res, err = m.fetch(ctx, query, name)
	if err != nil {
		return nil, err
	}

	return
}

func (m *BMIrepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.BMI, err error) {
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			logrus.Error(errRow)
		}
	}()

	result = make([]domain.BMI, 0)
	for rows.Next() {
		t := domain.BMI{}
		err = rows.Scan(
			&t.ID,
			&t.UserName,
			&t.Weight,
			&t.Height,
			&t.BMI,
			&t.CreatedAt,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}
