package service

import (
  "net/http"
  "strings"

  "github.com/satori/uuid"
  "github.com/sirupsen/logrus"
  "github.com/go-chi/chi"
  "github.com/go-chi/chi/middleware"

  "github.com/idlephysicist/cave-logger/pkg/worker"
)

const dir = "./html"

type Service struct {
  log  *logrus.Logger
  port string
}


func NewService(log *logrus.Logger, port string) *Service {
  return &Service{log: log, port: port}
}

func (s *Service) Run() {
  fs := http.FileServer(http.Dir(dir))
  s.log.Infof("main: Serving " + dir + " on http://localhost" + s.port)

  router := chi.NewRouter()
  router.Use(middleware.Recoverer)

  router.Get("/healthz", s.HealthZ())

  router.Get("/logs/{logID}", s.LogsHandler())
  router.Put("/logs/{logID}", s.LogsHandler())
  router.Delete("/logs/{logID}", s.LogsHandler())
  router.Post("/logs/new", s.LogsHandler())
  
  router.Get("/stats/{raw}", s.StatsHandler())

  // RESTy routes for "articles" resource
	/*router.Route("/logs", func(router chi.Router) {
		router.With(paginate).Get("/", ListArticles)
		router.Post("/", CreateArticle)       // POST /articles
		router.Get("/search", SearchArticles) // GET /articles/search

		router.Route("/{logID}", func(router chi.Router) {
			router.Use(ArticleCtx)            // Load the *Article on the request context
			router.Get("/", GetArticle)       // GET /articles/123
			router.Put("/", UpdateArticle)    // PUT /articles/123
			router.Delete("/", DeleteArticle) // DELETE /articles/123
		})

		// GET /articles/whats-up
		//r.With(ArticleCtx).Get("/{articleSlug:[a-z-]+}", GetArticle)
	})*/


  // NOTE: I guess that the below should be a GET
  router.Get("/", func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Add("Cache-Control", "no-cache")
		if strings.HasSuffix(req.URL.Path, ".wasm") {
			resp.Header().Set("content-type", "application/wasm")
		}
		fs.ServeHTTP(resp, req)
  })

  err := http.ListenAndServe(s.port, router)
  if err != nil {
    s.log.Fatalf("main.http: Error=%v", err)
  }

}

func (s *Service) HealthZ() http.HandlerFunc {
  return func (w http.ResponseWriter, r *http.Request) {
  }
}

func (s *Service) LogHandler() http.HandlerFunc {
  return func (w http.ResponseWriter, r *http.Request) {

    uuid, err := uuid.NewV4()
    if err != nil {
      w.WriteHeader(http.StatusInternalServerError)
      return
    }

    log := s.log.WithFields(logrus.Fields{
      "id": uuid.String(), "ip": r.RemoteAddr, "method": r.Method, })

    w := func () *worker.Worker {
      worker.NewWorker(log, r)
    }()

    w.FigureIt()
  }
}

func (s *Service) StatsHandler() http.HandlerFunc {
  return func (w http.ResponseWriter, r *http.Request) {
  }
}
