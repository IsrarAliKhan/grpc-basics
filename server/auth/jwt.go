package auth

import (
	"encoding/json"
	"item/server/config"
	"item/server/db"
	"item/server/db/models"
	"item/server/log"
	"time"

	"github.com/go-chi/jwtauth"
	"github.com/golang-jwt/jwt"
	jwtclaim "github.com/golang-jwt/jwt/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var TokenAuth *jwtauth.JWTAuth

func init() {
	TokenAuth = jwtauth.New("HS256", []byte(config.JwtSecretKey), nil)
}

type JWT struct {
	Username         string `json:"username"`
	Role             string `json:"role"`
	RegisteredClaims jwtclaim.RegisteredClaims
}

func Generate(username string, role string) (string, error) {
	// generate jwt token
	_, tokenString, err := TokenAuth.Encode(jwt.MapClaims{
		"user": JWT{
			Username: username,
			Role:     role,
			RegisteredClaims: jwtclaim.RegisteredClaims{
				ExpiresAt: jwtclaim.NewNumericDate(time.Now().Add(30 * time.Minute)),
			},
		},
	})
	if err != nil {
		return "", status.Errorf(codes.Internal, err.Error())
	}

	// respond
	return tokenString, nil
}

func VerifyToken(tokenString string) (*JWT, error) {
	var loggedUser JWT
	token, err := TokenAuth.Decode(tokenString)
	if err != nil {
		log.Error.Println(err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	jsonString, err := json.Marshal(token.PrivateClaims()["user"])
	if err != nil {
		log.Error.Println(err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	err = json.Unmarshal(jsonString, &loggedUser)
	if err != nil {
		log.Error.Println(err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	var user models.Users
	err = db.Conn().
		Where("username=? AND role=?", loggedUser.Username, loggedUser.Role).
		First(&user).
		Error
	if err != nil {
		log.Error.Println(err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	if loggedUser.RegisteredClaims.ExpiresAt.Before(time.Now()) {
		log.Error.Println("token expired")
		return nil, status.Errorf(codes.Unauthenticated, "token expired")
	}

	return &JWT{
		Username:         loggedUser.Username,
		Role:             loggedUser.Role,
		RegisteredClaims: loggedUser.RegisteredClaims,
	}, nil
}
