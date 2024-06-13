package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kylerequez/make-you-work-app/src/api/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetJWTKey() (*string, error) {
	key, err := GetEnv("JWT_KEY")
	if err != nil {
		return nil, err
	}
	return key, nil
}

func GenerateJWT(user *models.User) (*string, error) {
	key, err := GetJWTKey()
	if err != nil {
		return nil, err
	}
	signingKey := []byte(*key)

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["expiration"] = time.Now().Add(10 * time.Minute)
	claims["id"] = user.ID
	claims["username"] = user.Username
	claims["email"] = user.Email
	claims["authorities"] = user.Authorities

	jwtString, err := token.SignedString(signingKey)
	if err != nil {
		return nil, err
	}

	return &jwtString, nil
}

func ValidateJwt(tokenString string) (bool, error) {
	token, err := RetriveJwtToken(tokenString)
	if err != nil && token == nil {
		return false, err
	}

	claims, err := RetrieveClaims(token)
	if err != nil && claims == nil {
		return false, err
	}

	isExpired, err := IsExpired(claims)
	if isExpired {
		return false, err
	}

	return true, nil
}

func CheckIfUserHasAuthorities(tokenString string, authorities []string) (bool, error) {
	token, err := RetriveJwtToken(tokenString)
	if err != nil && token == nil {
		return false, err
	}

	claims, err := RetrieveClaims(token)
	if err != nil && claims == nil {
		return false, err
	}

	currentAuth := RetrieveAuthorities(claims)
	if currentAuth == nil {
		return false, errors.New("user does not have any authorities")
	}

	hasValue := HasValues(currentAuth, authorities)
	if !hasValue {
		return false, errors.New("user unauthorized")
	}
	return true, nil
}

func RetrieveAuthorities(claims jwt.MapClaims) []interface{} {
	return claims["authorities"].([]interface{})
}

func RetriveJwtToken(tokenString string) (*jwt.Token, error) {
	key, err := GetJWTKey()
	if err != nil {
		return nil, err
	}
	signingKey := []byte(*key)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("error in validating jwt")
		}
		return signingKey, nil
	})
	if err != nil {
		return nil, err
	}

	return token, nil
}

func RetrieveClaims(token *jwt.Token) (jwt.MapClaims, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("auth token is not valid")
	}

	return claims, nil
}

func IsExpired(claims jwt.MapClaims) (bool, error) {
	tokenExpiration, err := time.Parse(time.RFC3339, claims["expiration"].(string))
	if err != nil {
		return true, err
	}

	if time.Now().After(tokenExpiration) {
		return true, errors.New("auth token has expired")
	}

	return false, nil
}

func RetrieveUserIdFromJwt(tokenString string) (*primitive.ObjectID, error) {
	token, err := RetriveJwtToken(tokenString)
	if err != nil && token == nil {
		return nil, err
	}

	claims, err := RetrieveClaims(token)
	if err != nil && claims == nil {
		return nil, err
	}

	id := claims["id"].(string)
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	return &oid, nil
}
