// работа с БД PostgreSQL
package storage

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4/pgxpool"
)

// База данных.
type DB struct {
	pool *pgxpool.Pool
}

// Публикация, получаемая из RSS.
type Article struct {
	ID        int    // номер записи
	Title     string // заголовок публикации
	Content   string // содержание публикации
	PubTime   int64  // время публикации
	Url       string // ссылка на источник
	Publisher string // название источника
	Autor     string // Имя автора
}

func New() (*DB, error) {
	//connstr := os.Getenv("agrigatordb")
	connstr := "user=postgres password=plazma dbname=agrigatordb sslmode=disable"
	if connstr == "" {
		return nil, errors.New("не указано подключение к БД")
	}
	pool, err := pgxpool.Connect(context.Background(), connstr)
	if err != nil {
		return nil, err
	}
	db := DB{
		pool: pool,
	}
	return &db, nil
}

func (db *DB) SaveArticles(articles []Article) error {
	for _, article := range articles {
		_, err := db.pool.Exec(context.Background(), `
		INSERT INTO articles (title, content, pub_time, a_url, publisher, author)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (title) DO NOTHING`,
			article.Title,
			article.Content,
			article.PubTime,
			article.Url,
			article.Publisher,
			article.Autor,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// Articles возвращает последние новости из БД.
func (db *DB) LastArticles(n int) ([]Article, error) {
	if n == 0 {
		n = 10
	}
	rows, err := db.pool.Query(context.Background(), `
	SELECT id, title, content, pub_time, a_url, publisher, author FROM articles
	ORDER BY pub_time DESC
	LIMIT $1
	`,
		n,
	)
	if err != nil {
		return nil, err
	}
	var alist []Article
	for rows.Next() {
		var p Article
		err = rows.Scan(
			&p.ID,
			&p.Title,
			&p.Content,
			&p.PubTime,
			&p.Url,
			&p.Publisher,
			&p.Autor,
		)
		if err != nil {
			return nil, err
		}
		alist = append(alist, p)
	}
	return alist, rows.Err()
}
