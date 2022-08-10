package lnx

type CreateIndexRequest struct {
	OverrideIfExists bool                    `json:"override_if_exists"`
	Index            CreateIndexRequestIndex `json:"index"`
}

type CreateIndexRequestIndex struct {
	Name                    string                `json:"name"`
	StorageType             string                `json:"storage_type"`
	SetConjunctionByDefault bool                  `json:"set_conjunction_by_default"`
	StripStopWords          bool                  `json:"strip_stop_words"`
	Fields                  map[string]IndexField `json:"fields"`
	SearchFields            []string              `json:"search_fields"`
	ReaderThreads           int                   `json:"reader_threads"`
	MaxConcurrency          int                   `json:"max_concurrency"`
	WriterBuffer            int                   `json:"writer_buffer"`
	WriterThreads           int                   `json:"writer_threads"`
}

type IndexField struct {
	Type    string `json:"type"`
	Stored  bool   `json:"stored"`
	Indexed bool   `json:"indexed"`
	Fast    bool   `json:"fast"`
}
