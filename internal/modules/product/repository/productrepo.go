package repository

import (
	"context"
	"database/sql"
	"strconv"
	"strings"

	"github.com/afif0808/sagara-test/dataquery"
	"github.com/afif0808/sagara-test/errors"
	"github.com/afif0808/sagara-test/internal/domain"
	"github.com/afif0808/sagara-test/structs"
	"github.com/jmoiron/sqlx"
)

const (
	dbName = "products"
)

type ProductSQLRepository struct {
	readDB, writeDB *sqlx.DB
}

func NewProductSQLRepository(readDB, writeDB *sqlx.DB) ProductSQLRepository {
	return ProductSQLRepository{
		readDB:  readDB,
		writeDB: writeDB,
	}
}

func (repo *ProductSQLRepository) InsertProduct(ctx context.Context, u *domain.Product) error {
	var sql strings.Builder
	sql.WriteString("INSERT INTO " + dbName + " ")
	fields := structs.GetStructTagValues(domain.Product{}, "db")
	sql.WriteString("(" + strings.Join(fields, ",") + ")")
	for i, f := range fields {
		fields[i] = ":" + f
	}
	sql.WriteString("VALUES(" + strings.Join(fields, ",") + ")")
	_, err := repo.writeDB.NamedExecContext(ctx, sql.String(), u)
	if err != nil {
		return err
	}

	return nil
}

func (repo *ProductSQLRepository) UpdateProduct(ctx context.Context, u *domain.Product) error {
	var sql strings.Builder
	sql.WriteString("UPDATE " + dbName + " SET ")
	fields := structs.GetStructTagValues(domain.Product{}, "db")
	for _, f := range fields {
		if f == "id" || f == "created_at" {
			continue
		}
		sql.WriteString(f + "=:" + f + ",")
	}

	_, err := repo.writeDB.NamedExecContext(ctx, strings.TrimSuffix(sql.String(), ","), u)
	if err != nil {
		return err
	}

	return nil
}
func (repo *ProductSQLRepository) DeleteProduct(ctx context.Context, id int64) error {
	_, err := repo.writeDB.Exec("DELETE FROM "+dbName+" WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}

func (repo *ProductSQLRepository) GetProduct(ctx context.Context, id int64) (domain.Product, error) {
	var p domain.Product
	err := repo.readDB.Get(&p, "SELECT * FROM "+dbName+" WHERE id = ?", id)
	if err == sql.ErrNoRows {
		return domain.Product{}, errors.ErrNotFound
	}
	if err != nil {
		return domain.Product{}, err
	}
	return p, nil
}

func (repo *ProductSQLRepository) GetProductList(ctx context.Context, dq dataquery.DataQuery) ([]domain.Product, error) {
	var data []domain.Product
	sql := strings.Builder{}
	dq.Search = "'%" + dq.Search + "%'"
	sql.WriteString("SELECT * FROM " + dbName + " WHERE name LIKE  " + dq.Search)
	if dq.OrderBy == "" {
		dq.OrderBy = "id"
	}
	if dq.Sort == "" {
		dq.Sort = "DESC"
	}

	sql.WriteString(" ORDER BY " + dq.OrderBy + " " + dq.Sort)

	if !dq.ShowAll {
		sql.WriteString(" LIMIT " + strconv.Itoa(dq.Limit) + " OFFSET " + strconv.Itoa(dq.CalculateOffset()))
	}
	err := repo.readDB.SelectContext(ctx, &data, sql.String())
	return data, err
}
func (repo *ProductSQLRepository) CountProduct(ctx context.Context, dq dataquery.DataQuery) (int, error) {
	var count int
	sql := strings.Builder{}
	dq.Search = "'%" + dq.Search + "%'"
	sql.WriteString("SELECT count(*) FROM " + dbName + " WHERE name LIKE " + dq.Search)

	if dq.OrderBy == "" {
		dq.OrderBy = "id"
	}
	if dq.Sort == "" {
		dq.Sort = "DESC"
	}

	sql.WriteString(" ORDER BY " + dq.OrderBy + " " + dq.Sort)

	if !dq.ShowAll {
		sql.WriteString(" LIMIT " + strconv.Itoa(dq.Limit) + " OFFSET " + strconv.Itoa(dq.CalculateOffset()))
	}
	err := repo.readDB.GetContext(ctx, &count, sql.String())
	return count, err
}
