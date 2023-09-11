# Go-ing Gopher Hunting with Elasticsearch and Go

This repository provides an introductory example of using the Elasticsearch Go client to find documents in Elasticsearch. Specifically, it covers three types of search:

1. Traditional keyword search
2. Vector search, making use of the <MODEL> from Hugging Face
3. Hybrid search combinging the keyword and vector approaches

# How to Run

## Pre-requisites

This script requires setting the essential environment variables before running the script. I recommend using something like `direnv`, invoked via `.envrc` and then adding the variables to a top-level `.env` file. 

The following environment variables are required:

- `ELASTIC_CLOUD_ID=<MY_INSTANCE_CLOUD_ID>` 
- `ELASTIC_API_KEY=<MY_API_KEY>`

## Starting the server

Running `server.go` will start a `net/http` server that you can use to query Elasticsearch:

```bash
cd server
go run .
```

# Slides

# Sources
