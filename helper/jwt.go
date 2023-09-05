package helper

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
)

func GenJWT(sub string, exp time.Duration) (string, error) {

	hmac_secret := []byte(viper.Get("JWT_SECRET_KEY").(string))
	// NOTE that JWT must be generated on backend side of your application!
	// Here we are generating it on client side only for example simplicity.
	claims := jwt.MapClaims{
		"sub": sub,
		"exp": jwt.NewNumericDate(time.Now().Add(exp)),
	}

	t, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(hmac_secret)

	if err != nil {
		return "", err
	}

	return t, nil
}

func ValidateJWT(tokenString string) (jwt.MapClaims, error) {
	return checkJWT(tokenString, viper.Get("JWT_SECRET_KEY").(string))
}

func checkJWT(tokenString string, secret string) (jwt.MapClaims, error) {
	var JWT_SIGNING_METHOD = jwt.SigningMethodHS256
	var JWT_SIGNATURE_KEY = []byte(secret)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method invalid")
		} else if method != JWT_SIGNING_METHOD {
			return nil, fmt.Errorf("signing method invalid")
		}

		return JWT_SIGNATURE_KEY, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, err
	}

	claims := token.Claims.(jwt.MapClaims)
	return claims, nil
}
