package domain

type Comment struct {
	ID     int      `json:"id" db:"id"`
	PostID int      `json:"postID" db:"post_id"`
	Author UserCard `json:"author"`
	Text   string   `json:"text" db:"text"`
}

type CommentCredentials struct {
	UserID int    `json:"userID"`
	PostID int    `json:"postID"`
	Text   string `json:"text"`
}
