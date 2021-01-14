package public

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
)

func JWTEncode(claims jwt.StandardClaims) (string, error) {
	mySignKey := []byte(JWTSignKey)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(mySignKey)
}

func JWTDeCode(tokenString string) (*jwt.StandardClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWTSignKey), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*jwt.StandardClaims); ok {
		//if claims.ExpiresAt < time.Now().Unix() {
		//	return nil, errors.New("request expired")
		//}
		return claims, nil
	} else {
		return nil, errors.New("token is not jwt.StanderClaims")
	}
}
