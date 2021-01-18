package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"net/http"
)

type VerifyUserInput struct {
	UserEmail          string
	UserPassword       string
	UserSessionStorage string
}

type VerifyUserOutput struct {
	Result       string
	Content      string
	SessionLogin string
}

func checkLogin(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	var data VerifyUserInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		fmt.Println("Error parsing data: " + err.Error())
		var responseData VerifyUserOutput
		responseData.Result = "nok"
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		return
	}
	fmt.Println("Verifying user started for " + data.UserEmail + ":" + data.UserPassword + ":" + data.UserSessionStorage)
	if len(data.UserSessionStorage) != 0 {
		fmt.Println("Session storage not empty")
		db, err := gorm.Open(mysql.Open(databaseConnection), &gorm.Config{})
		sqlDB, _ := db.DB()
		defer sqlDB.Close()
		if err != nil {
			fmt.Println("Cannot connect to database")
			var responseData VerifyUserOutput
			responseData.Result = "nok"
			writer.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(writer).Encode(responseData)
			return
		}
		var user User
		db.Where("password = ?", data.UserSessionStorage).Find(&user)
		if user.Password == data.UserSessionStorage {
			fmt.Println("User matches")
			file, err := ioutil.ReadFile("html/content.html")
			if err != nil {
				fmt.Println("Error reading file: " + err.Error())
				var responseData VerifyUserOutput
				responseData.Result = "nok"
				writer.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(writer).Encode(responseData)
				return
			}
			var responseData VerifyUserOutput
			responseData.Result = "ok"
			responseData.Content = string(file)
			responseData.SessionLogin = user.Password
			writer.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(writer).Encode(responseData)
			return
		} else {
			fmt.Println("User does not match")
			var responseData VerifyUserOutput
			responseData.Result = "nok"
			writer.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(writer).Encode(responseData)
			return
		}
	} else {
		fmt.Println("Session storage empty")
		if len(data.UserEmail) == 0 || len(data.UserPassword) == 0 {
			fmt.Println("Password/email empty")
			var responseData VerifyUserOutput
			responseData.Result = "nok"
			writer.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(writer).Encode(responseData)
			return
		} else {
			fmt.Println("Password/email not empty")
			db, err := gorm.Open(mysql.Open(databaseConnection), &gorm.Config{})
			sqlDB, _ := db.DB()
			defer sqlDB.Close()
			if err != nil {
				fmt.Println("Cannot connect to database")
				var responseData VerifyUserOutput
				responseData.Result = "nok"
				writer.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(writer).Encode(responseData)
				return
			}
			var user User
			db.Where("email = ?", data.UserEmail).Find(&user)
			userMatchesPassword := comparePasswords(user.Password, []byte(data.UserPassword))
			if userMatchesPassword {
				fmt.Println("User matches")
				file, err := ioutil.ReadFile("html/content.html")
				if err != nil {
					fmt.Println("Error reading file: " + err.Error())
					var responseData VerifyUserOutput
					responseData.Result = "nok"
					writer.Header().Set("Content-Type", "application/json")
					_ = json.NewEncoder(writer).Encode(responseData)
					return
				}
				var responseData VerifyUserOutput
				responseData.Result = "ok"
				responseData.Content = string(file)
				responseData.SessionLogin = user.Password
				writer.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(writer).Encode(responseData)
				return
			} else {
				fmt.Println("User does not match")
				var responseData VerifyUserOutput
				responseData.Result = "nok"
				writer.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(writer).Encode(responseData)
				return
			}
		}
	}
}

func comparePasswords(hashedPwd string, plainPwd []byte) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
