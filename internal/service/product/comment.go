package product

import (
	"context"
	"shopnexus-go-service/internal/model"
	"shopnexus-go-service/internal/repository"
)

func (s *ProductService) GetComment(ctx context.Context, id int64) (model.Comment, error) {
	comment, err := s.repo.GetComment(ctx, id)
	if err != nil {
		return model.Comment{}, err
	}

	return comment, nil
}

type ListCommentsParams = repository.ListCommentsParams

func (s *ProductService) ListComments(ctx context.Context, params ListCommentsParams) (result model.PaginateResult[model.Comment], err error) {
	total, err := s.repo.CountComments(ctx, params)
	if err != nil {
		return result, err
	}

	comments, err := s.repo.ListComments(ctx, params)
	if err != nil {
		return result, err
	}

	return model.PaginateResult[model.Comment]{
		Data:       comments,
		Limit:      params.Limit,
		Page:       params.Page,
		Total:      total,
		NextPage:   params.NextPage(total),
		NextCursor: nil,
	}, nil
}

type CreateCommentParams struct {
	AccountID int64
	Type      model.CommentType
	DestID    int64
	Body      string
	Resources []string
}

func (s *ProductService) CreateComment(ctx context.Context, params CreateCommentParams) (model.Comment, error) {
	txRepo, err := s.repo.Begin(ctx)
	if err != nil {
		return model.Comment{}, err
	}
	defer txRepo.Rollback(ctx)

	comment, err := txRepo.CreateComment(ctx, model.Comment{
		Type:      params.Type,
		AccountID: params.AccountID,
		DestID:    params.DestID,
		Body:      params.Body,
		Resources: params.Resources,
	})
	if err != nil {
		return model.Comment{}, err
	}

	if err := txRepo.Commit(ctx); err != nil {
		return model.Comment{}, err
	}

	return comment, nil
}

// TODO: always check user only modify their resources
type UpdateCommentParams struct {
	Role      model.Role
	AccountID int64
	ID        int64
	Body      *string
	Resources *[]string
}

func (s *ProductService) UpdateComment(ctx context.Context, params UpdateCommentParams) error {
	txRepo, err := s.repo.Begin(ctx)
	if err != nil {
		return err
	}
	defer txRepo.Rollback(ctx)

	repoParams := repository.UpdateCommentParams{
		ID:        params.ID,
		Body:      params.Body,
		Resources: params.Resources,
	}

	// User only can modify their own comment
	if params.Role == model.RoleUser {
		repoParams.AccountID = &params.AccountID
	}

	err = txRepo.UpdateComment(ctx, repoParams)
	if err != nil {
		return err
	}

	if err := txRepo.Commit(ctx); err != nil {
		return err
	}

	return nil
}

type DeleteCommentParams struct {
	Role      model.Role
	AccountID int64
	ID        int64
}

func (s *ProductService) DeleteComment(ctx context.Context, params DeleteCommentParams) error {
	repoParams := repository.DeleteCommentParams{
		ID: params.ID,
	}

	// User only can delete their own comment
	if params.Role == model.RoleUser {
		repoParams.AccountID = &params.AccountID
	}

	return s.repo.DeleteComment(ctx, repoParams)
}
