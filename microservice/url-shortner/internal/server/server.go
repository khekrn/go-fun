package server

import (
	"coding2fun.in/url-shortner/internal/config"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Server struct {
	router *gin.Engine
	config *config.Config
	db     *gorm.DB
}

func NewServer(config *config.Config, db *gorm.DB) *Server {
	gin.SetMode(config.Server.Mode)
	server := &Server{
		router: gin.Default(),
		config: config,
		db:     db,
	}
	return server
}

func (s *Server) setUp() {
	s.router.GET("/health", s.defaultHandler)
}

func (s *Server) defaultHandler(ctx *gin.Context) {}

func (s *Server) Run() {

}
