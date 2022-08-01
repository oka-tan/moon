package lnx

import (
	"fmt"
	"moon/db"
	"strings"
)

func buildDeleteRequest(posts []db.Post) deleteRequest {
	var b strings.Builder

	fmt.Fprintf(&b, "post_number:%d", posts[0].PostNumber)

	if len(posts) > 1 {
		for _, p := range posts[1:] {
			fmt.Fprintf(&b, " OR post_number:%d", p.PostNumber)
		}
	}

	return deleteRequest{
		Query: query{normalQuery{Ctx: b.String()}},
		Limit: len(posts),
	}
}
