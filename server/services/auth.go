package services

import (
	"context"
	"item/pb"
	jwt "item/server/auth"
	"item/server/db"
	"item/server/db/models"
	"item/server/log"
)

type AuthServer struct {
	pb.UnimplementedAuthServer
}

func (s *AuthServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	log.Info.Printf("Recieved: %v\n", req)

	// get user form db
	var user models.Users
	err := db.Conn().
		Where("username=? AND password=?", req.Username, req.Password).
		First(&user).
		Error
	if err != nil {
		log.Error.Println(err)
		return nil, err
	}

	// get jwt token
	token, err := jwt.Generate(req.Username, "ADMIN")
	if err != nil {
		log.Error.Println(err)
		return nil, err
	}

	// build response
	res := &pb.LoginResponse{
		AccessToken: token,
	}

	// respond
	return res, nil
}
