package model

import "gorm.io/gorm"

type ClientPayment struct {
	ClientID        uint
	SeminarID       uint
	PaymentStatusID uint
	PaymentStatus   PaymentStatus
	Type            PaymentType
}

type PaymentType struct {
	gorm.Model
	Code string
	Name string
}

type PaymentStatus struct {
	gorm.Model
	Code string
	Name string
}
