package service

import (
	"discord-server-go/model"
	"discord-server-go/model/apperrors"
	"log"
	"mime/multipart"
)

type userService struct {
	UserRepository model.UserRepository
	FileRepository model.FileRepository
}

type USConfig struct {
	UserRepository model.UserRepository
	FileRepository model.FileRepository
}

func NewUserService(c *USConfig) model.UserService {
	return &userService{
		UserRepository: c.UserRepository,
		FileRepository: c.FileRepository,
	}
}
func (s *userService) Register(user *model.User) (*model.User, error) {
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		log.Printf("Unable to signup user for email: %v\n", user.Email)
		return nil, apperrors.NewInternal()
	}
	user.ID = GenerateId()
	user.Password = hashedPassword
	return s.UserRepository.Create(user)
}

func (s *userService) Login(email, password string) (*model.User, error) {
	user, err := s.UserRepository.FindByEmail(email)

	if err != nil {
		return nil, apperrors.NewAuthorization(apperrors.InvalidCredentials)
	}

	match, err := comparePasswords(user.Password, password)

	if err != nil {
		return nil, apperrors.NewInternal()
	}

	if !match {
		return nil, apperrors.NewAuthorization(apperrors.InvalidCredentials)
	}

	return user, nil
}

func (s *userService) Get(id string) (*model.User, error) {
	return s.UserRepository.FindByID(id)
}
func (s *userService) ChangeAvatar(header *multipart.FileHeader, directory string) (string, error) {

	url, _, err := s.FileRepository.UploadAvatar(header, directory)
	if err != nil {
		return "", err
	}

	return url, nil
}
func (s *userService) DeleteImage(key string) error {
	return s.FileRepository.DeleteImage(key)
}
