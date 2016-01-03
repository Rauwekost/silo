package http

import (
	"crypto/hmac"
	"crypto/sha1"
	"hash"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/bmizerany/pat"
	"github.com/boltdb/bolt"
	"github.com/jacobstr/confer"
	"github.com/justinas/alice"
	"github.com/rauwekost/gosigner"
	"github.com/rauwekost/silo/service"
	"github.com/rauwekost/silo/storage"
	"github.com/rauwekost/silo/store"
)

var (
	httpPathHealth   = "/health"
	httpPathDisplay  = "/display"
	httpPathDownload = "/download"
	httpPathUpload   = "/upload"
)

type Server struct {
	StorageService *service.Storage
	Config         *confer.Config
	Signer         *gosigner.Signer
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

	//create signer
	logrus.Infof("Secret: %s", c.GetString("silo.signing.secret"))
	signer := gosigner.New(hmac.New(func() hash.Hash {
		return sha1.New()
	}, []byte(c.GetString("silo.signing.secret"))), gosigner.Options{})

	return &Server{
		StorageService: srv,
		Config:         c,
		Signer:         signer,
	}, nil
}

func (s *Server) HTTPHandler() http.Handler {
	mux := pat.New()
	chain := alice.New(routeLogger)
	secChain := alice.New(routeLogger, gosigner.NewMidddleware(s.Signer, nil).Handler)

	//routing
	mux.Add("GET", httpPathHealth, chain.Then(HttpHandler(s.healthHandler)))
	mux.Add("GET", httpPathDisplay, secChain.Then(HttpHandler(s.displayHandler)))
	mux.Add("GET", httpPathDownload, secChain.Then(HttpHandler(s.downloadHandler)))
	mux.Add("POST", httpPathUpload, secChain.Then(HttpHandler(s.uploadHandler)))

	return mux
}
