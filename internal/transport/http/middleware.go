package http

import (
	"github.com/gin-gonic/gin"
	"github.com/stepanbukhtii/easy-tools/crypto"
	easymiddleware "github.com/stepanbukhtii/easy-tools/rest/middleware"
)

func (s *Server) registerMiddlewares(router *gin.Engine) {}

func (s *Server) registerJWTMiddleware() error {
	ed25519PublicKey, err := crypto.ParseED25519PublicKey(s.config.JWT.PublicKey)
	if err != nil {
		return err
	}

	s.jwtMiddleware = easymiddleware.NewJWTAuth(
		ed25519PublicKey,
		s.config.JWT.Issuer,
		s.config.JWT.Audience,
		s.config.JWT.Enabled,
	)

	return nil
}
