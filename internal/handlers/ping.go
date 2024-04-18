package handlers

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func MakePingHandler(dDSN string, logger *zap.SugaredLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db, err := sql.Open("postgres", dDSN)
		if err != nil {
			logger.Infof("Failed to connect to the database: %v", err)
			return
		}
		defer db.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = db.PingContext(ctx)
		if err != nil {
			logger.Infof("Unable to ping database: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
