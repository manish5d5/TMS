package https

import (
	"TMS/config"
	"context"
	"fmt"

	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	config config.Config
}

func NewServer(config config.Config) *Server {
	return &Server{
		config: config,
	}
}
func (s *Server) Listen(ctx context.Context, port string) error {

	r := chi.NewRouter()

	// // Global middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// API group
	r.Route("/api/v1", func(r chi.Router) {

	})
	addr := port
	server := &http.Server{
		Addr:    addr,
		Handler: r,
	}
	errCh := make(chan error, 1)

	go func() {
		fmt.Printf("🚀 Server running on %s\n", addr)
		errCh <- server.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		return err

	case <-ctx.Done():
		fmt.Println("\n🛑 Shutting down gracefully...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return server.Shutdown(shutdownCtx)
	}
}
