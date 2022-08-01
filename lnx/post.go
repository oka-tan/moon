package lnx

import (
	"encoding/base64"
	"moon/db"
	"time"
)

type Post struct {
	PostNumber     int64     `json:"post_number"`
	ThreadNumber   int64     `json:"thread_number"`
	Op             int64     `json:"op"`
	Deleted        int64     `json:"deleted"`
	TimePosted     time.Time `json:"time_posted"`
	Name           string    `json:"name"`
	Tripcode       string    `json:"tripcode"`
	Capcode        string    `json:"capcode"`
	PosterID       string    `json:"poster_id"`
	Country        string    `json:"country"`
	Flag           string    `json:"flag"`
	Email          string    `json:"email"`
	Subject        string    `json:"subject"`
	Comment        string    `json:"comment"`
	HasMedia       int64     `json:"has_media"`
	MediaDeleted   int64     `json:"media_deleted"`
	Media4chanHash string    `json:"media_4chan_hash"`
	MediaExtension string    `json:"media_extension"`
	MediaFileName  string    `json:"media_file_name"`
	Spoiler        int64     `json:"spoiler"`
	Sticky         int64     `json:"sticky"`
	Since4Pass     int64     `json:"since4pass"`
}

func DbPostToLnxPost(p *db.Post) Post {
	var media4chanHash string

	if p.Media4chanHash != nil {
		media4chanHash = base64.StdEncoding.EncodeToString(*p.Media4chanHash)
	}

	return Post{
		PostNumber:     p.PostNumber,
		ThreadNumber:   p.ThreadNumber,
		Op:             boolToInt64(p.Op),
		Deleted:        boolToInt64(p.Deleted),
		TimePosted:     p.TimePosted,
		Name:           derefString(p.Name),
		Tripcode:       derefString(p.Tripcode),
		Capcode:        derefString(p.Capcode),
		PosterID:       derefString(p.PosterID),
		Country:        derefString(p.Country),
		Flag:           derefString(p.Flag),
		Email:          derefString(p.Email),
		Subject:        derefString(p.Subject),
		Comment:        derefString(p.Comment),
		HasMedia:       boolToInt64(p.HasMedia),
		MediaDeleted:   boolPointerToInt64(p.MediaDeleted),
		Media4chanHash: media4chanHash,
		MediaExtension: derefString(p.MediaExtension),
		MediaFileName:  derefString(p.MediaFileName),
		Spoiler:        boolPointerToInt64(p.Spoiler),
		Sticky:         boolPointerToInt64(p.Sticky),
		Since4Pass:     derefInt16(p.Since4Pass),
	}
}

func DbPostsToLnxPosts(posts []db.Post) []Post {
	result := make([]Post, 0, len(posts))

	for i := range posts {
		result = append(result, DbPostToLnxPost(&posts[i]))
	}

	return result
}

func boolToInt64(b bool) int64 {
	if b {
		return 1
	} else {
		return 0
	}
}

func boolPointerToInt64(b *bool) int64 {
	if b != nil && *b {
		return 1
	} else {
		return 0
	}
}

func derefInt16(i *int16) int64 {
	if i != nil {
		return int64(*i)
	} else {
		return 0
	}
}

func derefInt64(i *int64) int64 {
	if i != nil {
		return *i
	} else {
		return 0
	}
}

func derefString(s *string) string {
	if s != nil {
		return *s
	} else {
		return ""
	}
}
