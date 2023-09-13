package search

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

type Rodent struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Url   string `json:"url"`
}

// Elasticsearch init
var cloudID = os.Getenv("ELASTIC_CLOUD_ID")
var apiKey = os.Getenv("ELASTIC_API_KEY")

var client, err = elasticsearch.NewTypedClient(elasticsearch.Config{
	CloudID: cloudID,
	APIKey:  apiKey,
})

// Traditional keyword search example
func KeywordSearch(term string) []Rodent {
	res, err := client.Search().
		Index("search-rodents").
		Query(&types.Query{
			Match: map[string]types.MatchQuery{
				"title": {Query: term},
			},
		}).
		From(0).
		Size(10).
		Do(context.Background())

	if err != nil {
		log.Fatal(err)
		return nil
	}

	return GetRodents(res.Hits.Hits)
}

// Vector search example
func VectorSearch(term string) []Rodent {
	res, err := client.Search().
		Index("search-rodents").
		Knn(types.KnnQuery{
			Field:         "ml.inference.predicted_value",
			K:             10,
			NumCandidates: 10,
			QueryVectorBuilder: &types.QueryVectorBuilder{
				TextEmbedding: &types.TextEmbedding{
					ModelId:   "sentence-transformers__msmarco-minilm-l-12-v3",
					ModelText: term,
				},
			}}).Do(context.Background())

	if err != nil {
		log.Fatal(err)
		return nil
	}

	return GetRodents(res.Hits.Hits)
}

func GetRodents(hits []types.Hit) []Rodent {
	var rodents []Rodent

	for _, hit := range hits {
		var currentRodent Rodent
		err := json.Unmarshal(hit.Source_, &currentRodent)

		if err != nil {
			log.Fatal(err)
			return nil
		}

		currentRodent.ID = hit.Id_
		rodents = append(rodents, currentRodent)
	}

	return rodents
}
