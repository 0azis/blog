package store

import (
	"blog/internal/core/domain"

	"github.com/jmoiron/sqlx"
)

type tag struct {
	db *sqlx.DB
}

func (t tag) Create(tag domain.Tag) (int64, error) {
	var lastID int64

	_, err := t.db.Query(`delete from tags where post_id = ?`, tag.PostID)
	if err != nil {
		return lastID, err
	}

	for _, s := range tag.Tags {
		sqlResult, err := t.db.Exec(`insert into tags values (?, ?)`, tag.PostID, s)
		if err != nil {
			return lastID, err
		}
		lastID, err = sqlResult.RowsAffected()
		if err != nil || lastID == 0 {
			return lastID, err
		}
	}

	return lastID, nil
}

func (t tag) GetByPostID(postID int) (domain.Tag, error) {
	var tags domain.Tag
	rows, err := t.db.Query(`select * from tags where post_id = ?`, postID)
	if err != nil {
		return tags, err
	}

	for rows.Next() {
		var tag string
		var postID int
		err = rows.Scan(&postID, &tag)
		if err != nil {
			return tags, err
		}

		tags.PostID = postID
		tags.Tags = append(tags.Tags, tag)
	}

	return tags, nil
}

func (t tag) GetByPopularity() (domain.Tags, error) {
	var tags domain.Tags
	rows, err := t.db.Query(`select tag from tags group by tag order by count(*) desc limit 3`)
	if err != nil {
		return tags, err
	}

	for rows.Next() {
		var tag string
		err = rows.Scan(&tag)
		if err != nil {
			return tags, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}
