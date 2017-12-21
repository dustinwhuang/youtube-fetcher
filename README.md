# youtube-fetcher

A utility for fetching YouTube video metadata written in Go

## Usage

> go build -o youtube-fetcher

> export YOUTUBE_KEY=<YOUTUBE_API_KEY>

> ./youtube-fetcher <snippet_output_file> <contentDetails_output_file> <statistics_output_file> <topicDetails_output_file> <query_file>

> <query_file> is a single-line comma separated list of topics to be used for searching
