package lnx

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"moon/config"
	"moon/db"
	"net/http"
	"time"
)

//Service wraps writes and upserts to Lnx
type Service struct {
	host           string
	client         http.Client
	readerThreads  int
	maxConcurrency int
	writerBuffer   int
}

//NewService constructs and returns a Service
func NewService(conf config.LnxConfig) Service {
	return Service{
		host: fmt.Sprintf("%s:%d/indexes", conf.Host, conf.Port),
		client: http.Client{
			Timeout: 30 * time.Second,
		},
		readerThreads:  conf.ReaderThreads,
		maxConcurrency: conf.MaxConcurrency,
		writerBuffer:   conf.WriterBuffer,
	}
}

//Upsert upserts an array of posts into Lnx
func (s *Service) Upsert(posts []db.Post, board string, previousScrape time.Time) error {
	deletables := make([]db.Post, 0, 10)

	for _, p := range posts {
		if p.CreatedAt.Before(previousScrape) {
			deletables = append(deletables, p)
		}
	}

	if len(deletables) > 0 {
		deleteRequest := buildDeleteRequest(deletables)

		for i := 0; ; i++ {
			pipeReader, pipeWriter := io.Pipe()

			go func() {
				jsonEncoder := json.NewEncoder(pipeWriter)
				err := jsonEncoder.Encode(&deleteRequest)
				pipeWriter.CloseWithError(err)
			}()

			r, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/post_%s/documents/query", s.host, board), pipeReader)
			resp, err := s.client.Do(r)

			if err != nil {
				if i < 3 {
					log.Printf("Error performing deletion request: %s", err)
					time.Sleep(30 * time.Second)
					continue
				} else {
					return fmt.Errorf("Error performing deletion request: %s", err)
				}
			}

			resp.Body.Close()

			if resp.StatusCode != 200 {
				return fmt.Errorf("Error deleting old posts: request received status %s", resp.Status)
			}

			break
		}
	}

	lnxPosts := DbPostsToLnxPosts(posts)

	for i := 0; ; i++ {
		pipeReader, pipeWriter := io.Pipe()

		go func() {
			jsonEncoder := json.NewEncoder(pipeWriter)
			err := jsonEncoder.Encode(&lnxPosts)
			pipeWriter.CloseWithError(err)
		}()

		resp, err := s.client.Post(fmt.Sprintf("%s/post_%s/documents", s.host, board), "application/json", pipeReader)

		if err != nil {
			if i < 3 {
				log.Printf("Error performing insertion request: %s", err)
				time.Sleep(30 * time.Second)

				continue
			} else {
				return fmt.Errorf("Error performing insertion request: %s", err)
			}
		}

		resp.Body.Close()

		if resp.StatusCode != 200 {
			return fmt.Errorf("Error inserting posts: request received status %s", resp.Status)
		}

		break
	}

	return nil
}

//Rollback rolls back index modifications
func (s *Service) Rollback(board string) error {
	for i := 0; ; i++ {
		resp, err := s.client.Post(fmt.Sprintf("%s/post_%s/rollback", s.host, board), "", nil)

		if err != nil {
			if i < 3 {
				log.Printf("Error performing rollback: %s", err)
				time.Sleep(30 * time.Second)

				continue
			} else {
				return fmt.Errorf("Error performing rollback: %s", err)
			}
		}

		resp.Body.Close()

		if resp.StatusCode != 200 {
			return fmt.Errorf("Rollback request received status %s", resp.Status)
		}

		return nil
	}
}

//Commit commits index modifications
func (s *Service) Commit(board string) error {
	for i := 0; ; i++ {
		resp, err := s.client.Post(fmt.Sprintf("%s/post_%s/commit", s.host, board), "", nil)

		if err != nil {
			if i < 3 {
				log.Printf("Error performing commit: %s", err)
				time.Sleep(30 * time.Second)

				continue
			} else {
				return fmt.Errorf("Error performing commit: %s", err)
			}
		}

		resp.Body.Close()

		if resp.StatusCode != 200 {
			return fmt.Errorf("Request received status %s", resp.Status)
		}

		return nil
	}
}

//CreateIndex creates the index described by the configuration passed
func (s *Service) CreateIndex(conf config.BoardConfig) {
	createIndexRequest := CreateIndexRequest{
		OverrideIfExists: conf.ForceRecreate,
		Index: CreateIndexRequestIndex{
			Name:                    fmt.Sprintf("post_%s", conf.Name),
			StorageType:             "filesystem",
			StripStopWords:          false,
			SetConjunctionByDefault: true,
			Fields: map[string]IndexField{
				"post_number": {
					Type:     "i64",
					Stored:   true,
					Indexed:  true,
					Fast:     true,
					Required: true,
				},
				"thread_number": {
					Type:     "i64",
					Stored:   false,
					Indexed:  true,
					Fast:     false,
					Required: true,
				},
				"op": {
					Type:     "i64",
					Stored:   false,
					Indexed:  true,
					Fast:     false,
					Required: true,
				},
				"deleted": {
					Type:     "i64",
					Stored:   false,
					Indexed:  true,
					Fast:     false,
					Required: true,
				},
				"time_posted": {
					Type:     "date",
					Stored:   false,
					Indexed:  true,
					Fast:     false,
					Required: true,
				},
				"name": {
					Type:     "text",
					Stored:   false,
					Indexed:  true,
					Required: false,
				},
				"tripcode": {
					Type:     "string",
					Stored:   false,
					Indexed:  true,
					Required: false,
				},
				"capcode": {
					Type:     "string",
					Stored:   false,
					Indexed:  true,
					Required: false,
				},
				"poster_id": {
					Type:     "string",
					Stored:   false,
					Indexed:  true,
					Required: false,
				},
				"country": {
					Type:     "string",
					Stored:   false,
					Indexed:  true,
					Required: false,
				},
				"flag": {
					Type:     "string",
					Stored:   false,
					Indexed:  true,
					Required: false,
				},
				"email": {
					Type:     "string",
					Stored:   false,
					Indexed:  true,
					Required: false,
				},
				"subject": {
					Type:     "text",
					Stored:   false,
					Indexed:  true,
					Required: false,
				},
				"comment": {
					Type:     "text",
					Stored:   false,
					Indexed:  true,
					Required: false,
				},
				"has_media": {
					Type:     "i64",
					Stored:   false,
					Indexed:  true,
					Fast:     false,
					Required: true,
				},
				"media_deleted": {
					Type:     "i64",
					Stored:   false,
					Indexed:  true,
					Fast:     false,
					Required: false,
				},
				"media_4chan_hash": {
					Type:     "string",
					Stored:   false,
					Indexed:  true,
					Required: false,
				},
				"media_extension": {
					Type:     "string",
					Stored:   false,
					Indexed:  true,
					Required: false,
				},
				"media_file_name": {
					Type:     "text",
					Stored:   false,
					Indexed:  true,
					Required: false,
				},
				"spoiler": {
					Type:     "i64",
					Stored:   false,
					Indexed:  true,
					Fast:     false,
					Required: false,
				},
				"sticky": {
					Type:     "i64",
					Stored:   false,
					Indexed:  true,
					Fast:     false,
					Required: false,
				},
				"since4pass": {
					Type:     "i64",
					Stored:   false,
					Indexed:  true,
					Fast:     false,
					Required: false,
				},
			},
			SearchFields:   []string{"comment", "subject", "name", "media_file_name"},
			ReaderThreads:  s.readerThreads,
			MaxConcurrency: s.maxConcurrency,
			WriterBuffer:   s.writerBuffer,
			WriterThreads:  1,
		},
	}

	reader, writer := io.Pipe()

	go func() {
		jsonEncoder := json.NewEncoder(writer)
		err := jsonEncoder.Encode(createIndexRequest)
		writer.CloseWithError(err)
	}()

	resp, err := http.Post(s.host, "application/json", reader)

	if err != nil {
		log.Fatalf("Error creating index: %v", err)
	}

	resp.Body.Close()

	if resp.StatusCode == 400 {
		log.Printf("Received status 400 creating index for %s\n", conf.Name)
		return
	}

	if resp.StatusCode != 200 {
		log.Fatalf("Received status %s creating index\n", resp.Status)
	}
}
