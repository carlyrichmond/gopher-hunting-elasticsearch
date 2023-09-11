module server

go 1.20

replace search => ../search

require search v0.0.0-00010101000000-000000000000

require (
	github.com/elastic/elastic-transport-go/v8 v8.0.0-20230329154755-1a3c63de0db6 // indirect
	github.com/elastic/go-elasticsearch/v8 v8.9.0 // indirect
)
