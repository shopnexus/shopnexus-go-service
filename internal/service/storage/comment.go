package storage

import (
	"context"
	"shopnexus-go-service/gen/sqlc"
	"shopnexus-go-service/internal/model"

	"github.com/jackc/pgx/v5/pgtype"
)

func (r *ServiceImpl) GetComment(ctx context.Context, id int64) (model.Comment, error) {
	row, err := r.sqlc.GetComment(ctx, id)
	if err != nil {
		return model.Comment{}, err
	}

	return model.Comment{
		ID:          row.ID,
		Type:        model.CommentType(row.Type),
		AccountID:   row.AccountID,
		DestID:      row.DestID,
		Body:        row.Body,
		Upvote:      row.Upvote,
		Downvote:    row.Downvote,
		Score:       row.Score,
		DateCreated: row.DateCreated.Time.UnixMilli(),
		DateUpdated: row.DateUpdated.Time.UnixMilli(),
		Resources:   row.Resources,
	}, nil
}

type ListCommentsParams struct {
	model.PaginationParams
	AccountID       *int64
	Type            *model.CommentType
	DestID          *int64
	Body            *string
	UpvoteFrom      *int64
	UpvoteTo        *int64
	DownvoteFrom    *int64
	DownvoteTo      *int64
	ScoreFrom       *int64
	ScoreTo         *int64
	DateCreatedFrom *int64
	DateCreatedTo   *int64
}

func (r *ServiceImpl) CountComments(ctx context.Context, params ListCommentsParams) (int64, error) {
	return r.sqlc.CountComments(ctx, sqlc.CountCommentsParams{
		AccountID:     *PtrToPgtype(&pgtype.Int8{}, params.AccountID),
		Type:          *PtrBrandedToPgType(&sqlc.NullProductCommentType{}, params.Type),
		DestID:        *PtrToPgtype(&pgtype.Int8{}, params.DestID),
		Body:          *PtrToPgtype(&pgtype.Text{}, params.Body),
		UpvoteFrom:    *PtrToPgtype(&pgtype.Int8{}, params.UpvoteFrom),
		UpvoteTo:      *PtrToPgtype(&pgtype.Int8{}, params.UpvoteTo),
		DownvoteFrom:  *PtrToPgtype(&pgtype.Int8{}, params.DownvoteFrom),
		DownvoteTo:    *PtrToPgtype(&pgtype.Int8{}, params.DownvoteTo),
		ScoreFrom:     *PtrToPgtype(&pgtype.Int4{}, params.ScoreFrom),
		ScoreTo:       *PtrToPgtype(&pgtype.Int4{}, params.ScoreTo),
		CreatedAtFrom: *PtrToPgtype(&pgtype.Timestamptz{}, params.DateCreatedFrom),
		CreatedAtTo:   *PtrToPgtype(&pgtype.Timestamptz{}, params.DateCreatedTo),
	})
}

func (r *ServiceImpl) ListComments(ctx context.Context, params ListCommentsParams) ([]model.Comment, error) {
	rows, err := r.sqlc.ListComments(ctx, sqlc.ListCommentsParams{
		Limit:         params.Limit,
		Offset:        params.Offset(),
		AccountID:     *PtrToPgtype(&pgtype.Int8{}, params.AccountID),
		Type:          *PtrBrandedToPgType(&sqlc.NullProductCommentType{}, params.Type),
		DestID:        *PtrToPgtype(&pgtype.Int8{}, params.DestID),
		Body:          *PtrToPgtype(&pgtype.Text{}, params.Body),
		UpvoteFrom:    *PtrToPgtype(&pgtype.Int8{}, params.UpvoteFrom),
		UpvoteTo:      *PtrToPgtype(&pgtype.Int8{}, params.UpvoteTo),
		DownvoteFrom:  *PtrToPgtype(&pgtype.Int8{}, params.DownvoteFrom),
		DownvoteTo:    *PtrToPgtype(&pgtype.Int8{}, params.DownvoteTo),
		ScoreFrom:     *PtrToPgtype(&pgtype.Int4{}, params.ScoreFrom),
		ScoreTo:       *PtrToPgtype(&pgtype.Int4{}, params.ScoreTo),
		CreatedAtFrom: *PtrToPgtype(&pgtype.Timestamptz{}, params.DateCreatedFrom),
		CreatedAtTo:   *PtrToPgtype(&pgtype.Timestamptz{}, params.DateCreatedTo),
	})
	if err != nil {
		return nil, err
	}

	comments := make([]model.Comment, 0, len(rows))
	for _, row := range rows {
		comments = append(comments, model.Comment{
			ID:          row.ID,
			Type:        model.CommentType(row.Type),
			AccountID:   row.AccountID,
			DestID:      row.DestID,
			Body:        row.Body,
			Upvote:      row.Upvote,
			Downvote:    row.Downvote,
			Score:       row.Score,
			DateCreated: row.DateCreated.Time.UnixMilli(),
			DateUpdated: row.DateUpdated.Time.UnixMilli(),
			Resources:   row.Resources,
		})
	}

	return comments, nil
}

func (r *ServiceImpl) CreateComment(ctx context.Context, comment model.Comment) (model.Comment, error) {
	row, err := r.sqlc.CreateComment(ctx, sqlc.CreateCommentParams{
		AccountID: comment.AccountID,
		Type:      sqlc.ProductCommentType(comment.Type),
		DestID:    comment.DestID,
		Body:      comment.Body,
		Upvote:    comment.Upvote,
		Downvote:  comment.Downvote,
		Score:     comment.Score,
	})
	if err != nil {
		return model.Comment{}, err
	}

	if err = r.AddResources(ctx, row.ID, model.ResourceTypeComment, comment.Resources); err != nil {
		return model.Comment{}, err
	}

	return model.Comment{
		ID:          row.ID,
		Type:        model.CommentType(row.Type),
		AccountID:   row.AccountID,
		DestID:      row.DestID,
		Body:        row.Body,
		Upvote:      row.Upvote,
		Downvote:    row.Downvote,
		Score:       row.Score,
		DateCreated: row.DateCreated.Time.UnixMilli(),
		DateUpdated: row.DateUpdated.Time.UnixMilli(),
		Resources:   comment.Resources,
	}, nil
}

type UpdateCommentParams struct {
	ID        int64
	AccountID *int64
	Body      *string
	Upvote    *int64
	Downvote  *int64
	Score     *int64
	Resources *[]string
}

func (r *ServiceImpl) UpdateComment(ctx context.Context, params UpdateCommentParams) error {
	row, err := r.sqlc.UpdateComment(ctx, sqlc.UpdateCommentParams{
		ID:        params.ID,
		AccountID: *PtrToPgtype(&pgtype.Int8{}, params.AccountID),
		Body:      *PtrToPgtype(&pgtype.Text{}, params.Body),
		Upvote:    *PtrToPgtype(&pgtype.Int8{}, params.Upvote),
		Downvote:  *PtrToPgtype(&pgtype.Int8{}, params.Downvote),
		Score:     *PtrToPgtype(&pgtype.Int4{}, params.Score),
	})
	if err != nil {
		return err
	}

	if params.Resources != nil {
		if err = r.UpdateResources(ctx, row.ID, model.ResourceTypeComment, *params.Resources); err != nil {
			return err
		}
	}

	return nil
}

type DeleteCommentParams struct {
	ID        int64
	AccountID *int64
}

func (r *ServiceImpl) DeleteComment(ctx context.Context, params DeleteCommentParams) error {
	return r.sqlc.DeleteComment(ctx, sqlc.DeleteCommentParams{
		ID:        params.ID,
		AccountID: *PtrToPgtype(&pgtype.Int8{}, params.AccountID),
	})
}
