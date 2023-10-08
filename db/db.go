package db

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"srbolab_cpc/config"
	"srbolab_cpc/model"
	"strconv"
)

var Client *gorm.DB

func Connect(conf config.Config) error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", conf.DbHost, conf.DbUser, conf.DbPassword, conf.DbName, strconv.Itoa(conf.DbPort))
	var err error
	Client, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	err = Client.AutoMigrate(
		&model.Location{},
		&model.File{},
		&model.User{},
		&model.Client{},
		&model.Seminar{},
		&model.SeminarDay{},
		&model.SeminarClass{},
		&model.ClassRoom{},
		&model.ClientSeminar{},
		&model.ClientPresence{},
		&model.Question{},
		&model.Answer{},
		&model.Test{},
		&model.ClientTest{},
		&model.SeminarClassName{},
	)
	if err != nil {
		return err
	}

	return nil
}
