package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Email    string
	Password string
}

func CheckDatabase() {
	databaseCheck := false
	for !databaseCheck {
		productionDatabase, err := gorm.Open(mysql.Open(databaseConnection), &gorm.Config{})
		productionDB, _ := productionDatabase.DB()
		if err != nil {
			fmt.Println("Problem opening database, looks like it does not exist")
			database, err := gorm.Open(mysql.Open(connection), &gorm.Config{})
			sqlDB, _ := database.DB()
			if err != nil {
				fmt.Println("Problem opening main mysql database: " + err.Error())
				continue
			}
			fmt.Println("Creating jake database")
			database.Exec("CREATE DATABASE medium;")
			sqlDB.Close()
		}
		fmt.Println("Medium database already exists")
		if !productionDatabase.Migrator().HasTable(&User{}) {
			fmt.Println("Creating table User")
			err := productionDatabase.Migrator().CreateTable(&User{})
			if err != nil {
				fmt.Println("Cannot create table: " + err.Error())
				return
			}
			fmt.Println("Creating user admin")

			password, err := hashPasswordFromString([]byte("54321"))
			if err != nil {
				fmt.Println("Cannot hash password: " + err.Error())
				return
			}
			user := User{
				Email:    "admin@admin.com",
				Password: password,
			}
			productionDatabase.Create(&user)
		} else {
			fmt.Println("Updating table User")
			err := productionDatabase.Migrator().AutoMigrate(&User{})
			if err != nil {
				fmt.Println("Cannot update table: " + err.Error())
				return
			}
		}
		productionDB.Close()
		databaseCheck = true
		time.Sleep(1 * time.Second)
	}
	fmt.Println("Checking database done")
}

func hashPasswordFromString(pwd []byte) (string, error) {
	fmt.Println("Hashing password")
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		fmt.Println("Cannot hash password: " + err.Error())
		return "", err
	}
	fmt.Println("Password hashed")
	return string(hash), nil
}
