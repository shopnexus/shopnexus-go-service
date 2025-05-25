package storage

import (
	"context"
	"shopnexus-go-service/gen/sqlc"
	"shopnexus-go-service/internal/model"
	"slices"
)

func (r *ServiceImpl) GetResources(ctx context.Context, ownerID int64, resourceType model.ResourceType) ([]string, error) {
	return r.sqlc.GetResources(ctx, sqlc.GetResourcesParams{
		OwnerID: ownerID,
		Type:    sqlc.ProductResourceType(resourceType),
	})
}

func (r *ServiceImpl) AddResources(ctx context.Context, ownerID int64, resourceType model.ResourceType, resources []string) error {
	var args []sqlc.AddResourcesParams
	for i, resource := range resources {
		args = append(args, sqlc.AddResourcesParams{
			Type:    sqlc.ProductResourceType(resourceType),
			OwnerID: ownerID,
			Url:     resource,
			Order:   int32(i),
		})
	}

	_, err := r.sqlc.AddResources(ctx, args)
	return err
}

func (r *ServiceImpl) UpdateResources(ctx context.Context, ownerID int64, resourceType model.ResourceType, resources []string) error {
	current, err := r.GetResources(ctx, ownerID, resourceType)
	if err != nil {
		return err
	}

	if slices.Equal(current, resources) {
		return nil
	}

	if err := r.EmptyResources(ctx, ownerID, resourceType); err != nil {
		return err
	}

	return r.AddResources(ctx, ownerID, resourceType, resources)
}

func (r *ServiceImpl) EmptyResources(ctx context.Context, ownerID int64, resourceType model.ResourceType) error {
	return r.sqlc.EmptyResources(ctx, sqlc.EmptyResourcesParams{
		OwnerID: ownerID,
		Type:    sqlc.ProductResourceType(resourceType),
	})
}
