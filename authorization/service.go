package authorization

import (
	"kedaiprogrammer/helpers"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

type Services interface {
	GenerateJWT(uuid string, username string) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}

type jwtServices struct {
}

var jwtKey = []byte(viper.GetString("SECRET.SECRET_KEY_JWT"))

type JWTClaim struct {
	Uuid     string `json:"uuid"`
	Username string `json:"username"`
	jwt.StandardClaims
}

func NewServices() *jwtServices {
	return &jwtServices{}
}

func (s *jwtServices) GenerateJWT(uuid string, username string) (tokenString string, err error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	uuidHashed, _ := helpers.EncryptUUID(uuid)

	claims := &JWTClaim{
		Username: username,
		Uuid:     uuidHashed,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(jwtKey)
	return
}

func (s *jwtServices) ValidateToken(signedToken string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(signedToken, &JWTClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	if err != nil {
		return nil, err
	}
	return token, err
}
