package usecase

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/afif0808/sagara-test/internal/domain"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var (
	jwtSigningMethod = jwt.SigningMethodHS256
	jwtSecretKey     = os.Getenv("JWT_SECRET_KEY")
)

type repository interface {
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
}

type AuthUsecase struct {
	repo repository
}

func NewAuthUsecase(repo repository) AuthUsecase {
	return AuthUsecase{
		repo: repo,
	}
}

func (au *AuthUsecase) generateJWT(u domain.User) (string, error) {
	expireAt := time.Now().Add(time.Hour * 48)
	claims := struct {
		UserID    string `json:"user_id"`
		UserName  string `json:"user_name"`
		UserEmail string `json:"user_email"`
		jwt.StandardClaims
	}{
		StandardClaims: jwt.StandardClaims{ExpiresAt: expireAt.Unix()},
		UserID:         strconv.FormatInt(u.ID, 10),
		UserName:       u.Name,
		UserEmail:      u.Email,
	}
	token, err := jwt.NewWithClaims(jwtSigningMethod, claims).SignedString([]byte(jwtSecretKey))
	if err != nil {
		return "", err
	}
	return token, nil

}

func (au *AuthUsecase) Login(ctx context.Context, crd domain.LoginCredentials) (domain.User, string, error) {
	var token string
	user, err := au.repo.GetUserByEmail(ctx, crd.Identity)
	if err != nil {
		return domain.User{}, "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(user.PasswordSalt+crd.Password))
	if err != nil {
		return domain.User{}, "", errors.New("password invalid")
	}
	token, err = au.generateJWT(user)
	if err != nil {
		return domain.User{}, "", err
	}
	return user, token, nil
}

func (au *AuthUsecase) validateJWT(token string) (domain.User, error) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwtSigningMethod {
			return nil, errors.New("signing method invalid")
		}
		return []byte(jwtSecretKey), nil
	})
	if err != nil {
		return domain.User{}, err
	}
	if !t.Valid {
		return domain.User{}, errors.New("token is invalid")
	}
	claims := t.Claims.(jwt.MapClaims)
	userID, err := strconv.ParseInt(fmt.Sprint(claims["user_id"]), 10, 64)
	if err != nil {
		return domain.User{}, errors.New("user id in claim is expected to be integer")
	}

	user := domain.User{
		ID:    userID,
		Name:  fmt.Sprint(claims["user_name"]),
		Email: fmt.Sprint(claims["user_email"]),
	}
	return user, nil
}

func (au *AuthUsecase) Authenticate(ctx context.Context, token string) (domain.User, error) {
	return au.validateJWT(token)
}
