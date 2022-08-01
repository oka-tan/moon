package lnx

import (
	"fmt"
	"moon/config"
)

type CreateIndexRequest struct {
	OverrideIfExists bool                    `json:"override_if_exists"`
	Index            CreateIndexRequestIndex `json:"index"`
}

type CreateIndexRequestIndex struct {
	Name           string                `json:"name"`
	StorageType    string                `json:"storage_type"`
	StripStopWords bool                  `json:"strip_stop_words"`
	Fields         map[string]IndexField `json:"fields"`
	SearchFields   []string              `json:"search_fields"`
	ReaderThreads  int                   `json:"reader_threads"`
	MaxConcurrency int                   `json:"max_concurrency"`
	WriterBuffer   int                   `json:"writer_buffer"`
	WriterThreads  int                   `json:"writer_threads"`
}

type IndexField struct {
	Type    string `json:"type"`
	Stored  bool   `json:"stored"`
	Indexed bool   `json:"indexed"`
	Fast    bool   `json:"fast"`
}

func IndexRequestFromConfiguration(board string, indexConfiguration config.IndexConfiguration) CreateIndexRequest {
	return CreateIndexRequest{
		OverrideIfExists: false,
		Index: CreateIndexRequestIndex{
			Name:           fmt.Sprintf("post_%s", board),
			StorageType:    "filesystem",
			StripStopWords: false,
			Fields: map[string]IndexField{
				"post_number": {
					Type:    "i64",
					Stored:  true,
					Indexed: true,
					Fast:    true,
				},
				"thread_number": {
					Type:    "i64",
					Stored:  false,
					Indexed: true,
					Fast:    false,
				},
				"op": {
					Type:    "i64",
					Stored:  false,
					Indexed: true,
					Fast:    false,
				},
				"deleted": {
					Type:    "i64",
					Stored:  false,
					Indexed: true,
					Fast:    false,
				},
				"time_posted": {
					Type:    "i64",
					Stored:  false,
					Indexed: true,
					Fast:    false,
				},
				"name": {
					Type:    "text",
					Stored:  false,
					Indexed: true,
				},
				"tripcode": {
					Type:    "string",
					Stored:  false,
					Indexed: true,
				},
				"capcode": {
					Type:    "string",
					Stored:  false,
					Indexed: true,
				},
				"poster_id": {
					Type:    "string",
					Stored:  false,
					Indexed: true,
				},
				"country": {
					Type:    "string",
					Stored:  false,
					Indexed: true,
				},
				"flag": {
					Type:    "string",
					Stored:  false,
					Indexed: true,
				},
				"email": {
					Type:    "string",
					Stored:  false,
					Indexed: true,
				},
				"subject": {
					Type:    "text",
					Stored:  false,
					Indexed: true,
				},
				"comment": {
					Type:    "text",
					Stored:  false,
					Indexed: true,
				},
				"has_media": {
					Type:    "i64",
					Stored:  false,
					Indexed: true,
					Fast:    false,
				},
				"media_deleted": {
					Type:    "i64",
					Stored:  false,
					Indexed: true,
					Fast:    false,
				},
				"media_4chan_hash": {
					Type:    "string",
					Stored:  false,
					Indexed: true,
				},
				"media_extension": {
					Type:    "string",
					Stored:  false,
					Indexed: true,
				},
				"media_file_name": {
					Type:    "text",
					Stored:  false,
					Indexed: true,
				},
				"spoiler": {
					Type:    "i64",
					Stored:  false,
					Indexed: true,
					Fast:    false,
				},
				"sticky": {
					Type:    "i64",
					Stored:  false,
					Indexed: true,
					Fast:    false,
				},
				"since4pass": {
					Type:    "i64",
					Stored:  false,
					Indexed: true,
					Fast:    false,
				},
			},
			SearchFields:   []string{"comment", "subject", "name", "media_file_name"},
			ReaderThreads:  indexConfiguration.ReaderThreads,
			MaxConcurrency: indexConfiguration.MaxConcurrency,
			WriterBuffer:   indexConfiguration.WriterBuffer,
			WriterThreads:  1,
		},
	}
}
