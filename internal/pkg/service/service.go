package service

import (
  "net/http"
  "strings"

  "github.com/satori/uuid"
  "github.com/sirupsen/logrus"
  "github.com/gorilla/mux"

  "github.com/idlephysicist/cave-logger/internal/pkg/worker"
)

const dir = "./html"

type Service struct {
  log  *logrus.Logger
  port string
}


func NewService(log *logrus.Logger, port string) *Service {
  return &Service{log: log, port: port}
}

func (service *Service) Run() {
  service.log.Infof("main: Serving " + dir + " on http://localhost" + service.port)

  router := mux.NewRouter()

  router.HandleFunc("/healthz", service.HealthZ())
  sub_logs := router.PathPrefix("/logs/").Subrouter()
  sub_logs.HandleFunc("/", service.General()).Methods("GET")
  sub_logs.HandleFunc("/{key}", service.Detail()).Methods("GET", "PUT", "DELETE")
  sub_logs.HandleFunc("/new", service.NewEntry()).Methods("POST")

  router.HandleFunc("/", service.wasm()).Methods("GET")

  err := http.ListenAndServe(service.port, router)
  if err != nil {
    service.log.Fatalf("main.http: Error=%v", err)
  }
}

//
// Serve the website
func (service *Service) wasm() http.HandlerFunc {
  return func (resp http.ResponseWriter, req *http.Request) {
    fs := http.FileServer(http.Dir(dir))

    resp.Header().Add("Cache-Control", "no-cache")
    if strings.HasSuffix(req.URL.Path, ".wasm") {
      resp.Header().Set("content-type", "application/wasm")
    }
    
    fs.ServeHTTP(resp, req)
  }
}

//
// Get all entries
func (service *Service) General() http.HandlerFunc {
  return func (w http.ResponseWriter, r *http.Request) {

    uuid, err := uuid.NewV4()
    if err != nil {
      w.WriteHeader(http.StatusInternalServerError)
      return
    }

    log := service.log.WithFields(logrus.Fields{
      "id": uuid.String(), "ip": r.RemoteAddr, "method": r.Method,
    })

    wrkr := func () *worker.Worker {
      return worker.NewWorker(log, r)
    }()

    wrkr.ListEntries()
  }
}

//
// Get a specific entry
func (service *Service) Detail() http.HandlerFunc {
  return func (w http.ResponseWriter, r *http.Request) {

    uuid, err := uuid.NewV4()
    if err != nil {
      w.WriteHeader(http.StatusInternalServerError)
      return
    }

    log := service.log.WithFields(logrus.Fields{
      "id": uuid.String(), "ip": r.RemoteAddr, "method": r.Method,
    })

    wrkr := func () *worker.Worker {
      return worker.NewWorker(log, r)
    }()

    wrkr.Get()
  }
}

//
// Create a new entry
func (service *Service) NewEntry() http.HandlerFunc {
  return func (w http.ResponseWriter, r *http.Request) {
  }
}

func (service *Service) StatsHandler() http.HandlerFunc {
  return func (w http.ResponseWriter, r *http.Request) {
  }
}

func (service *Service) HealthZ() http.HandlerFunc {
  return func (w http.ResponseWriter, r *http.Request) {
  }
}
