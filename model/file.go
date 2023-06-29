package model

import "gorm.io/gorm"

type File struct {
	gorm.Model
	Name    string
	Content []byte
}
