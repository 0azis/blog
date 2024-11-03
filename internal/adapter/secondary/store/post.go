package store

import (
	"blog/internal/core/domain"
	"blog/internal/core/utils"

	"github.com/jmoiron/sqlx"
)

var (
	tagPopular   = "select posts.id, posts.title, posts.preview, posts.date, posts.content, users.id, users.username, users.name, users.avatar, count(distinct views.user_id), count(distinct comments.id) from posts inner join users on posts.user_id= users.id inner join tags on posts.id = tags.post_id left join comments on posts.id = comments.post_id left join views on posts.id = views.post_id where tags.tag = ? and posts.public = 1 group by posts.id order by count(distinct views.user_id) desc limit ? offset ?"
	tagNewest    = "select posts.id, posts.title, posts.preview, posts.date, posts.content, users.id, users.username, users.name, users.avatar, count(distinct views.user_id), count(distinct comments.id) from posts inner join users on posts.user_id = users.id inner join tags on posts.id = tags.post_id left join comments on posts.id = comments.post_id left join views on posts.id = views.post_id where tags.tag = ? and posts.public = 1 group by posts.id order by date asc limit ? offset ?"
	noTagPopular = "select posts.id, posts.title, posts.preview, posts.date, posts.content, users.id, users.username, users.name, users.avatar, count(distinct views.user_id), count(distinct comments.id) from posts inner join users on posts.user_id = users.id left join views on posts.id = views.post_id left join comments on posts.id = comments.post_id where posts.public = 1 group by posts.id order by count(distinct views.user_id) desc limit ? offset ?"
	noTagNewest  = "select posts.id, posts.title, posts.preview, posts.date, posts.content, users.id, users.username, users.name, users.avatar, count(distinct views.user_id), count(distinct comments.id) from posts inner join users on posts.user_id = users.id left join views on posts.id = views.post_id left join comments on posts.id = comments.post_id where posts.public = 1 group by posts.id order by posts.date asc limit ? offset ?"
)

type post struct {
	db *sqlx.DB
}

func (p post) Create(post domain.PostCredentials) (int, error) {
	sqlResult, err := p.db.Exec(`insert into posts (user_id) values (?)`, post.UserID)
	if err != nil {
		return 0, err
	}
	postID, err := sqlResult.LastInsertId()
	return int(postID), err
}

func (p post) Update(postID int, post domain.PostCredentials) (int, error) {
	sqlResult, err := p.db.Exec(`update posts set title = ?, preview = ?, content = ? where user_id = ? and id = ?`, post.Title, post.Preview, post.Content, post.UserID, postID)
	if err != nil {
		return 0, err
	}
	lastID, err := sqlResult.RowsAffected()
	return int(lastID), err
}

func (p post) GetPostsByUser(userID int) ([]*domain.Post, error) {
	var posts []*domain.Post
	rows, err := p.db.Query(`select posts.id, posts.title, posts.preview, posts.date, posts.content, users.id, users.username, users.name, users.avatar, count(distinct views.user_id), count(distinct comments.post_id) from posts inner join users on posts.user_id = users.id left join views on posts.id = views.post_id left join comments on posts.id = comments.post_id where public = 1 and posts.user_id = ? group by posts.id`, userID)
	if err != nil {
		return posts, err
	}

	for rows.Next() {
		var post domain.Post
		err = rows.Scan(&post.ID, &post.Title, &post.Preview, &post.Date, &post.Content, &post.Author.ID, &post.Author.Username, &post.Author.Name, &post.Author.Avatar, &post.Counter.Views, &post.Counter.Comments)
		if err != nil {
			return posts, err
		}
		posts = append(posts, &post)
	}

	return posts, nil
}

func (p post) GetPostByID(postID int) (domain.Post, error) {
	var post domain.Post
	rows, err := p.db.Query(`select posts.id, posts.title, posts.preview, posts.date, posts.content, users.id, users.username, users.name, users.avatar, count(distinct views.user_id), count(distinct comments.post_id) from posts inner join users on posts.user_id = users.id left join views on posts.id = views.post_id left join comments on posts.id = comments.post_id where public = 1 and posts.id = ? group by posts.id`, postID)
	if err != nil {
		return post, err
	}

	for rows.Next() {
		err = rows.Scan(&post.ID, &post.Title, &post.Preview, &post.Date, &post.Content, &post.Author.ID, &post.Author.Username, &post.Author.Name, &post.Author.Avatar, &post.Counter.Views, &post.Counter.Comments)
		if err != nil {
			return post, err
		}
	}

	return post, nil
}

func (p post) GetDrafts(userID int) ([]*domain.Post, error) {
	var drafts []*domain.Post
	err := p.db.Select(&drafts, `select posts.id, posts.title, posts.preview, posts.date, posts.content from posts inner join users on posts.user_id = users.id where public = 0 and posts.user_id= ? group by posts.id`, userID)
	return drafts, err
}

func (p post) GetDraft(userID, postID int) (domain.Post, error) {
	var draft domain.Post
	err := p.db.Get(&draft, `select posts.id, posts.title, posts.preview, posts.date, posts.content from posts inner join users on posts.user_id = users.id where posts.id = ? and posts.user_id = ? and public = 0`, postID, userID)
	return draft, err
}

func (p post) GetPosts(queryMap *utils.QueryMap) ([]*domain.Post, error) {
	var posts []*domain.Post
	if queryMap.Queries["tag"] != "" {
		if queryMap.Queries["sort"] == "newest" {
			rows, err := p.db.Query(tagNewest, queryMap.Queries["tag"], queryMap.Pq.Limit, queryMap.Pq.Offset)
			if err != nil {
				return posts, err
			}

			for rows.Next() {
				var post domain.Post
				err = rows.Scan(&post.ID, &post.Title, &post.Preview, &post.Date, &post.Content, &post.Author.ID, &post.Author.Username, &post.Author.Name, &post.Author.Avatar, &post.Counter.Views, &post.Counter.Comments)
				if err != nil {
					return posts, err
				}
				posts = append(posts, &post)
			}

			return posts, nil
		}
		rows, err := p.db.Query(tagPopular, queryMap.Queries["tag"], queryMap.Pq.Limit, queryMap.Pq.Offset)
		if err != nil {
			return posts, err
		}

		for rows.Next() {
			var post domain.Post
			err = rows.Scan(&post.ID, &post.Title, &post.Preview, &post.Date, &post.Content, &post.Author.ID, &post.Author.Username, &post.Author.Name, &post.Author.Avatar, &post.Counter.Views, &post.Counter.Comments)
			if err != nil {
				return posts, err
			}
			posts = append(posts, &post)
		}

		return posts, nil
	} else {
		if queryMap.Queries["sort"] == "newest" {
			rows, err := p.db.Query(noTagNewest, queryMap.Pq.Limit, queryMap.Pq.Offset)
			if err != nil {
				return posts, err
			}

			for rows.Next() {
				var post domain.Post
				err = rows.Scan(&post.ID, &post.Title, &post.Preview, &post.Date, &post.Content, &post.Author.ID, &post.Author.Username, &post.Author.Name, &post.Author.Avatar, &post.Counter.Views, &post.Counter.Comments)
				if err != nil {
					return posts, err
				}
				posts = append(posts, &post)
			}

			return posts, nil
		}
		rows, err := p.db.Query(noTagPopular, queryMap.Pq.Limit, queryMap.Pq.Offset)
		if err != nil {
			return posts, err
		}

		for rows.Next() {
			var post domain.Post
			err = rows.Scan(&post.ID, &post.Title, &post.Preview, &post.Date, &post.Content, &post.Author.ID, &post.Author.Username, &post.Author.Name, &post.Author.Avatar, &post.Counter.Views, &post.Counter.Comments)
			if err != nil {
				return posts, err
			}
			posts = append(posts, &post)
		}

		return posts, nil
	}
}

func (p post) Publish(postID, userID int) (int, error) {
	sqlResult, err := p.db.Exec(`update posts set public = 1 where id = ? and user_id = ?`, postID, userID)
	if err != nil {
		return 0, err
	}
	lastID, err := sqlResult.RowsAffected()
	return int(lastID), err
}
