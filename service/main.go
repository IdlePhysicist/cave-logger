package main

import (
  "net/http"

  "github.com/sirupsen/logrus"
  "github.com/gorilla/mux"
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

  router := mux.NewRouter()

  router.HandleFunc("/healthz", s.HealthZ()).Methods("GET")
  router.HandleFunc("/logs/{event}", s.ViewLog).Methods("GET")
  router.HandleFunc("/logs", s.CreateLog).Methods("POST")
  router.HandleFunc("/logs", s.UpdateLog).Methods("PUT")
  router.HandleFunc("/stats/{raw}", Stats).Methods("GET")
  router.PathPrefix("/").Handler(func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Add("Cache-Control", "no-cache")
		if strings.HasSuffix(req.URL.Path, ".wasm") {
			resp.Header().Set("content-type", "application/wasm")
		}
		fs.ServeHTTP(resp, req)
  }

  err := http.ListenAndServe(port, router)
  if err != nil {
    log.Fatalf("main.http: Error=%v", err)
  }

}

func (s *Service) HealthZ() http.HandlerFunc {
  return func (w http.ResponseWriter, r *http.Request) {
  }
}


