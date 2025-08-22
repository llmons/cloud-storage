package mw

import (
	"cloud-storage/biz/dal/entity"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/jwt"
)

var JwtMiddleware *jwt.HertzJWTMiddleware

func InitJwt() {
	var err error
	JwtMiddleware, err = jwt.New(&jwt.HertzJWTMiddleware{
		Key: []byte("llmons"),
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if userBasic, ok := data.(entity.UserBasic); ok {
				return jwt.MapClaims{
					"id":       userBasic.ID,
					"identity": userBasic.Identity,
					"name":     userBasic.Name,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(ctx context.Context, c *app.RequestContext) interface{} {
			claims := jwt.ExtractClaims(ctx, c)
			return map[string]interface{}{
				"id":       claims["id"],
				"identity": claims["identity"],
				"name":     claims["name"],
			}
		},
	})
	if err != nil {
		panic(err)
	}
}
