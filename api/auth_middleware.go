package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/oriventi/simplebank/token"
)

const (
	AuthorizationHeaderKey  = "authorization"
	AuthorizationPayloadKey = "username"
)

func authMiddleware(maker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader(AuthorizationHeaderKey)

		//check if header exists
		if len(authHeader) == 0 {
			noHeaderErr := errors.New("no Authorization-Header provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(noHeaderErr))
			return
		}

		splitAuthHeader := strings.Fields(authHeader)
		if len(splitAuthHeader) < 2 {
			wrongFormat := errors.New("wrong header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(wrongFormat))
			return
		}

		authType := splitAuthHeader[0]
		token := splitAuthHeader[1]

		if strings.ToLower(authType) != "bearer" {
			wrongTypeErr := errors.New("authorization type not supported")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(wrongTypeErr))
			return
		}

		payload, verifyErr := maker.VerifyToken(token)
		if verifyErr != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(verifyErr))
			return
		}

		ctx.Set(AuthorizationPayloadKey, payload)
		ctx.Next()
	}
}
