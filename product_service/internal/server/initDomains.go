package server

import (
	"net/http"
	"product_service/internal/domain/product/handler"
	"product_service/internal/domain/product/service"
	"product_service/internal/middleware"
	"product_service/utility/respond"

	"github.com/go-chi/chi/v5"
)

func (s *Server) InitDomains() {
	s.initVersion()
	// s.initSwagger()
	// s.initHealth()
	s.initProduct()

}

func (s *Server) initVersion() {
	s.router.Route("/version", func(router chi.Router) {
		router.Use(middleware.Json)

		router.Get("/", func(w http.ResponseWriter, r *http.Request) {
			respond.Json(w, http.StatusOK, map[string]string{"version": s.Version})
		})
	})
}

// func (s *Server) initSwagger() {
// 	if s.Config().Api.RunSwagger {
// 		docsPath, err := fs.Sub(swaggerDocsAssetPath, "docs")
// 		if err != nil {
// 			panic(err)
// 		}

// 		fileServer := http.FileServer(http.FS(docsPath))

// 		s.router.HandleFunc("/swagger", func(w http.ResponseWriter, r *http.Request) {
// 			http.Redirect(w, r, "/swagger/", http.StatusMovedPermanently)
// 		})
// 		s.router.Handle("/swagger/", http.StripPrefix("/swagger", middleware.ContentType(fileServer)))
// 		s.router.Handle("/swagger/*", http.StripPrefix("/swagger", middleware.ContentType(fileServer)))
// 	}
// }

func (s *Server) initProduct() {
	// newBookRepo := bookRepo.New(s.sqlx)
	productService := service.New()
	handler.RegisterHTTPEndPoints(s.router, s.validator, productService)
}
