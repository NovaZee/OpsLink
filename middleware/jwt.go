package middleware

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
)

var (
	TokenExpired   = errors.New("token expired")
	TokenNotValid  = errors.New("token not valid")
	TokenMalformed = errors.New("token is malformed")
	TokenInvalid   = errors.New("invalid")
)

type UserInfo struct {
	//ID int `json:"id"` // userId
	//Email string `json:"email"` // user email
	Name     string `json:"name"` // username
	PassWord string `json:"pwd"`
	//Role int    `json:"role"`
}

const SignKey = "OpsLink"

var JWTUtil = NewJWT(SignKey)

type CustomClaims struct {
	UserInfo
	jwt.StandardClaims
}

type JWT struct {
	SignKey []byte
}

func NewJWT(key string) *JWT {
	return &JWT{[]byte(key)}
}

func (j *JWT) GetSignKey() string {
	return string(j.SignKey)
}

func (j *JWT) SetSignKey(key string) {
	j.SignKey = []byte(key)
}

func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	return token.SignedString(j.SignKey)
}

func (j *JWT) ParseToken(t string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(t, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SignKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, TokenExpired
			}
		}
		return nil, err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, TokenInvalid
}
