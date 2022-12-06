package auth

import (
	"context"
	"item/server/constants"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var accessRoles = map[string][]string{
	constants.METHOD_AUTH + "Login":      {constants.ROLE_USER, constants.ROLE_ADMIN},
	constants.METHOD_ITEM + "GetItems":   {constants.ROLE_USER, constants.ROLE_ADMIN},
	constants.METHOD_ITEM + "GetItem":    {constants.ROLE_USER, constants.ROLE_ADMIN},
	constants.METHOD_ITEM + "CreateItem": {constants.ROLE_ADMIN},
	constants.METHOD_ITEM + "UpdateItem": {constants.ROLE_ADMIN},
	constants.METHOD_ITEM + "DeleteItem": {constants.ROLE_USER, constants.ROLE_ADMIN},
}

func Authorize(ctx context.Context, method string) error {
	roles, ok := accessRoles[method]
	if !ok {
		log.Println("no permission to access this service")
		return status.Errorf(codes.PermissionDenied, "no permission to access this service")
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Println("metadata is not provided")
		return status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := md["authorization"]
	if len(values) == 0 {
		log.Println("authorization token is not provided")
		return status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	accessToken := values[0]
	claims, err := VerifyToken(accessToken)
	if err != nil {
		log.Println("access token is not valid")
		return status.Errorf(codes.Unauthenticated, "access token is not valid: %v", err)
	}

	for _, role := range roles {
		if role == claims.Role {
			return nil
		}
	}

	log.Println("THIRD")
	return status.Errorf(codes.PermissionDenied, "no permission to access this service")
}
