package model

import (
	"gorm.io/gorm"
	"strconv"
)

type Person struct {
	FirstName   string `json:"first_name"`
	MiddleName  string `json:"middle_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}

func (p Person) FullName() string {
	return p.FirstName + " " + p.LastName
}

func (p Person) FullNameWithMiddleName() string {
	return p.FirstName + " (" + p.MiddleName + ") " + p.LastName
}

func (c Client) GetBirthDate() string {
	if len(*c.JMBG) < 13 {
		return "Nevalidan jmbg"
	}
	str := *c.JMBG
	year, _ := strconv.Atoi(str[4:7])
	yearPrefix := 1
	if year < 100 {
		yearPrefix++
	}
	return str[0:2] + "." + str[2:4] + "." + strconv.Itoa(yearPrefix) + str[4:7] + "."
}

type User struct {
	gorm.Model
	Person   Person `json:"person" gorm:"embedded"`
	Password string `json:"password"`
	Roles    []Role `json:"roles" gorm:"many2many:user_role;"`
}

func (u *User) AfterFind(tx *gorm.DB) (err error) {
	u.Password = ""
	return nil
}

type Role struct {
	gorm.Model
	Code string `json:"code"`
	Name string `json:"name"`
}
