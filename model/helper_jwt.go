package model

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
)

var AppSecret = "J%df4e8hcjvbkjclkjkklfgki843895iojfdnvufh98"
var AppIss = "sshfortress"
var ExpireTime = time.Hour * 24 * 30

func jwtGenerateToken(m *User) (*jwtObj, error) {
	m.Password = ""
	m.InputSshPassword = ""
	m.InputPassword = ""
	expireTime := time.Now().Add(ExpireTime)
	stdClaims := jwt.StandardClaims{
		ExpiresAt: expireTime.Unix(),
		IssuedAt:  time.Now().Unix(),
		Id:        fmt.Sprintf("%d", m.Id),
		Issuer:    AppIss,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, stdClaims)
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(AppSecret))
	if err != nil {
		logrus.WithError(err).Error("config is wrong, can not generate jwt")
	}
	data := &jwtObj{User: *m, Token: tokenString, Expire: expireTime, ExpireTs: expireTime.Unix()}
	return data, err
}

type jwtObj struct {
	User
	Token    string    `json:"token"`
	Expire   time.Time `json:"expire"`
	ExpireTs int64     `json:"expire_ts"`
}

//JwtParseUser
func JwtParseUser(tokenString string) (uint, error) {
	if tokenString == "" {
		return 0, errors.New("no token is found in Authorization Bearer")
	}
	claims := jwt.StandardClaims{}
	_, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(AppSecret), nil
	})
	if err != nil {
		return 0, err
	}
	if claims.VerifyExpiresAt(time.Now().Unix(), true) == false {
		return 0, errors.New("token is expired")
	}
	if !claims.VerifyIssuer(AppIss, true) {
		return 0, errors.New("token's issuer is wrong")
	}
	uid, err := strconv.ParseUint(claims.Id, 10, 64)
	return uint(uid), err
}
