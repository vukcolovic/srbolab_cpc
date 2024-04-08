package model

import "gorm.io/gorm"

type Partner struct {
	gorm.Model
	Name string `json:"name"`
}
