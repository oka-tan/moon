package lnx

import (
	"encoding/base64"
	"moon/db"
	"time"
)

var one int64 = 1
var zero int64 = 0

//Post is an intermediate struct that helps
//us unmarshal posts in the db to JSON we can
//send to Lnx
type Post struct {
	PostNumber     int64     `json:"post_number"`
	ThreadNumber   int64     `json:"thread_number"`
	Op             int64     `json:"op"`
	Deleted        int64     `json:"deleted"`
	TimePosted     time.Time `json:"time_posted"`
	Name           *string   `json:"name,omitempty"`
	Tripcode       *string   `json:"tripcode,omitempty"`
	Capcode        *string   `json:"capcode,omitempty"`
	PosterID       *string   `json:"poster_id,omitempty"`
	Country        *string   `json:"country,omitempty"`
	Flag           *string   `json:"flag,omitempty"`
	Email          *string   `json:"email,omitempty"`
	Subject        *string   `json:"subject,omitempty"`
	Comment        *string   `json:"comment,omitempty"`
	HasMedia       int64     `json:"has_media"`
	MediaDeleted   *int64    `json:"media_deleted,omitempty"`
	Media4chanHash *string   `json:"media_4chan_hash,omitempty"`
	MediaExtension *string   `json:"media_extension,omitempty"`
	MediaFileName  *string   `json:"media_file_name,omitempty"`
	Spoiler        *int64    `json:"spoiler,omitempty"`
	Sticky         *int64    `json:"sticky,omitempty"`
	Since4Pass     *int64    `json:"since4pass,omitempty"`
}

//DbPostToLnxPost turns a db.Post into a Post
func DbPostToLnxPost(p *db.Post) Post {
	var media4chanHash *string

	if p.Media4chanHash != nil {
		media4chanHashV := base64.StdEncoding.EncodeToString(*p.Media4chanHash)
		media4chanHash = &media4chanHashV
	}

	return Post{
		PostNumber:     p.PostNumber,
		ThreadNumber:   p.ThreadNumber,
		Op:             boolToInt64(p.Op),
		Deleted:        boolToInt64(p.Deleted),
		TimePosted:     p.TimePosted,
		Name:           p.Name,
		Tripcode:       p.Tripcode,
		Capcode:        p.Capcode,
		PosterID:       p.PosterID,
		Country:        p.Country,
		Flag:           p.Flag,
		Email:          p.Email,
		Subject:        p.Subject,
		Comment:        p.Comment,
		HasMedia:       boolToInt64(p.HasMedia),
		MediaDeleted:   boolPointerToInt64Pointer(p.MediaDeleted),
		Media4chanHash: media4chanHash,
		MediaExtension: p.MediaExtension,
		MediaFileName:  p.MediaFileName,
		Spoiler:        boolPointerToInt64Pointer(p.Spoiler),
		Sticky:         boolPointerToInt64Pointer(p.Sticky),
		Since4Pass:     int16PointerToInt64Pointer(p.Since4Pass),
	}
}

//DbPostsToLnxPosts converts an array of db.Post into
//an array of Post
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
	}

	return 0
}

func boolPointerToInt64Pointer(b *bool) *int64 {
	if b != nil && *b {
		return &one
	}

	return &zero
}

func int16PointerToInt64Pointer(i *int16) *int64 {
	if i == nil {
		return nil
	}

	v := int64(*i)
	return &v
}
