package config

type IndexConfiguration struct {
	ReaderThreads  int
	MaxConcurrency int
	WriterBuffer   int
}
