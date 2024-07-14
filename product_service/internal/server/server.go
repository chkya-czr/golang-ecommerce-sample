package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"product_service/config"
	"product_service/external/database"
	"product_service/external/validate"
	"product_service/internal/middleware"
	"product_service/logger"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/rs/cors"
	"gorm.io/gorm"
)

type Server struct {
	Version string
	cfg     *config.Config

	db   *sql.DB
	gorm *gorm.DB

	validator *validator.Validate
	cors      *cors.Cors
	router    *chi.Mux

	httpServer *http.Server
}

func (s *Server) Config() *config.Config {
	return s.cfg
}

func New() *Server {
	return &Server{
		cfg:    config.New(),
		router: chi.NewRouter(),
	}

}

func (s *Server) initLog() {
	slog.SetDefault(slog.New(logger.NewTraceHandler(
		os.Stdout,
		&slog.HandlerOptions{},
	)))
}

func (s *Server) setCors() {
	s.cors = cors.New(
		cors.Options{
			AllowedOrigins: s.cfg.Cors.AllowedOrigins,
			AllowedMethods: []string{
				http.MethodHead,
				http.MethodGet,
				http.MethodPost,
				http.MethodPut,
				http.MethodPatch,
				http.MethodDelete,
			},
			AllowedHeaders:   []string{"*"},
			AllowCredentials: true,
		})
}

func (s *Server) newValidator() {
	s.validator = validate.New()
}

func (s *Server) NewDatabase() {
	if s.cfg.Database.Driver == "" {
		log.Fatal("please fill in database credentials in .env file or set in environment variable")
	}

	s.gorm = database.NewGorm(s.cfg.Database)

	db, err := s.gorm.DB()

	if err != nil {
		log.Fatal("error while connecting to DB")
	}

	s.db = db
	s.db.SetMaxOpenConns(s.cfg.Database.MaxConnectionPool)
	s.db.SetMaxIdleConns(s.cfg.Database.MaxIdleConnections)
	s.db.SetConnMaxLifetime(s.cfg.Database.ConnectionsMaxLifeTime)
}

func (s *Server) Init() {
	s.initLog()
	s.setCors()
	s.NewDatabase()
	s.newValidator()
	s.newRouter()
	s.setGlobalMiddleware()
}

func (s *Server) newRouter() {
	s.router = chi.NewRouter()
}

func (s *Server) setGlobalMiddleware() {
	s.router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"message": "endpoint not found"}`))
	})
	s.router.Use(s.cors.Handler)
	s.router.Use(middleware.Json)
	s.router.Use(middleware.Audit)
	if s.cfg.Api.RequestLog {
		s.router.Use(chiMiddleware.Logger)
	}
	s.router.Use(chiMiddleware.Recoverer)
}

func (s *Server) Run() {
	s.httpServer = &http.Server{
		Addr:              s.cfg.Api.Host + ":" + s.cfg.Api.Port,
		Handler:           s.router,
		ReadHeaderTimeout: s.cfg.Api.ReadHeaderTimeout,
	}

	fmt.Println(`
	░██████╗░░█████╗░░░░░░░██╗░░██╗░█████╗░██████╗░████████╗
	██╔════╝░██╔══██╗░░░░░░██║░██╔╝██╔══██╗██╔══██╗╚══██╔══╝
	██║░░██╗░██║░░██║█████╗█████═╝░███████║██████╔╝░░░██║░░░
	██║░░╚██╗██║░░██║╚════╝██╔═██╗░██╔══██║██╔══██╗░░░██║░░░
	╚██████╔╝╚█████╔╝░░░░░░██║░╚██╗██║░░██║██║░░██║░░░██║░░░
	░╚═════╝░░╚════╝░░░░░░░╚═╝░░╚═╝╚═╝░░╚═╝╚═╝░░╚═╝░░░╚═╝░░░`)
	go func() {
		start(s)
	}()

	_ = gracefulShutdown(context.Background(), s)
}

func start(s *Server) {
	log.Printf("Serving at %s:%s\n", s.cfg.Api.Host, s.cfg.Api.Port)
	err := s.httpServer.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func gracefulShutdown(ctx context.Context, s *Server) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	log.Println("Shutting down...")

	ctx, shutdown := context.WithTimeout(ctx, s.Config().Api.GracefulTimeout*time.Second)
	defer shutdown()

	err := s.httpServer.Shutdown(ctx)
	if err != nil {
		log.Println(err)
	}

	return nil
}
