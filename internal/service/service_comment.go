package service

import (
	"comment_tree/internal/model"
	"context"
	"fmt"
	"strings"
)

// CreateComment  создает новый коментарий или ответ на родительский
func (s *Service) CreateComment(ctx context.Context, parentID *int, content string) (*model.Comment, error) {
	content = strings.TrimSpace(content)
	if content == "" {
		return nil, fmt.Errorf("[service] content can't be empty")
	}
	c, err := s.storage.CreateComment(ctx, parentID, content)
	if err != nil {
		return nil, fmt.Errorf("[service] error creating comment: %w", err)
	}
	return c, nil

}

// GetSubtree возвращает коментарий и древо коментариев при наличии
func (s *Service) GetSubtree(ctx context.Context, id int) ([]model.Comment, error) {
	tree, err := s.storage.GetSubtree(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("[service] error getting subtree: %w", err)
	}
	return tree, nil
}

// DeleteSubtree удаляет коментарий и все вложенные
func (s *Service) DeleteSubtree(ctx context.Context, id int) error {
	err := s.storage.DeleteSubtree(ctx, id)
	if err != nil {
		return fmt.Errorf("[service] error of deleting subtree: %w", err)
	}
	return nil
}

// SearchComments ищет коментарии по ключевым словам
func (s *Service) SearchComments(ctx context.Context, queryText string, limit, offset int) ([]model.Comment, error) {
	queryText = strings.TrimSpace(queryText)
	if queryText == "" {
		return nil, fmt.Errorf("[service] serach query can't be empty")
	}
	results, err := s.storage.SearchComments(ctx, queryText, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("[service] error serching comments: %w", err)
	}
	return results, nil
}

// ListRootComments возвращает корневые коментарии
func (s *Service) ListRootComments(ctx context.Context, limit, offset int) ([]model.Comment, error) {
	comments, err := s.storage.ListRootComments(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("[service] error listing root comments: %w", err)
	}
	return comments, nil
}
