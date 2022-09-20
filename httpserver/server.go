package httpserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/bratteby/go-service-template/httpserver/middleware"
	"github.com/bratteby/go-service-template/logging"
)

type Server struct {
	Address        string
	ExampleService exampleService
	Logger         *logging.Logger
}

func (s *Server) Start() error {

	r := s.setupHandler()

	// TODO: Setup server with timeouts etc.
	// s := http.Server {
	// 	Addr: s.Address,
	// 	Handler: r,
	//  timeouts...
	// }

	return http.ListenAndServe(s.Address, r)
}

// setupHandler will setup all routes and return the http handler.
func (s Server) setupHandler() http.Handler {
	r := chi.NewRouter()

	// Middlewares.
	r.Use(chimiddleware.RequestID)
	r.Use(middleware.RequestLogger(*s.Logger, &middleware.RequestLoggerOptions{
		Verbose: true,
	}))
	r.Use(chimiddleware.Recoverer)

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("healthy"))
	})

	e := encoder{
		Logger: s.Logger,
	}

	exampleHandler := exampleHandler{
		exampleService: s.ExampleService,
		encoder:        e,
	}

	r.Route("/api", func(r chi.Router) {
		r.Use(chimiddleware.BasicAuth("Example", map[string]string{
			"username": "nOt_saFE_PWD",
		}))

		r.Route("/example", exampleHandler.GetRoutes())
	})

	return r
}
