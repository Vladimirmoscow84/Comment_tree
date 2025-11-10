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

// GetSubtree возвращает комментарий и все дерево рекурсивно
func (p *Postgres) GetSubtree(ctx context.Context, id int) ([]model.Comment, error) {
	query := `
		WITH RECURSIVE subtree AS (
			SELECT * FROM comments WHERE id = $1
			UNION ALL
			SELECT c.* FROM comments c
			INNER JOIN subtree s ON c.parent_id = s.id
		)
		SELECT id, parent_id, content, created_at
		FROM subtree
		ORDER BY created_at ASC;
	`
	var comments []model.Comment
	if err := p.DB.SelectContext(ctx, &comments, query, id); err != nil {
		return nil, fmt.Errorf("[postgres]error of get subtree: %w", err)
	}
	return comments, nil
}

// DeleteSubtree удаляет комментарий и все дерево
func (p *Postgres) DeleteSubtree(ctx context.Context, id int) error {
	query := `
		WITH RECURSIVE subtree AS (
			SELECT id FROM comments WHERE id = $1
			UNION ALL
			SELECT c.id FROM comments c
			INNER JOIN subtree s ON c.parent_id = s.id
		)
		DELETE FROM comments WHERE id IN (SELECT id FROM subtree);
	`
	if _, err := p.DB.ExecContext(ctx, query, id); err != nil {
		return fmt.Errorf("delete subtree: %w", err)
	}
	return nil
}

// SearchComments выполняет полнотекстовый поиск по контенту
func (p *Postgres) SearchComments(ctx context.Context, queryText string, limit, offset int) ([]model.Comment, error) {
	query := `
		SELECT id, parent_id, content, created_at
		FROM comments
		WHERE search_vector @@ plainto_tsquery('simple', unaccent($1))
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3;
	`
	var results []model.Comment
	if err := p.DB.SelectContext(ctx, &results, query, queryText, limit, offset); err != nil {
		return nil, fmt.Errorf("search comments: %w", err)
	}
	return results, nil
}

// ListRootComments возвращает корневые комментарии с пагинацией
func (p *Postgres) ListRootComments(ctx context.Context, limit, offset int) ([]model.Comment, error) {
	query := `
		SELECT id, parent_id, content, created_at
		FROM comments
		WHERE parent_id IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2;
	`
	var comments []model.Comment
	if err := p.DB.SelectContext(ctx, &comments, query, limit, offset); err != nil {
		return nil, fmt.Errorf("list roots: %w", err)
	}
	return comments, nil
}
