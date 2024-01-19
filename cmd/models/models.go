package models

import (
	"context"
	"database/sql"
	"time"

	"github.com/gofrs/uuid"
)

type Snippet struct {
	ID      string    `json:"id"`
	Title   string    `json:"title,omitempty"`
	Content string    `json:"content,omitempty"`
	Created time.Time `json:"created"`
	Expires time.Time `json:"expires"`
}

type SnippetService interface {
	Get(id string) (*Snippet, error)
	Create(snippet *Snippet) error
}

type PostgresSnippetService struct {
	db           *sql.DB
	queryTimeout time.Duration
}

func (p *PostgresSnippetService) Create(snippet *Snippet) error {
	query := `
        INSERT INTO snippets (title, content, created, expires) 
        VALUES ($1, $2, current_timestamp, current_timestamp + '365 days')
        RETURNING id, title, created`

	newV4, err := uuid.NewV4()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), p.queryTimeout)
	defer cancel()

	return p.db.QueryRowContext(ctx, query, newV4.String(), snippet.Title, snippet.Content).
		Scan(&snippet.ID, &snippet.Title, &snippet.Created)
}

func (p *PostgresSnippetService) Get(id string) (*Snippet, error) {
	query := `
        SELECT id, title, content, created, expires
        FROM snippets
        WHERE id = $1`

	var snippet Snippet

	ctx, cancel := context.WithTimeout(context.Background(), p.queryTimeout)
	defer cancel()

	err := p.db.QueryRowContext(ctx, query, id).Scan(
		&snippet.ID,
		&snippet.Title,
		&snippet.Content,
		&snippet.Created,
		&snippet.Expires,
	)

	if err != nil {
		return nil, err
	}

	return &snippet, nil
}
