package middleware

import (
	"log"
	"net/http"
	"os"
	"time"
)

var logFile *os.File

func InitLogFile() {
	var err error
	logFile, err = os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Error opening log file:", err)
	}
	log.SetOutput(logFile)
}

func LoggingToFileMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		log.Printf("[%s] %s %s %s", time.Now().Format("2006-01-02 15:04:05"), r.Method, r.RequestURI, time.Since(start))
	})
}
