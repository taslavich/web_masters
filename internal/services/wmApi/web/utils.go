package wmApiWeb

import "fmt"

func getWmApiGroupByExpression(groupBy string) (string, error) {
	switch groupBy {
	case "geo":
		return "geo", nil
	case "date":
		return "event_date", nil
	case "site":
		return "site_id", nil
	default:
		return "", fmt.Errorf("Invalid group_by value")
	}
}

func mergeSspFeeds(sspFeeds ...map[string]string) map[string]string {
	res := make(map[string]string)

	for _, feeds := range sspFeeds {
		for feed, sspDomain := range feeds {
			res[feed] = sspDomain
		}
	}

	return res
}
