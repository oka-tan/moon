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

type Service struct {
	host   string
	client http.Client
}

func NewService(host string, port int) Service {
	return Service{
		host: fmt.Sprintf("%s:%d/indexes", host, port),
		client: http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *Service) Upsert(posts []db.Post, board string, previousScrape time.Time) error {
	deletables := make([]db.Post, 0, 10)

	for _, p := range posts {
		if p.CreatedAt.Before(previousScrape) {
			deletables = append(deletables, p)
		}
	}

	if len(deletables) > 0 {
		deleteRequest := buildDeleteRequest(deletables)

		for i := 1; ; i++ {
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

	for i := 1; ; i++ {
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

func (s *Service) Rollback(board string) error {
	for i := 1; ; i++ {
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

func (s *Service) Commit(board string) error {
	for i := 1; ; i++ {
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

func (s *Service) CreateIndex(board string, indexConfiguration config.IndexConfiguration, forceRecreate bool) {
	createIndexRequest := CreateIndexRequest{
		OverrideIfExists: forceRecreate,
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
					Type:    "date",
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

	if resp.StatusCode == 200 {
		return
	}

	if resp.StatusCode == 400 {
		log.Printf("Received status 400 creating index for %s\n", board)
		return
	}

	log.Fatalf("Received status %s creating index\n", resp.Status)
	return
}
