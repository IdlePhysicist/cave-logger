package worker

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
	//"github.com/go-chi/chi"

	"github.com/idlephysicist/cave-logger/internal/pkg/keeper"
	//"github.com/idlephysicist/cave-logger/internal/pkg/model"
)

type Worker struct {
	log *logrus.Logger
	req *http.Request
	db	*keeper.Keeper
}

func NewWorker(log *logrus.Logger, r *http.Request) *Worker {
	return &Worker{log:log, req:r}
}

/*func (worker *Worker) FigureIt()  {
  return func (w http.ResponseWriter, r *http.Request) {
		/* Well 
			if Method == POST: it's a new log entry
			else they MUST have a logID
		* /
		logID := chi.URLParam(worker.req, "logID")
		if logID == "" {
			if worker.req.Method == "GET" {
				worker.log.Debugf("worker.figureit: Calling ListLogs")
				
			} else {
				worker.log.Debugf("worker.figureit: Calling Create")
				return worker.Create(w,r)
			}
		} else {
			switch worker.req.Method {
			case "GET":
				worker.log.WithField("logID", logID).Debugf("worker.figureit: Calling Get")
				return worker.Get(w,r)
			case "PUT":
				worker.log.WithField("logID", logID).Debugf("worker.figureit: Calling Update")
				return worker.Update(w,r)
			case "DELETE":
				w.log.WithField("logID", logID).Debugf("worker.figureit: Calling Delete")
				return worker.Delete(w,r)
			}
		}
	}
}*/

func (worker *Worker) ListEntries() http.HandlerFunc {
  return func (w http.ResponseWriter, r *http.Request) {
    /*err := render.RenderList(w, r, NewLogListResponse(Logs))
    if err != nil {
      render.Render(w, r, ErrRender(err))
      return
		}*/
		

		data, err := worker.db.QueryLogs(`-1`)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)

		return
  }
}

/*func (worker *Worker) Create() (error) {
	// So the request body will contain the log data
	var data *model.Row
	err := json.Unmarshal(w.req.Body, &data)
	if err != nil {
		worker.log.Errorf("worker.create.json: Failed to unmarshal request body", err)
	}

	// TODO: wuts all this here ?
	if err = render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	err = worker.db.Insert(data)

	//render.Status(r, http.StatusCreated)
	//render.Render(w, r, NewLogResponse(Log))
}*/

// GetLog returns the specific Log
func (worker *Worker) Get() http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		// Assume if we've reach this far, we can access the Log
		// context because this handler is a child of the LogCtx
		// middleware. The worst case, the recoverer middleware will save us.
		//Log := r.Context().Value("Log").(*Log)
		key := r.FormValue(`key`)
		if key == `` {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		/*if err := render.Render(w, r, NewLogResponse(Log)); err != nil {
			render.Render(w, r, ErrRender(err))
			return
		}*/

		data, err := worker.db.QueryLogs(key)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)

		return
	}
}

// UpdateLog updates an existing Log in our persistent store.
/*func (worker *Worker) Update(w http.ResponseWriter, r *http.Request) {
	Log := r.Context().Value("Log").(*Log)

	data := &LogRequest{Log: Log}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	Log = data.Log
	dbUpdateLog(Log.ID, Log)

	render.Render(w, r, NewLogResponse(Log))
}

// DeleteLog removes an existing Log from our persistent store.
func (worker *Worker) Delete(w http.ResponseWriter, r *http.Request) {
	var err error

	//Log := r.Context().Value("Log").(*Log) // REVIEW: ?
	err = worker.db.DeleteRow(logID)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.Render(w, r, NewLogResponse(Log))
}
*/