package algoliautils

import (
	"fmt"
	"strconv"
)

type AlgoliaMatch struct {
	Value      string
	MatchLevel string
}

type AlgoliaResult struct {
	Object             map[string]interface{}
	HighlightedResults map[string]AlgoliaMatch
}

func (this *AlgoliaResult) ObjectID() string {
	o, ok := this.Object["objectID"]
	if ok == true {
		return String(o)
	}

	return ""
}

/**
  "hits": [
    {
      "firstname": "Jimmie",
      "lastname": "Barninger",
      "objectID": "433",
      "_highlightResult": {
        "firstname": {
          "value": "<em>Jimmie</em>",
          "matchLevel": "partial"
        },
        "lastname": {
          "value": "Barninger",
          "matchLevel": "none"
        },
        "company": {
          "value": "California <em>Paint</em> & Wlpaper Str",
          "matchLevel": "partial"
        }
      }
    }
  ],
  "page": 0,
  "nbHits": 1,
  "nbPages": 1,
  "hitsPerPage": 20,
  "processingTimeMS": 1,
  "query": "jimmie paint",
  "params": "query=jimmie+paint&attributesToRetrieve=firstname,lastname&hitsPerPage=50"

*/
type AlgoliaSearchResponse struct {
	Hits             []AlgoliaResult
	Page             int64
	NBHits           int64
	NBPages          int64
	ProcessingTimeMS int64
	Query            string
	Params           string
}

func NewAlgoliaSearchResponse(r interface{}) AlgoliaSearchResponse {
	out := AlgoliaSearchResponse{}

	if m, ok := r.(map[string]interface{}); ok == true {
		out.Page = Int(m["page"])
		out.NBHits = Int(m["nbHits"])
		out.NBPages = Int(m["nbPages"])
		out.ProcessingTimeMS = Int(m["processingTimeMS"])
		out.Query = String(m["query"])
		out.Params = String(m["params"])

		if results, has_results := m["hits"].([]interface{}); has_results == true {
			for _, result := range results {
				if result_map, has_result_map := result.(map[string]interface{}); has_result_map == true {
					out.Hits = append(out.Hits, NewAlgoliaResult(result_map))
				}
			}
		}
	}

	return out
}

func NewAlgoliaResult(results_map map[string]interface{}) AlgoliaResult {

	out := AlgoliaResult{}
	out.HighlightedResults = make(map[string]AlgoliaMatch)
	out.Object = make(map[string]interface{})

	for key, value := range results_map {
		if key == "_highlightResult" {
			if casted, ok := value.(map[string]interface{}); ok == true {
				for match_key, match := range casted {
					if match_map, has_match_map := match.(map[string]interface{}); has_match_map == true {
						out.HighlightedResults[match_key] = NewAlgoliaMatch(match_map)
					}
				}
			}
		} else {
			out.Object[key] = value
		}
	}

	return out
}

func NewAlgoliaMatch(match_map map[string]interface{}) AlgoliaMatch {

	out := AlgoliaMatch{}
	out.Value = String(match_map["value"])
	out.MatchLevel = String(match_map["matchLevel"])
	return out
}

func Int(r interface{}) int64 {
	switch casted := r.(type) {
	case int:
		return int64(casted)
	case int64:
		return casted
	case string:
		out, _ := strconv.ParseInt(casted, 10, 64)
		return out
	default:
		return 0
	}
}

func String(r interface{}) string {
	switch casted := r.(type) {
	case string:
		return casted
	case fmt.Stringer:
		return casted.String()
	default:
		return ""
	}
}
