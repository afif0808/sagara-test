package repository

import (
	"context"
	"database/sql"
	"strings"

	"github.com/afif0808/sagara-test/errors"
	"github.com/afif0808/sagara-test/internal/domain"

	"github.com/afif0808/sagara-test/structs"
	"github.com/jmoiron/sqlx"
)

const (
	dbName = "users"
)

type UserSQLRepository struct {
	readDB, writeDB *sqlx.DB
}

func NewUserSQLRepository(readDB, writeDB *sqlx.DB) UserSQLRepository {
	return UserSQLRepository{readDB: readDB, writeDB: writeDB}
}

func (repo *UserSQLRepository) InsertUser(ctx context.Context, u *domain.User) error {
	var sql strings.Builder
	sql.WriteString("INSERT INTO " + dbName + " ")
	fields := structs.GetStructTagValues(domain.User{}, "db")
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

func (repo *UserSQLRepository) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	var user domain.User
	err := repo.readDB.Get(&user, "SELECT * FROM "+dbName+" WHERE email = ?", email)
	if err == sql.ErrNoRows {
		return domain.User{}, errors.ErrNotFound
	}
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}
