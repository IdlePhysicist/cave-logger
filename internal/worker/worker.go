package logs

import (
	"net/http"


	"github.com/go-chi/chi"

	"github.com/idlephysicist/cave-logger/pkg/keeper"
	"github.com/idlephysicist/cave-logger/pkg/model"
)

type Worker struct {
	log *logrus.Logger
	req *http.Request
	db	*keeper.Keeper
}

func NewWorker(log *logrus.Logger, r *http.Request) *Worker {
	return &Worker{log:log, req:r}
}

func (w *Worker) FigureIt() {
	/* Well 
		if Method == POST: it's a new log entry
		else they MUST have a logID
	*/
	logID := chi.URLParam(r, "logID")
	if logID == "" {
		if w.req.Method == "GET" {
			w.log.Debugf("worker.figureit: Calling ListLogs")
		} else {
			w.log.Debugf("worker.figureit: Calling Create")
			return Create()
		}
	} else {
		switch w.req.Method {
		case "GET":
			w.log.WithFields(logrus.Fields{"logID":logID}).Debugf("worker.figureit: Calling Get")
			return Get()
		case "PUT":
			w.log.WithFields(logrus.Fields{"logID":logID}).Debugf("worker.figureit: Calling Update")
			return Update()
		case "DELETE":
			w.log.WithFields(logrus.Fields{"logID":logID}).Debugf("worker.figureit: Calling Delete")
			return Delete()
		}
	}
}

func (w *Worker) ListLogs() http.HandlerFunc {
  return func (w http.ResponseWriter, r *http.Request) {
    err := render.RenderList(w, r, NewLogListResponse(Logs))
    if err != nil {
      render.Render(w, r, ErrRender(err))
      return
	  }
  }
}

func (w *Worker) Create() (error) {
		// So the request body will contain the log data
		var data *model.Row
		err := json.Unmarshal(w.req.Body, &data)
		if err != nil {
			w.log.Errorf("worker.create.json: Failed to unmarshal request body", err)
		}

		// TODO: wuts all this here ?
		if err = render.Bind(r, data); err != nil {
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}

		err = w.db.Insert(data)

		//render.Status(r, http.StatusCreated)
		//render.Render(w, r, NewLogResponse(Log))
	}
}

// GetLog returns the specific Log
func (w *Worker) Get(w http.ResponseWriter, r *http.Request) {
	// Assume if we've reach this far, we can access the Log
	// context because this handler is a child of the LogCtx
	// middleware. The worst case, the recoverer middleware will save us.
	Log := r.Context().Value("Log").(*Log)

	if err := render.Render(w, r, NewLogResponse(Log)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

// UpdateLog updates an existing Log in our persistent store.
func (w *Worker) Update(w http.ResponseWriter, r *http.Request) {
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
func (w *Worker) Delete(w http.ResponseWriter, r *http.Request) {
	var err error

	//Log := r.Context().Value("Log").(*Log) // REVIEW: ?
	err = w.db.DeleteRow(logID)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.Render(w, r, NewLogResponse(Log))
}
