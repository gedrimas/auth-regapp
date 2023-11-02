package helpers

import (
	"log"
	jwt "github.com/dgrijalva/jwt-go"
	"time"
	"fmt"
)

type JwtDetails struct {
	Email string
	Username string
	Uid string
	User_type string
	User_id string
	Company_id string
	jwt.StandardClaims
}

var SECRET_JWT_KEY = EnvFileVal("SECRET_KEY")

func GenerateAllTokens(
	userName string,
	userType string,
	userId string,
	companyId string,
) (
	token string,
	refreshToken string,
	err error,
) {

	fmt.Println("userType", userType)

	claims := &JwtDetails{
		Username: userName,
		User_type: userType,
		User_id: userId,
		Company_id: companyId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(50)).Unix(),
		},
	}

	refreshClaims := &JwtDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(100)).Unix(),
		},
	}

	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_JWT_KEY))
	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_JWT_KEY))

	if err != nil {
		log.Panic(err)
		return
	}
	return token, refreshToken, err
}




