package service

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"srbolab_cpc/logoped"
	"srbolab_cpc/model"
	"strconv"
	"time"
)

var (
	AuthService authServiceInterface = &authService{}
)

type authService struct {
}

type authServiceInterface interface {
	GetUserIDByToken(token string) (int, error)
	Login(request model.LoginRequest) (model.LoginResponse, error)
}

func (s *authService) Login(request model.LoginRequest) (model.LoginResponse, error) {
	user, err := UsersService.GetUserByEmail(request.Email)
	if err != nil {
		logoped.ErrorLog.Println("Error loging " + err.Error())
		return model.LoginResponse{}, errors.New("korisnički email ili šifra nisu dobri")
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		logoped.ErrorLog.Println("Error login user, error comparing hashes: ", err)
		return model.LoginResponse{}, err
	}

	roles, err := RoleService.GetRolesByUserID(user.ID)
	if err != nil {
		logoped.ErrorLog.Println("Error login user, error getting roles: ", err)
		return model.LoginResponse{}, err
	}

	roleCodes := []string{}
	for _, r := range roles {
		roleCodes = append(roleCodes, r.Code)
	}

	claims := jwt.MapClaims{
		"Id":        strconv.Itoa(int(user.ID)),
		"ExpiresAt": time.Now().Add(time.Hour * 8).Unix(),
		"Roles":     roleCodes,
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := jwtToken.SignedString([]byte("secret"))
	if err != nil {
		logoped.ErrorLog.Println("Error login user, error signing token: ", err)
		return model.LoginResponse{}, err
	}

	return model.LoginResponse{Token: token, FirstName: user.Person.FirstName, LastName: user.Person.LastName, Roles: roleCodes}, nil
}

func (s *authService) GetUserIDByToken(token string) (int, error) {
	claims := jwt.MapClaims{}

	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		logoped.ErrorLog.Println("Error getting user by token, error parse claims: ", err)
		return 0, err
	}

	id, err := strconv.Atoi(claims["Id"].(string))
	if err != nil {
		logoped.ErrorLog.Println("Error getting user by token: ", err)
		return 0, err
	}

	return id, nil
}
