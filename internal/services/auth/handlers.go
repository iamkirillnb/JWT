package auth

import (
	"auth-service/sdk/auth/proto"
	"context"
	"fmt"
	"github.com/golang-jwt/jwt"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"os"
	"strings"
	"time"
)

type Claims struct {
	Username string `json:"username"`
	jwt.MapClaims
}

func (s *service) Authentication(ctx context.Context, input *proto.AuthenticationRequest) (*proto.AuthenticationResponse, error) {

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = jwt.MapClaims{
		"exp": time.Now().Unix() + s.tokenExp,
		"sub": input.GetUuid(),
	}
	token.Header["kid"] = "signin_1"
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshToken.Claims = jwt.MapClaims{
		"exp": time.Now().Unix() + s.refreshTokenExp,
		"sub": input.GetUuid(),
	}
	refreshToken.Header["kid"] = "signin_2"

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(s.jwtSK))
	if err != nil {
		log.Println(err)
	}
	refreshTokenString, err := refreshToken.SignedString([]byte(s.jwtSK))
	if err != nil {
		log.Println(err)
	}

	return &proto.AuthenticationResponse{
		Token:        tokenString,
		RefreshToken: refreshTokenString,
	}, nil

}

func (s *service) ValidateToken(ctx context.Context, input *proto.ValidateTokenRequest) (*proto.ValidateTokenResponse, error) {

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(input.Token, claims, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Println(ok)
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(s.jwtSK), nil

	})

	if val, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if token.Header["kid"] != nil && token.Header["kid"] == "signin_1" {
			isExp := val.VerifyExpiresAt(int64(val["exp"].(float64)), false)
			if isExp {
				return &proto.ValidateTokenResponse{Valid: true}, err
			}
			return &proto.ValidateTokenResponse{Valid: false}, nil
		}
	}
	if err != nil {
		log.Println(err)
		return &proto.ValidateTokenResponse{Valid: false}, err
	}

	return &proto.ValidateTokenResponse{Valid: false}, nil
}

func (s *service) RefreshToken(ctx context.Context, input *emptypb.Empty) (*proto.RefreshTokenResponse, error) {
	var (
		authHeader []string
		tkn        = ""
	)
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		authHeader = md.Get("Authorization")
	}
	log.Printf("authHeader= %v", authHeader)
	if authHeader == nil {
		return nil, fmt.Errorf("no %v token in header", "Authorization")
	}
	if len(authHeader) > 0 {
		tkn = strings.Split(authHeader[0], " ")[1]
	}

	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(tkn, claims, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Println(ok)
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		if sk, ok := os.LookupEnv("JWT_SECRET"); ok {
			return []byte(sk), nil
		}
		return []byte(""), nil

	})

	if val, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Create a new token object, specifying signing method and the claims
		// you would like it to contain.
		tokenNew := jwt.New(jwt.SigningMethodHS256)
		tokenNew.Claims = jwt.MapClaims{
			"exp": time.Now().Unix() + s.tokenExp,
			"sub": val["uuid"],
		}
		tokenNew.Header["kid"] = "signin_1"
		refreshTokenNew := jwt.New(jwt.SigningMethodHS256)
		refreshTokenNew.Claims = jwt.MapClaims{
			"exp": time.Now().Unix() + s.refreshTokenExp,
			"sub": val["uuid"],
		}
		refreshTokenNew.Header["kid"] = "signin_2"

		// Sign and get the complete encoded token as a string using the secret
		tokenString, err := tokenNew.SignedString([]byte(s.jwtSK))
		if err != nil {
			log.Println(err)
		}
		refreshTokenString, err := refreshTokenNew.SignedString([]byte(s.jwtSK))
		if err != nil {
			log.Println(err)
		}
		return &proto.RefreshTokenResponse{
			Token:        tokenString,
			RefreshToken: refreshTokenString,
		}, nil

	}
	if err != nil {
		log.Println(err)
		return &proto.RefreshTokenResponse{}, err
	}

	return &proto.RefreshTokenResponse{}, nil

}
