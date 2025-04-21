package repository

import (
	"context"
	"shopnexus-go-service/gen/sqlc"
	pgxutil "shopnexus-go-service/internal/db/pgx"
	"shopnexus-go-service/internal/model"

	"github.com/jackc/pgx/v5/pgtype"
)

func (r *RepositoryImpl) GetComment(ctx context.Context, id int64) (model.Comment, error) {
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

func (r *RepositoryImpl) CountComments(ctx context.Context, params ListCommentsParams) (int64, error) {
	return r.sqlc.CountComments(ctx, sqlc.CountCommentsParams{
		AccountID:     *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.AccountID),
		Type:          *pgxutil.PtrBrandedToPgType(&sqlc.NullProductCommentType{}, params.Type),
		DestID:        *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.DestID),
		Body:          *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Body),
		UpvoteFrom:    *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.UpvoteFrom),
		UpvoteTo:      *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.UpvoteTo),
		DownvoteFrom:  *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.DownvoteFrom),
		DownvoteTo:    *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.DownvoteTo),
		ScoreFrom:     *pgxutil.PtrToPgtype(&pgtype.Int4{}, params.ScoreFrom),
		ScoreTo:       *pgxutil.PtrToPgtype(&pgtype.Int4{}, params.ScoreTo),
		CreatedAtFrom: *pgxutil.PtrToPgtype(&pgtype.Timestamptz{}, params.DateCreatedFrom),
		CreatedAtTo:   *pgxutil.PtrToPgtype(&pgtype.Timestamptz{}, params.DateCreatedTo),
	})
}

func (r *RepositoryImpl) ListComments(ctx context.Context, params ListCommentsParams) ([]model.Comment, error) {
	rows, err := r.sqlc.ListComments(ctx, sqlc.ListCommentsParams{
		Limit:         params.Limit,
		Offset:        params.Offset(),
		AccountID:     *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.AccountID),
		Type:          *pgxutil.PtrBrandedToPgType(&sqlc.NullProductCommentType{}, params.Type),
		DestID:        *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.DestID),
		Body:          *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Body),
		UpvoteFrom:    *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.UpvoteFrom),
		UpvoteTo:      *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.UpvoteTo),
		DownvoteFrom:  *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.DownvoteFrom),
		DownvoteTo:    *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.DownvoteTo),
		ScoreFrom:     *pgxutil.PtrToPgtype(&pgtype.Int4{}, params.ScoreFrom),
		ScoreTo:       *pgxutil.PtrToPgtype(&pgtype.Int4{}, params.ScoreTo),
		CreatedAtFrom: *pgxutil.PtrToPgtype(&pgtype.Timestamptz{}, params.DateCreatedFrom),
		CreatedAtTo:   *pgxutil.PtrToPgtype(&pgtype.Timestamptz{}, params.DateCreatedTo),
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

func (r *RepositoryImpl) CreateComment(ctx context.Context, comment model.Comment) (model.Comment, error) {
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

func (r *RepositoryImpl) UpdateComment(ctx context.Context, params UpdateCommentParams) error {
	row, err := r.sqlc.UpdateComment(ctx, sqlc.UpdateCommentParams{
		ID:        params.ID,
		AccountID: *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.AccountID),
		Body:      *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Body),
		Upvote:    *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.Upvote),
		Downvote:  *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.Downvote),
		Score:     *pgxutil.PtrToPgtype(&pgtype.Int4{}, params.Score),
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

func (r *RepositoryImpl) DeleteComment(ctx context.Context, params DeleteCommentParams) error {
	return r.sqlc.DeleteComment(ctx, sqlc.DeleteCommentParams{
		ID:        params.ID,
		AccountID: *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.AccountID),
	})
}
