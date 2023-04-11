package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type renewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type renewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (server *Server) renewAccessToken(ctx *gin.Context) {
	var req renewAccessTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	refreshPayload, verifyErr := server.tokenMaker.VerifyToken(req.RefreshToken)
	if verifyErr != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(verifyErr))
		return
	}

	session, dbErr := server.store.GetSession(ctx, refreshPayload.ID)
	if dbErr != nil {
		if dbErr == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(dbErr))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(dbErr))
		return
	}

	if session.IsBlocked {
		err := fmt.Errorf("blocked sesion")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if session.ExpiresAt.Before(time.Now()) {
		err := fmt.Errorf("refresh token expired")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if session.Username != refreshPayload.Username {
		err := fmt.Errorf("session doesn't belong to user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	newAccessToken, newAccessTokenRes, newAccessTokenErr := server.tokenMaker.CreateToken(refreshPayload.Username, server.config.AccessTokenDuration)
	if newAccessTokenErr != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(newAccessTokenErr))
		return
	}

	response := renewAccessTokenResponse{
		AccessToken:          newAccessToken,
		AccessTokenExpiresAt: newAccessTokenRes.ExpiredAt,
	}
	ctx.JSON(http.StatusOK, response)
}
