package wmApiWeb

import (
	"context"
	"fmt"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

func selectWmApiStats(
	ctx context.Context,
	clickhouseConn driver.Conn,
	groupByExpr string,
	sspDomain string,
	dateStart time.Time,
	dateEnd time.Time,
) ([]getWmApiResponse, error) {
	query := fmt.Sprintf(`
WITH
clicks_part AS
(
    SELECT
        toString(%[1]s) as gr_key,
        count(*) AS clicks,
        sumIf(win_dsp_price, format = 'POP') / 1000
            + sumIf(win_dsp_price, format = 'IPP') AS cost
    FROM fact_clicks
    WHERE spp_domain = ?
      AND event_date >= ?
      AND event_date <= ?
    GROUP BY gr_key
),

impressions_part AS
(
    SELECT
        toString(%[1]s) as gr_key,
        count(*) AS impressions,
        sumIf(win_final_price, format IN ('BAN', 'NAT')) / 1000 AS cost
    FROM fact_impressions
    WHERE spp_domain = ?
      AND event_date >= ?
      AND event_date <= ?
    GROUP BY gr_key
)

SELECT
    coalesce(c.gr_key, i.gr_key) AS group_key,
    ifNull(i.impressions, 0) AS impressions,
    ifNull(c.clicks, 0) AS clicks,
    ifNull(c.cost, 0) + ifNull(i.cost, 0) AS cost
FROM clicks_part AS c
FULL OUTER JOIN impressions_part AS i
    ON c.gr_key = i.gr_key
ORDER BY group_key
`, groupByExpr)

	rows, err := clickhouseConn.Query(
		ctx,
		query,
		sspDomain,
		dateStart,
		dateEnd,
		sspDomain,
		dateStart,
		dateEnd,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]getWmApiResponse, 0)
	for rows.Next() {
		var item getWmApiResponse
		if err := rows.Scan(&item.GroupKey, &item.Impressions, &item.Clicks, &item.Cost); err != nil {
			return nil, err
		}
		res = append(res, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}
