# Go-ing Gopher Hunting with Elasticsearch and Go

This repository provides an introductory example of using the Elasticsearch Go client to find documents in Elasticsearch. Specifically, it covers three types of search:

1. Traditional keyword search.
2. Vector search, making use of the [sentence-transformers/msmarco-MiniLM-L-12-v3 model from Hugging Face](https://huggingface.co/sentence-transformers/msmarco-MiniLM-L-12-v3) to generate the embeddings.
3. Hybrid search combining the keyword and vector approaches.

## How to Run

# Elasticsearch Instance Setup

The quickest way to setup your own cluster is to register for a [free trial of Elastic Cloud](https://www.elastic.co/cloud/elasticsearch-service/signup). You'll need to perform these additional steps:

1. Note your Cloud ID
2. [Generate an API Key](https://www.elastic.co/guide/en/kibana/current/api-keys.html)
3. Populate your instance with data in the same format as those in the [Sources](https://github.com/carlyrichmond/gopher-hunting-elasticsearch#sources) section below
4. [Upload your model from Hugging Face using Eland](https://www.elastic.co/guide/en/elasticsearch/client/eland/current/machine-learning.html#ml-nlp-pytorch)
5. Enriching your ingested documents using an [ingest pipeline](https://www.elastic.co/guide/en/elasticsearch/reference/current/ingest.html)

### Pre-requisites

This script requires setting the essential environment variables before running the script. I recommend using something like `direnv`, invoked via `.envrc` and then adding the variables to a top-level `.env` file. Alternatively you can explicitly set the environment variables in your current session according to your operating system.

The following environment variables are required:

- `ELASTIC_CLOUD_ID=<MY_INSTANCE_CLOUD_ID>` 
- `ELASTIC_API_KEY=<MY_API_KEY>`

### Starting the server

Running `server.go` will start a `net/http` server on port `80` that you can use to query Elasticsearch:

```bash
cd server
go run .
```

Navigate to the below URLs to obtain the Gopher search results for each search type:

* Keyword: [http://localhost/gophers](http://localhost/gophers)
* Vector: [http://localhost/vector-gophers](http://localhost/vector-gophers)
* Vector with filter: [http://localhost/filtered-gophers](http://localhost/filtered-gophers)
* Hybrid search with manual boosting: [http://localhost/hybrid-gophers](http://localhost/hybrid-gophers)
* Hybrid search with RRF: [http://localhost/rrf-gophers](http://localhost/rrf-gophers)

## Slides

The slides from the [Women Who Go meetup @ Elastic](https://www.meetup.com/women-who-go-london/events/295633460/) are available in the [docs/slides](./docs/slides/) folder.

## Sources

The below set of rodent-focused [Wikipedia](https://en.wikipedia.org/wiki/Main_Page) pages have been extracted to Elasticsearch using the [Elastic Web Crawler](https://www.elastic.co/web-crawler):

* [Rodent | Wikipedia](https://en.wikipedia.org/wiki/Rodent)
* [Gopher | Wikipedia](https://en.wikipedia.org/wiki/Gopher)
* [Rat | Wikipedia](https://en.wikipedia.org/wiki/Rat)
* [Prairie Dog | Wikipedia](https://en.wikipedia.org/wiki/Prairie_dog)
* [Porcupine | Wikipedia](https://en.wikipedia.org/wiki/Porcupine)
* [Guinea Pig | Wikipedia](https://en.wikipedia.org/wiki/Guinea_pig)
* [Hamster | Wikipedia](https://en.wikipedia.org/wiki/Hamster)
* [Capybara | Wikipedia](https://en.wikipedia.org/wiki/Capybara)
* [Pedetes | Wikipedia](https://en.wikipedia.org/wiki/Pedetes)
* [Beaver | Wikipedia](https://en.wikipedia.org/wiki/Beaver)
* [House Mouse | Wikipedia](https://en.wikipedia.org/wiki/House_mouse)
* [Squirrel | Wikipedia](https://en.wikipedia.org/wiki/Squirrel)

If you're new to Go and would like to build your own Web Crawler, I recommend having a stab at [this exercise in the Tour of Go](https://go.dev/tour/concurrency/10) where you can build your own concurrent web crawler.

## Resources

Check out the below resources to learn more about Elasticsearch, Keyword Search and Vector Search.

### Elasticsearch

1. [Elasticsearch](https://www.elastic.co/elasticsearch/)
2. [Elasticsearch Go Client](https://www.elastic.co/guide/en/elasticsearch/client/go-api/current/index.html)
3. [Understanding Analysis in Elasticsearch (Analyzers) by Bo Andersen | #CodingExplained](https://codingexplained.com/coding/elasticsearch/understanding-analysis-in-elasticsearch-analyzers)

### Vector Search

1. [code.sajari.com/word2vec](https://pkg.go.dev/code.sajari.com/word2vec)
2. [huggingface | pkg.go.dev](https://pkg.go.dev/github.com/nlpodyssey/spago/pkg/nlp/transformers/huggingface)
3. [What is Vector Search | Elastic](https://www.elastic.co/what-is/vector-search)

### LLMs and Natural Language Processing

1. [BERT 101: State Of The Art NLP Model Explained | Hugging Face](https://huggingface.co/blog/bert-101)
2. [sentence-transformers/msmarco-MiniLM-L-12-v3 | Hugging Face](https://huggingface.co/sentence-transformers/msmarco-MiniLM-L-12-v3)