package middlewares

import (
	"errors"
	"fmt"
	"github.com/dimasbagussusilo/gin-golang-boilerplate/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	AuthorizationHeaderKey  = "authorization"
	AuthorizationTypeBearer = "bearer"
	AuthorizationPayloadKey = "authorization_payload"
)

// AuthMiddleware creates a gin middleware for authorization
func AuthMiddleware(tokenMaker utils.TokenMaker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(AuthorizationHeaderKey)

		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")

			ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.ResponseData(utils.ResponseStatusError, err.Error(), nil))
			return
		}

		var accessToken string
		fields := strings.Fields(authorizationHeader)
		//if len(fields) < 2 {
		//	err := errors.New("invalid authorization header format")
		//	ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.ResponseData(utils.ResponseStatusError, err.Error(), nil))
		//	 return
		//}

		if len(fields) > 1 {
			authorizationType := strings.ToLower(fields[0])
			if authorizationType != AuthorizationTypeBearer {
				err := fmt.Errorf("unsupported authorization type %s", authorizationType)
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.ResponseData(utils.ResponseStatusError, err.Error(), nil))
				return
			}

			accessToken = fields[1]
		}

		if len(fields) == 1 {
			accessToken = fields[0]
		}

		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.ResponseData(utils.ResponseStatusError, err.Error(), nil))
			return
		}

		ctx.Set(AuthorizationPayloadKey, payload)
		ctx.Next()
	}
}
