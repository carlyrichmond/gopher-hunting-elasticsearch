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

var client = GetElasticsearchClient()

// Elasticsearch init
func GetElasticsearchClient() *elasticsearch.TypedClient {
	var cloudID = os.Getenv("ELASTIC_CLOUD_ID")
	var apiKey = os.Getenv("ELASTIC_API_KEY")

	var es, err = elasticsearch.NewTypedClient(elasticsearch.Config{
		CloudID: cloudID,
		APIKey:  apiKey,
	})

	if err != nil {
		log.Fatalf("Unable to connect: %s", err)
		os.Exit(3)
	}

	return es
}

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
		Index("vector-search-rodents").
		Knn(types.KnnQuery{
			Field:         "text_embedding.predicted_value",
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

// Vector search example with filter
func VectorSearchWithFilter(term string) []Rodent {
	res, err := client.Search().
		Index("vector-search-rodents").
		Knn(types.KnnQuery{
			Field: "text_embedding.predicted_value",
			K:     10,
			Filter: []types.Query{
				{
					Match: map[string]types.MatchQuery{
						"body_content": {Query: "rodent"},
					},
				},
			},
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

// Hybrid search example
func HybridSearchWithBoost(term string) []Rodent {
	var knnBoost float32 = 0.2

	res, err := client.Search().
		Index("vector-search-rodents").
		Knn(types.KnnQuery{
			Field:         "text_embedding.predicted_value",
			Boost:         &knnBoost,
			K:             10,
			NumCandidates: 10,
			QueryVectorBuilder: &types.QueryVectorBuilder{
				TextEmbedding: &types.TextEmbedding{
					ModelId:   "sentence-transformers__msmarco-minilm-l-12-v3",
					ModelText: term,
				},
			}}).
		Query(&types.Query{
			Match: map[string]types.MatchQuery{
				"title": {Query: term},
			},
		}).
		From(0).
		Size(2).
		Do(context.Background())

	if err != nil {
		log.Fatal(err)
		return nil
	}

	return GetRodents(res.Hits.Hits)
}

// Hybrid search with RRF
func HybridSearchWithRRF(term string) []Rodent {
	var windowSize int64 = 10 // min required
	var rankConstant int64 = 42

	res, err := client.Search().
		Index("vector-search-rodents").
		Knn(types.KnnQuery{
			Field:         "text_embedding.predicted_value",
			K:             10,
			NumCandidates: 10,
			QueryVectorBuilder: &types.QueryVectorBuilder{
				TextEmbedding: &types.TextEmbedding{
					ModelId:   "sentence-transformers__msmarco-minilm-l-12-v3",
					ModelText: term,
				},
			}}).
		Query(&types.Query{
			Match: map[string]types.MatchQuery{
				"title": {Query: term},
			},
		}).
		Rank(&types.RankContainer{
			Rrf: &types.RrfRank{
				WindowSize:   &windowSize,
				RankConstant: &rankConstant,
			},
		}).
		Do(context.Background())

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
