package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"backend/internal/database"
	"backend/internal/service"
)

type Server struct {
	port int

	db                    database.Service
	dropboxRefreshService *service.DropboxRefreshService
}

func NewServer() *Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	server := &Server{
		port: port,
		db:   database.New(),
	}

	return server
}

func (s *Server) GetHTTPServer() *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf(":%d", s.port),
		Handler:      s.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}

func (s *Server) StopDropboxRefreshService() {
	if s.dropboxRefreshService != nil {
		s.dropboxRefreshService.Stop()
	}
}
