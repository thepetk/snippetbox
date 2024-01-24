package models

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: no matching record found")

type Snippet struct {
	ID      int       `json:"id"`
	Title   string    `json:"title,omitempty"`
	Content string    `json:"content,omitempty"`
	Created time.Time `json:"created"`
	Expires time.Time `json:"expires"`
}

type SnippetService interface {
	Get(id string) (*Snippet, error)
	Create(snippet *Snippet) error
}

type SnippetModel struct {
	db           *sql.DB
	queryTimeout time.Duration
}

func (s *SnippetModel) SetDB(db *sql.DB) {
	s.db = db
	s.queryTimeout = 5 * time.Second
}

func (s *SnippetModel) Create(snippet *Snippet) (int, error) {
	query := `
        INSERT INTO snippets (title, content, created, expires) 
        VALUES ($1, $2, current_timestamp, current_timestamp + '365 days')
        RETURNING id, title, created`

	ctx, cancel := context.WithTimeout(context.Background(), s.queryTimeout)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, snippet.Title, snippet.Content).
		Scan(&snippet.ID, &snippet.Title, &snippet.Created)

	if err != nil {
		return 0, err
	}

	return snippet.ID, err
}

func (s *SnippetModel) Get(id string) (*Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
		WHERE expires > current_timestamp AND id=$1`

	row := s.db.QueryRow(stmt, id)

	sn := &Snippet{}

	err := row.Scan(&sn.ID, &sn.Title, &sn.Content, &sn.Created, &sn.Expires)
	if err == sql.ErrNoRows {
		return nil, ErrNoRecord
	} else if err != nil {
		return nil, err
	}

	return sn, nil
}

func (s *SnippetModel) Latest() ([]*Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > current_timestamp ORDER BY created DESC LIMIT 10`

	rows, err := s.db.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	snippets := []*Snippet{}
	for rows.Next() {
		sn := &Snippet{}
		err = rows.Scan(&sn.ID, &sn.Title, &sn.Content, &sn.Created, &sn.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, sn)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
