package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/elastic/elastic-transport-go/v8/elastictransport"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

type Rodent struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Url   string `json:"url"`
}

// Elasticsearch init
func GetElasticsearchClient() (*elasticsearch.TypedClient, error) {
	var cloudID = os.Getenv("ELASTIC_CLOUD_ID")
	var apiKey = os.Getenv("ELASTIC_API_KEY")

	var es, err = elasticsearch.NewTypedClient(elasticsearch.Config{
		CloudID: cloudID,
		APIKey:  apiKey,
		Logger:  &elastictransport.ColorLogger{os.Stdout, true, true},
	})

	if err != nil {
		return nil, fmt.Errorf("unable to connect: %w", err)
	}

	return es, nil
}

// Traditional keyword search example
func KeywordSearch(client *elasticsearch.TypedClient, term string) ([]Rodent, error) {
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
		return nil, fmt.Errorf("could not search for rodents: %w", err)
	}

	return getRodents(res.Hits.Hits)
}

// Vector search example
func VectorSearch(client *elasticsearch.TypedClient, term string) ([]Rodent, error) {
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
		return nil, fmt.Errorf("error in rodents vector search: %w", err)
	}

	return getRodents(res.Hits.Hits)
}

// Vector search example with Hugging Face Generated Embeddings
func VectorSearchWithGeneratedQueryVector(client *elasticsearch.TypedClient, term string) ([]Rodent, error) {
	vector, err := GetTextEmbeddingForQuery(term)
	if err != nil {
		return nil, err
	}

	if vector == nil {
		return nil, fmt.Errorf("unable to generate vector: %w", err)
	}

	res, err := client.Search().
		Index("vector-search-rodents").
		Knn(types.KnnQuery{
			Field:         "text_embedding.predicted_value",
			K:             10,
			NumCandidates: 10,
			QueryVector:   vector,
		}).
		Do(context.Background())

	if err != nil {
		return nil, err
	}

	return getRodents(res.Hits.Hits)
}

// HuggingFace text embedding helper
func GetTextEmbeddingForQuery(term string) ([]float32, error) {
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
		return nil, err
	}

	token := os.Getenv("HUGGING_FACE_TOKEN")
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var post []float32
	derr := json.NewDecoder(res.Body).Decode(&post)

	if derr != nil {
		return nil, err
	}

	return post, nil
}

// Vector search example with filter
func VectorSearchWithFilter(client *elasticsearch.TypedClient, term string) ([]Rodent, error) {
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
		return nil, err
	}

	return getRodents(res.Hits.Hits)
}

// Hybrid search example
func HybridSearchWithBoost(client *elasticsearch.TypedClient, term string) ([]Rodent, error) {
	var knnBoost float32 = 0.2
	var queryBoost float32 = 0.8

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
				"title": {
					Query: term,
					Boost: &queryBoost,
				},
			},
		}).
		Do(context.Background())

	if err != nil {
		return nil, err
	}

	return getRodents(res.Hits.Hits)
}

// Hybrid search with RRF
func HybridSearchWithRRF(client *elasticsearch.TypedClient, term string) ([]Rodent, error) {
	// Minimum required window size for the default result size of 10
	var windowSize int64 = 10
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
		return nil, err
	}

	return getRodents(res.Hits.Hits)
}

func getRodents(hits []types.Hit) ([]Rodent, error) {
	var rodents []Rodent

	for _, hit := range hits {
		var currentRodent Rodent
		err := json.Unmarshal(hit.Source_, &currentRodent)

		if err != nil {
			return nil, fmt.Errorf("an error occurred while unmarshaling rodent %s: %w", hit.Id_, err)
		}

		currentRodent.ID = hit.Id_
		rodents = append(rodents, currentRodent)
	}

	return rodents, nil
}
