package postgres

import (
	"comment_tree/internal/model"
	"context"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

type Postgres struct {
	DB *sqlx.DB
}

// New  - конструктор БД
func New(databaseURI string) (*Postgres, error) {
	db, err := sqlx.Connect("pgx", databaseURI)
	if err != nil {
		return nil, fmt.Errorf("[postgres] failed to connect to DB: %w", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("[postgres] ping failed: %w", err)
	}
	log.Println("[postgres] connect to DB successfully")
	return &Postgres{
		DB: db,
	}, nil
}

// Close закрывает соединение с БД
func (p *Postgres) Close() error {
	if p.DB != nil {
		log.Println("[postgres] closing connection to DB")
		return p.DB.Close()
	}
	return nil
}

// CreateComment создает новый коментарий
func (p *Postgres) CreateComment(ctx context.Context, parentID *int, content string) (*model.Comment, error) {
	query := `
		INSERT INTO comments(parent_id, content)
		VALUES ($1, $2)	
		RETURNING id, parent_id, content, created_at;
	`
	var c model.Comment
	err := p.DB.GetContext(ctx, &c, query, parentID, content)
	if err != nil {
		return nil, fmt.Errorf("[postgres] failed creating comment: %w", err)
	}
	return &c, nil
}
