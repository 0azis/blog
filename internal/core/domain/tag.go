package domain

type Tag struct {
	PostID int      `json:"postID" db:"post_id"`
	Tags   []string `json:"tags" db:"tag"`
}

// type TagCredentials struct {
// 	PostID int
// 	Tags   []string
// }

type Tags []string
