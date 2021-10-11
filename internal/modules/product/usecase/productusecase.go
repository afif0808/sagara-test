package usecase

import (
	"context"
	"time"

	"github.com/afif0808/sagara-test/dataquery"
	"github.com/afif0808/sagara-test/internal/domain"
	"github.com/afif0808/sagara-test/meta"
	"github.com/bwmarrin/snowflake"
)

type repository interface {
	GetProductList(ctx context.Context, dq dataquery.DataQuery) (data []domain.Product, err error)
	CountProduct(ctx context.Context, dq dataquery.DataQuery) (count int, err error)
	GetProduct(ctx context.Context, id int64) (domain.Product, error)
	InsertProduct(ctx context.Context, p *domain.Product) error
	UpdateProduct(ctx context.Context, p *domain.Product) error
	DeleteProduct(ctx context.Context, id int64) error
}

type ProductUsecase struct {
	repo repository
}

func NewProductUsecase(repo repository) ProductUsecase {
	return ProductUsecase{repo: repo}
}

func (pu *ProductUsecase) CreateProduct(ctx context.Context, p *domain.Product) error {
	node, err := snowflake.NewNode(1)
	if err != nil {
		return err
	}
	p.ID = node.Generate().Int64()
	p.CreatedAt = time.Now()

	return pu.repo.InsertProduct(ctx, p)
}

func (pu *ProductUsecase) UpdateProduct(ctx context.Context, p *domain.Product) error {
	existing, err := pu.repo.GetProduct(ctx, p.ID)
	if err != nil {
		return err
	}

	existing.Name = p.Name
	existing.ImageURL = p.ImageURL

	err = pu.repo.UpdateProduct(ctx, &existing)
	if err != nil {
		return err
	}

	*p = existing

	return nil
}

func (pu *ProductUsecase) GetProduct(ctx context.Context, id int64) (domain.Product, error) {
	return pu.repo.GetProduct(ctx, id)
}
func (pu *ProductUsecase) GetProductList(ctx context.Context, dq dataquery.DataQuery) ([]domain.Product, meta.Meta, error) {
	data, err := pu.repo.GetProductList(ctx, dq)
	if err != nil {
		return nil, meta.Meta{}, err
	}
	dq.ShowAll = true
	count, err := pu.repo.CountProduct(ctx, dq)
	if err != nil {
		return nil, meta.Meta{}, err
	}

	meta := meta.NewMeta(dq.Page, dq.Limit, count)

	return data, meta, nil
}

func (pu *ProductUsecase) DeleteProduct(ctx context.Context, id int64) error {
	_, err := pu.GetProduct(ctx, id)
	if err != nil {
		return err
	}
	return pu.repo.DeleteProduct(ctx, id)
}
