package usecase

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/afif0808/sagara-test/errors"
	"github.com/afif0808/sagara-test/internal/domain"
	"github.com/bwmarrin/snowflake"
	"golang.org/x/crypto/bcrypt"
)

type repository interface {
	InsertUser(ctx context.Context, u *domain.User) error
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
}

type UserUsecase struct {
	repo repository
}

func NewUserUsecase(repo repository) UserUsecase {
	return UserUsecase{repo: repo}
}
func (uu *UserUsecase) generatePasswordSalt(length int) string {
	seed := rand.New(rand.NewSource(time.Now().UnixNano()))
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	salt := make([]byte, length)
	for i := range salt {
		salt[i] = charset[seed.Intn(len(charset))]
	}
	return string(salt)
}

func (uu *UserUsecase) CreateUser(ctx context.Context, u *domain.User) error {
	_, err := uu.repo.GetUserByEmail(ctx, u.Email)

	if err == nil {
		return errors.ErrUserExists
	} else if err != errors.ErrNotFound {
		return err
	}

	node, err := snowflake.NewNode(1)
	if err != nil {
		return err
	}

	u.PasswordSalt = uu.generatePasswordSalt(5)
	log.Println(u.PasswordSalt)
	u.ID = node.Generate().Int64()
	password, err := bcrypt.GenerateFromPassword([]byte(u.PasswordSalt+u.Password), 12)
	if err != nil {
		return err
	}
	u.Password = string(password)
	return uu.repo.InsertUser(ctx, u)
}
