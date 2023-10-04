package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

// Vector search example with HuggingFace Generated Embeddings
func VectorSearchWithGeneratedQueryVector(term string) []Rodent {
	var vector []float32 = GetTextEmbeddingForQuery(term)

	if vector == nil {
		log.Fatal("Unable to generate vector")
		return nil
	}

	res, err := client.KnnSearch("vector-search-rodents").
		Knn(&types.CoreKnnQuery{
			Field:         "text_embedding.predicted_value",
			K:             10,
			NumCandidates: 10,
			QueryVector:   vector,
		}).
		Do(context.Background())

	if err != nil {
		log.Fatal(err)
		return nil
	}

	return GetRodents(res.Hits.Hits)
}

// HuggingFace text embedding helper
func GetTextEmbeddingForQuery(term string) []float32 {
	// HTTP endpoint
	model := "sentence-transformers/msmarco-minilm-l-12-v3"
	posturl := fmt.Sprintf("https://api-inference.huggingface.co/pipeline/feature-extraction/%s", model)

	// JSON body
	body := []byte(fmt.Sprintf(`{
		"inputs": "%s",
		"options": {"wait_for_model":True}
	}`, term))

	// Create a HTTP post request
	r, err := http.NewRequest("POST", posturl, bytes.NewBuffer(body))

	if err != nil {
		log.Fatal(err)
		return nil
	}

	token := os.Getenv("HUGGING_FACE_TOKEN")
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	var post []float32
	derr := json.NewDecoder(res.Body).Decode(&post)

	if derr != nil {
		log.Fatal(derr)
		return nil
	}

	return post
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
