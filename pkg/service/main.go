package service

import (
  "net/http"

  "github.com/satori/uuid"
  "github.com/sirupsen/logrus"
  "github.com/go-chi/chi"
)

const dir = "./html"

type Service struct {
  log  *logrus.Logger
  port string
}




func NewService(log *logrus.Logger, port string) *Service {
  return &Service{log:log, port: port}
}

func (s *Service) Run() {

  fs := http.FileServer(http.Dir(dir))
  log.Infof("main: Serving " + dir + " on http://localhost" + port)

  router := chi.NewRouter()

  // TODO: Fix all these lads
  router.Get("/healthz", s.HealthZ()).Methods("GET")
  router.HandleFunc("/logs/{event}", s.LogsHandler).Methods("GET", "PUT")
  router.HandleFunc("/logs", s.LogsHandler).Methods("POST")
  router.HandleFunc("/stats/{raw}", s.StatsHandler).Methods("GET")
  router.PathPrefix("/").Handler(func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Add("Cache-Control", "no-cache")
		if strings.HasSuffix(req.URL.Path, ".wasm") {
			resp.Header().Set("content-type", "application/wasm")
		}
		fs.ServeHTTP(resp, req)
  })

  err := http.ListenAndServe(port, router)
  if err != nil {
    log.Fatalf("main.http: Error=%v", err)
  }

}

func (s *Service) HealthZ() http.HandlerFunc {
  return func (w http.ResponseWriter, r *http.Request) {
  }
}

func (s *Service) LogsHandler() http.HandlerFunc {
  return func (w. http.ResponseWriter, r *http.Request) {

    uuid, err := uuid.NewV4()
    if err != nil {
      w.WriteHeader(http.StatusInternalServerError)
      return
    }

    if (r.Method == "GET" || r.Method == "PUT") {
      event := chi.URLParam(r, "event")
      if event == "" {
        w.WriteHeader(http.StatusBadRequest)
        return
      }
    }



    log := service.log.WithFields(logrus.Fields{
      "id": uuid.String(), "ip": r.RemoteAddr, "method": r.Method, })


  }
}

