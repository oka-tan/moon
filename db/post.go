//Package db wraps database access by providing entities
package db

import (
	"time"

	"github.com/uptrace/bun"
)

//Post is a post in the database to be indexed
type Post struct {
	bun.BaseModel `bun:"table:post,alias:post"`

	Board                 string     `bun:"board,pk"`
	PostNumber            int64      `bun:"post_number,pk"`
	ThreadNumber          int64      `bun:"thread_number"`
	Op                    bool       `bun:"op"`
	Deleted               bool       `bun:"deleted"`
	Hidden                bool       `bun:"hidden"`
	TimePosted            time.Time  `bun:"time_posted"`
	LastModified          time.Time  `bun:"last_modified"`
	CreatedAt             time.Time  `bun:"created_at"`
	Name                  *string    `bun:"name"`
	Tripcode              *string    `bun:"tripcode"`
	Capcode               *string    `bun:"capcode"`
	PosterID              *string    `bun:"poster_id"`
	Country               *string    `bun:"country"`
	Flag                  *string    `bun:"flag"`
	Email                 *string    `bun:"email"`
	Subject               *string    `bun:"subject"`
	Comment               *string    `bun:"comment"`
	HasMedia              bool       `bun:"has_media"`
	MediaDeleted          *bool      `bun:"media_deleted"`
	TimeMediaDeleted      *time.Time `bun:"time_media_deleted"`
	MediaTimestamp        *int64     `bun:"media_timestamp"`
	Media4chanHash        *[]byte    `bun:"media_4chan_hash"`
	MediaInternalHash     *[]byte    `bun:"media_internal_hash"`
	ThumbnailInternalHash *[]byte    `bun:"thumbnail_internal_hash"`
	MediaExtension        *string    `bun:"media_extension" json:"media_extension"`
	MediaFileName         *string    `bun:"media_file_name" json:"media_file_name"`
	MediaSize             *int       `bun:"media_size" json:"media_size"`
	MediaHeight           *int16     `bun:"media_height" json:"media_height"`
	MediaWidth            *int16     `bun:"media_width" json:"media_width"`
	ThumbnailHeight       *int16     `bun:"thumbnail_height" json:"thumbnail_height"`
	ThumbnailWidth        *int16     `bun:"thumbnail_width" json:"thumbnail_width"`
	Spoiler               *bool      `bun:"spoiler" json:"spoiler"`
	CustomSpoiler         *int16     `bun:"custom_spoiler" json:"custom_spoiler"`
	Sticky                *bool      `bun:"sticky" json:"sticky"`
	Closed                *bool      `bun:"closed" json:"closed"`
	Posters               *int16     `bun:"posters" json:"posters"`
	Replies               *int16     `bun:"replies" json:"replies"`
	Since4Pass            *int16     `bun:"since4pass" json:"since4pass"`
}
