package http

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/boltdb/bolt"
	"github.com/jacobstr/confer"
	"github.com/justinas/alice"
	"github.com/rauwekost/silo/service"
	"github.com/rauwekost/silo/storage"
	"github.com/rauwekost/silo/store"
)

var (
	httpPathHealth   = "/health"
	httpPathDisplay  = "/display/:sum"
	httpPathDownload = "/download/:sum"
	httpPathUpload   = "/upload"
)

type Server struct {
	StorageService *service.Storage
	Config         *confer.Config
}

func NewServer(c *confer.Config) (*Server, error) {
	//create new file storage
	fs := storage.NewFileSystemStorage(c.GetString("silo.storage.location"))

	//create new bolt database
	b, err := bolt.Open(c.GetString("silo.database.location"), 0600, nil)
	if err != nil {
		return nil, err
	}

	//create database json wrapper
	db := store.NewJSONStore(b, []byte("metadata"))

	//create new service wrapper
	srv, err := service.New(fs, db)
	if err != nil {
		return nil, err
	}

	return &Server{
		StorageService: srv,
		Config:         c,
	}, nil
}

func (s *Server) HTTPHandler() http.Handler {
	mux := pat.New()
	chain := alice.New(routeLogger)

	//routing
	mux.Add("GET", httpPathHealth, chain.Then(HttpHandler(s.healthHandler)))
	mux.Add("GET", httpPathDisplay, chain.Then(HttpHandler(s.displayHandler)))
	mux.Add("GET", httpPathDownload, chain.Then(HttpHandler(s.downloadHandler)))
	mux.Add("POST", httpPathUpload, chain.Then(HttpHandler(s.uploadHandler)))

	return mux
}
