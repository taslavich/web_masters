package wmApiWeb

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/ggicci/httpin"
)

func getWmApiStats(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
	clickhouseConn driver.Conn,
	sspFeeds map[string]string,
) {
	input := r.Context().Value(httpin.Input).(*getWmApiRequest)

	sspDomain, ok := sspFeeds[input.Feed]
	if !ok {
		err := fmt.Errorf("Busy")
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	groupByExpr, err := getWmApiGroupByExpression(input.GroupBy)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dateStart, err := time.Parse("2006-01-02", input.DateStart)
	if err != nil {
		http.Error(w, "Invalid date_start", http.StatusBadRequest)
		return
	}

	dateEnd, err := time.Parse("2006-01-02", input.DateEnd)
	if err != nil {
		http.Error(w, "Invalid date_end", http.StatusBadRequest)
		return
	}

	if dateStart.After(dateEnd) {
		http.Error(w, "date_start cannot be greater than date_end", http.StatusBadRequest)
		return
	}

	res, err := selectWmApiStats(ctx, clickhouseConn, groupByExpr, sspDomain, dateStart, dateEnd)
	if err != nil {
		log.Printf("Cannot select wm api stats: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := rnr.JSON(w, http.StatusOK, res); err != nil {
		log.Printf("Cannot make HTTP response back: %v\n", err)
	}
}
