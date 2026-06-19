package wmApiWeb

import (
	"context"
	"net/http"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/ggicci/httpin"
	"github.com/ggicci/httpin/integration"
	"github.com/go-chi/chi/v5"
	"github.com/unrolled/render"
)

var rnr = render.New(render.Options{
	StreamingJSON: true,
	UnEscapeHTML:  true,
})

const (
	GetWmApiUrl = "/wm_api/"
)

type getWmApiRequest struct {
	Feed      string `in:"query=feed" required:"true"`
	GroupBy   string `in:"query=group_by" required:"true"`
	DateStart string `in:"query=date_start" required:"true"`
	DateEnd   string `in:"query=date_end" required:"true"`
}

type getWmApiResponse struct {
	GroupKey    string  `json:"group_key"`
	Impressions uint64  `json:"impressions"`
	Clicks      uint64  `json:"clicks"`
	Cost        float64 `json:"cost"`
}

func InitHttpRoutes(
	ctx context.Context,
	httpRouter *chi.Mux,
	clickhouseConn driver.Conn,
	sspFeedsPopAdl map[string]string,
	sspFeedsPopMc map[string]string,
	sspFeedsIppAdl map[string]string,
	sspFeedsIppMc map[string]string,
	sspFeedsBanAdl map[string]string,
	sspFeedsBanMc map[string]string,
	sspFeedsNatAdl map[string]string,
	sspFeedsNatMc map[string]string,
) {
	integration.UseGochiURLParam("path", chi.URLParam)

	sspFeeds := mergeSspFeeds(
		sspFeedsPopAdl,
		sspFeedsPopMc,
		sspFeedsIppAdl,
		sspFeedsIppMc,
		sspFeedsBanAdl,
		sspFeedsBanMc,
		sspFeedsNatAdl,
		sspFeedsNatMc,
	)

	httpRouter.With(
		httpin.NewInput(getWmApiRequest{}),
	).Get(GetWmApiUrl, func(w http.ResponseWriter, r *http.Request) {
		getWmApiStats(ctx, w, r, clickhouseConn, sspFeeds)
	})
}
