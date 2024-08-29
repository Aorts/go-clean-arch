package domain

import (
	"math"
	"time"
)

type BMI struct {
	ID        int64     `json:"id"`
	UserName  string    `json:"userName" validate:"required"`
	Weight    float64   `json:"weight" validate:"required"`
	Height    float64   `json:"height" validate:"required"`
	BMI       float64   `json:"bmi"`
	CreatedAt time.Time `json:"created_at"`
}

func (b *BMI) CalculateBMI() {
	if b.Height > 0 && b.Weight > 0 {
		b.BMI = (b.Weight / math.Pow(b.Height, 2))
	}
}
