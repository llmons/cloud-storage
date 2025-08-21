package jwt

import gjwt "github.com/golang-jwt/jwt/v5"

type UserClaims struct {
	ID       uint64
	Identity string
	Name     string
	gjwt.RegisteredClaims
}

var Key = []byte("llmons")
