package controllers

import (
	"context"
	"log"
	"time"
	"fmt"

	"auth-regapp/database"
	"auth-regapp/helpers"
	"auth-regapp/model"
	
	"github.com/gedrimas/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt" 
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var companyCollection *mongo.Collection = database.OpenCollection(database.Client, "company")

var validate = validator.New()

func EncodePassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func CheckPassword(enteredPassword string, fetchedPassword string) bool {

	err := bcrypt.CompareHashAndPassword([]byte(fetchedPassword), []byte(enteredPassword))
	check := true

	if err != nil {
		fmt.Println("err", err)
		check = false
	}
	return check
}

type SignupRequest struct {
	Username  string `json:"username" validate:"required,min=4,max=100"`
	Password  string `json:"password" validate:"required,min=8"`
	Email     string `json:"email" validate:"email,required"`
	User_type string `json:"user_type" validate:"required,eq=ADMIN|eq=USER"`
	Company   string `json:"company" validate:"required"`
	Contacts   string `json:"contacts"`
}



func Signup() gin.HandlerFunc {

	return func(c *gin.Context) {

		sendor := response.NewSendor(c)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var req SignupRequest
		if err := c.BindJSON(&req); err != nil {
			badRequest := response.GetResponse("badRequest")
			badRequest.SetData(err.Error())
			sendor.Send(badRequest)
			return
		}

		regexEmail := bson.M{"$regex": primitive.Regex{Pattern: req.Email, Options: "i"}}
		emailCount, emailErr := userCollection.CountDocuments(ctx, bson.M{"email": regexEmail})
		regexUsername := bson.M{"$regex": primitive.Regex{Pattern: req.Username, Options: "i"}}
		usernameCount, usernameErr := userCollection.CountDocuments(ctx, bson.M{"username": regexUsername})
		
		if emailErr != nil {
			emailGetErr := response.GetResponse("emailGetError")
			emailGetErr.SetData(emailErr)
			sendor.Send(emailGetErr)
			log.Panic(emailErr)
		}

		if emailCount > 0 {		
			emailExistsErr := response.GetResponse("emailExistsError")
			emailExistsErr.SetData(nil)
			sendor.Send(emailExistsErr)
			return
		}

		if usernameErr != nil {
			userGetErr := response.GetResponse("userGetError")
			userGetErr.SetData(usernameErr)
			sendor.Send(userGetErr)
			log.Panic(usernameErr)
		}

		if usernameCount > 0 {
			userExistsErr := response.GetResponse("userExistsError")
			userExistsErr.SetData(nil)
			sendor.Send(userExistsErr)
			return
		}

		var user model.User
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()

		var companyId interface{}
		if req.User_type == "ADMIN" {
			regexCompany := bson.M{"$regex": primitive.Regex{Pattern: req.Company, Options: "i"}}
			companyCount, companyErr := companyCollection.CountDocuments(ctx, bson.M{"company": regexCompany})

			if companyErr != nil {
				companyGetErr := response.GetResponse("companyGetError")
				companyGetErr.SetData(companyErr)
				sendor.Send(companyGetErr)
				log.Panic(companyErr)
			}

			if companyCount > 0 {
				companyExistsErr := response.GetResponse("companyExistsError")
				companyExistsErr.SetData(nil)
				sendor.Send(companyExistsErr)
				return
			}

			var company model.Company
			company.ID = primitive.NewObjectID()
			company.Conpmay_id = company.ID.Hex()
			company.Company = req.Company
			company.Admin_id = user.User_id

			result, err := companyCollection.InsertOne(ctx, company)

			if err != nil {
				companyCreationErr := response.GetResponse("companyCreationError")
				companyCreationErr.SetData(companyErr)
				sendor.Send(companyCreationErr)
				log.Panic(err)
			}
			companyId = result.InsertedID
		}

		if id, ok := companyId.(primitive.ObjectID); ok {
			user.Company_id = id.Hex()
		} else {
			user.Company_id = ""
		}

		password := EncodePassword(req.Password)
		user.Password = password
		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Username = req.Username
		user.Password = req.Password
		user.Email = req.Email
		user.User_type = req.User_type

		token, refreshToken, _ := helpers.GenerateAllTokens(user.Username, user.User_type, user.User_id, user.Company_id)
		user.Token = token
		user.Refresh_token = refreshToken

		newUser := model.User{
			ID:            user.ID,
			User_id:       user.User_id,
			Username:      req.Username,
			Contacts:	   req.Contacts,	
			User_type:     req.User_type,
			Email:         req.Email,
			Company_id:    user.Company_id,
			Password:      password,
			Token:         user.Token,
			Refresh_token: user.Refresh_token,
			Created_at:    user.Created_at,
			Updated_at:    user.Updated_at,
		}

		result, err := userCollection.InsertOne(ctx, newUser)

		if err != nil {
			userCreationErr := response.GetResponse("userCreationError")
			userCreationErr.SetData(err)
			sendor.Send(userCreationErr)
			log.Panic(err)
		}

		var newUserId string
		if id, ok := result.InsertedID.(primitive.ObjectID); ok {
			newUserId = id.Hex()
		}

		var createdUser model.User
		err = userCollection.FindOne(ctx, bson.D{{"user_id", newUserId }}).Decode(&createdUser)
		if err != nil {
			userGetErr := response.GetResponse("userGetError")
			userGetErr.SetData(err)
			sendor.Send(userGetErr)
			log.Panic(err)
		}

		
		data := map[string]string{"token": createdUser.Token, "user_type": createdUser.User_type}

		signupSuccess := response.GetResponse("signuppSuccess")
		signupSuccess.SetData(data)
		sendor.Send(signupSuccess)
		return

	}
}

type LoginRequest struct {
	Password  string `json:"password" validate:"required,min=8"`
	Email     string `json:"email" validate:"email,required"`
	User_type string `json:"user_type" validate:"required,eq=ADMIN|eq=USER"`
}


func Login() gin.HandlerFunc {

	return func(c *gin.Context) {

		sendor := response.NewSendor(c)

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		
		var req LoginRequest 
		var fetchedUser model.User

		if err := c.BindJSON(&req); err != nil {
			badRequest := response.GetResponse("badRequest")
			badRequest.SetData(err.Error())
			sendor.Send(badRequest)
			return
		}

		emailErr := userCollection.FindOne(ctx, bson.M{"email": req.Email}).Decode(&fetchedUser)
		if emailErr != nil {
			loginError := response.GetResponse("loginError")
			loginError.SetData(nil)
			sendor.Send(loginError)
			return
		}

		
		if passErr := CheckPassword(req.Password, fetchedUser.Password); passErr != true {
			loginError := response.GetResponse("loginError")
			loginError.SetData(nil)
			sendor.Send(loginError)
			return
		}

		token, refreshToken, _ := helpers.GenerateAllTokens(fetchedUser.Username, fetchedUser.User_type, fetchedUser.User_id, fetchedUser.Company_id)
		UpdateTokens(token, refreshToken, fetchedUser.User_id)

		userErr := userCollection.FindOne(ctx, bson.M{"user_id": fetchedUser.User_id}).Decode(&fetchedUser)
		if userErr != nil {
			userGetErr := response.GetResponse("userGetError")
			userGetErr.SetData(userErr)
			sendor.Send(userGetErr)
			log.Panic(userErr)
		}

		data := map[string]string{"token": fetchedUser.Token, "user_type": fetchedUser.User_type}
		loginSuccess := response.GetResponse("loginSuccess")
		loginSuccess.SetData(data)
		sendor.Send(loginSuccess)
		return

	}
}

func UpdateTokens(token string, refreshToken string, userId string) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var updateToken primitive.D
	updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	updateToken = append(
		updateToken,
		bson.E{Key: "token", Value: token},
		bson.E{Key: "refresh_token", Value: refreshToken},
		bson.E{Key: "updated_at", Value: updated_at},
	)

	upsert := true
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	//upset := true
	filter := bson.M{"user_id": userId}
	_, err := userCollection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: updateToken}}, &opt)

	if err != nil {
		log.Panic(err)
		return
	}
	return

}
