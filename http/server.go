package http

import (
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"hash"
	"net/http"

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
	DB             *bolt.DB
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

	//server
	server := &Server{
		StorageService: srv,
		Config:         c,
		DB:             b,
	}

	//create signer
	signer := gosigner.New(hmac.New(func() hash.Hash {
		return sha1.New()
	}, []byte(c.GetString("silo.signing.secret"))), gosigner.Options{
		CheckNonceFunc: server.CheckNonceFunc,
	})
	server.Signer = signer

	return server, nil
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

//function for the signer to check nonces and store them so a nonce can't be
//used twice, @TODO: clean the database once in a while
func (s *Server) CheckNonceFunc(n string) error {
	//create the database
	err := s.DB.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("nonces"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		nonce := b.Get([]byte(n))
		if nonce != nil {
			return fmt.Errorf("Nonce %s already exists", n)
		}
		return b.Put([]byte(n), []byte("1"))
	})
	return err
}
