package token

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

const accessTokenDuration = time.Minute * 15
const refreshTokenDuration = time.Hour * 24 * 7

func InitJWTKey() {
	jwtKey = []byte(os.Getenv("JWT_KEY"))
}

var jwtKey []byte

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

func GetNewTokens(userID uint) (string, string, error) {
	// TODO: Вынести TTL в конфиг
	accessToken, err := NewToken(userID, accessTokenDuration)
	if err != nil {
		log.Fatal(err)
		return "", "", nil
	}

	refreshToken, err := NewToken(userID, refreshTokenDuration)
	if err != nil {
		log.Fatal(err)
		return "", "", nil
	}

	return accessToken, refreshToken, nil
}

func NewToken(userID uint, ttl time.Duration) (string, error) {
	expirationTime := time.Now().Add(ttl)
	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func UpdateToken(refreshTokenString string) (string, string, error) {
	claims := &Claims{}
	refreshToken, err := jwt.ParseWithClaims(refreshTokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return "", "", err
	}

	if !refreshToken.Valid {
		return "", "", errors.New("invalid refresh token")
	}

	return GetNewTokens(claims.UserID)
}

func ValidateToken(refreshToken string) bool {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return false
	}
	return true
}

func VerifyToken(tokenString string) (uint, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Проверяем, что метод подписи совпадает с ожидаемым
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userIDFloat, ok := claims["user_id"].(float64); ok {
			userID := uint(userIDFloat)
			return userID, nil
		}
		return 0, fmt.Errorf("user_id not found in token claims")
	}

	return 0, fmt.Errorf("invalid token")
}

//func ParseToken(r *http.Request) (uint, error) {
//	authHeader := r.Header.Get("Authorization")
//
//	if authHeader == "" {
//		return 0, errors.New("Empty token")
//	}
//
//	parts := strings.Split(authHeader, " ")
//	if len(parts) != 2 || parts[0] != "Bearer" {
//		return 0, errors.New("Invalid Authorization header format")
//	}
//
//	accessToken := parts[1]
//
//	userID, err := VerifyToken(accessToken)
//	if err != nil {
//		return 0, err
//	}
//
//	return userID, nil
//}

func ParseToken(accessToken string) (uint, error) {
	//authHeader := r.Header.Get("Authorization")
	//
	//if authHeader == "" {
	//	return 0, errors.New("Empty token")
	//}
	//
	//parts := strings.Split(authHeader, " ")
	//if len(parts) != 2 || parts[0] != "Bearer" {
	//	return 0, errors.New("Invalid Authorization header format")
	//}
	//
	//accessToken := parts[1]

	userID, err := VerifyToken(accessToken)
	if err != nil {
		return 0, err
	}

	return userID, nil
}
