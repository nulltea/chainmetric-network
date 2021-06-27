package shared

import (
	"fmt"

	"github.com/timoth-y/chainmetric-core/utils"
)

// BuildQuery builds CouchDB query by given parameters:
//
// `selector`: a filter string declaring which documents to return
//
// `fields`: specifying which fields to be returned
//
// `sort`: expression containing how to sort selected records.
func BuildQuery(selector map[string]interface{}, sort ...string) string {
	query := map[string]interface{}{
		"selector": selector,
	}

	if len(sort) > 0 {
		query["sort"] = bindSortingPairs("asc", sort...)
	}

	fmt.Println(utils.MustEncode(query))

	return utils.MustEncode(query)
}

func bindSortingPairs(defaultValue string, kvs ...string) (res []map[string]string) {
	res = make([]map[string]string, 0)
	if len(kvs) % 2 != 0 {
		kvs = append(kvs, defaultValue)
	}

	for i := 1; i < len(kvs); i += 2 {
		sortMode := kvs[i]

		if sortMode != "asc" && sortMode != "desc" {
			sortMode = defaultValue
		}

		res = append(res, map[string]string{
			kvs[i - 1]: sortMode,
		})
	}

	return res
}
