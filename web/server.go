package web

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"github.com/greboid/ircclient/irc"
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var (
	//go:embed static/*
	staticFS embed.FS
)

type Server struct {
	lock              sync.Mutex
	httpServer        *http.Server
	connectionManager *irc.ConnectionManager
	activeServer      string
	activeWindow      string
	fixedPort         int
}

func NewServer(cm *irc.ConnectionManager, fixedPort int) *Server {
	mux := http.NewServeMux()
	server := &Server{
		fixedPort: fixedPort,
		lock:      sync.Mutex{},
		httpServer: &http.Server{
			Handler: mux,
		},
		connectionManager: cm,
	}
	server.addRoutes(mux)
	return server
}

func (s *Server) GetListenAddress() string {
	if s.httpServer.Addr != "" {
		return fmt.Sprintf("http://%s", s.httpServer.Addr)
	}
	ip, port, err := s.getPort()
	if err != nil {
		slog.Error("Unable to get free port", "error", err)
		return ""
	}
	s.httpServer.Addr = net.JoinHostPort(ip.String(), strconv.Itoa(port))
	return fmt.Sprintf("http://%s", s.httpServer.Addr)
}

func (s *Server) Start() string {
	clickAddr := s.GetListenAddress()
	slog.Info("Starting webserver", "url", clickAddr)
	if err := s.httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		slog.Error("Error starting server:", slog.String("error", err.Error()))
	}
	slog.Debug("Server stopped")
	return clickAddr
}

func (s *Server) Stop() {
	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 2*time.Second)
	defer shutdownRelease()
	if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
		slog.Error("Error shutting down server:", slog.String("error", err.Error()))
	}
}

func (s *Server) getPort() (net.IP, int, error) {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", s.fixedPort))
	if err != nil {
		return nil, -1, err
	}
	listen, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return nil, -1, err
	}
	defer func() { _ = listen.Close() }()
	lp := listen.Addr().(*net.TCPAddr)
	return lp.IP, lp.Port, nil
}
