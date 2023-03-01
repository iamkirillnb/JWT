package auth

import (
	"auth-service/sdk/auth/proto"
)

type StoreAuth interface {
}

type service struct {
	commit, buildAt, version, jwtSK string
	tokenExp, refreshTokenExp       int64
	store                           StoreAuth
	proto.UnimplementedAuthServer
}

func NewService(commit, buildAt, version, jwtSK string, tokenExp, refreshTokenExp int64) *service {
	return &service{
		commit:          commit,
		buildAt:         buildAt,
		version:         version,
		tokenExp:        tokenExp,
		refreshTokenExp: refreshTokenExp,
		jwtSK:           jwtSK,
	}
}

func (service *service) SetStore(store StoreAuth) *service {
	service.store = store
	return service
}
