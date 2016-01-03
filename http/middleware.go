package http

import (
	"net/http"
	"time"

	"github.com/rauwekost/silo/Godeps/_workspace/src/github.com/Sirupsen/logrus"
)

//routeLogger logs the route and request time
func routeLogger(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		next.ServeHTTP(w, r)
		t2 := time.Now()
		logrus.Printf("[%s] %q %v", r.Method, r.URL.String(), t2.Sub(t1))
	}

	return http.HandlerFunc(fn)
}
