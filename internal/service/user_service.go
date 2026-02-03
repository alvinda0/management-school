package service

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"project-go/internal/model"
	"project-go/internal/repository"
	"project-go/utils"
)

type UserService interface {
	Register(req *model.UserRequest) (*model.UserResponse, error)
	Login(req *model.LoginRequest) (string, error)
	GetUserByID(id string) (*model.UserResponse, error)
	GetAllUsers() ([]model.UserResponse, error)
	UpdateUser(id string, req *model.UserRequest) (*model.UserResponse, error)
	DeleteUser(id string) error
}

type userService struct {
	userRepo repository.UserRepository
	roleRepo repository.RoleRepository
}

func NewUserService(userRepo repository.UserRepository, roleRepo repository.RoleRepository) UserService {
	return &userService{
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

func (s *userService) Register(req *model.UserRequest) (*model.UserResponse, error) {
	// Check if user already exists
	_, err := s.userRepo.GetByEmail(req.Email)
	if err == nil {
		return nil, errors.New("user already exists")
	}

	// Validate role exists
	_, err = s.roleRepo.GetByID(req.RoleID)
	if err != nil {
		return nil, errors.New("invalid role")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		RoleID:   req.RoleID,
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	// Get user with role for response
	userWithRole, err := s.userRepo.GetByID(user.ID)
	if err != nil {
		return nil, err
	}

	return s.mapToUserResponse(userWithRole), nil
}

func (s *userService) Login(req *model.LoginRequest) (string, error) {
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := utils.GenerateJWT(user.ID, user.RoleID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *userService) GetUserByID(id string) (*model.UserResponse, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return s.mapToUserResponse(user), nil
}

func (s *userService) GetAllUsers() ([]model.UserResponse, error) {
	users, err := s.userRepo.GetAll()
	if err != nil {
		return nil, err
	}

	var userResponses []model.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, *s.mapToUserResponse(&user))
	}

	return userResponses, nil
}

func (s *userService) UpdateUser(id string, req *model.UserRequest) (*model.UserResponse, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Validate role exists
	_, err = s.roleRepo.GetByID(req.RoleID)
	if err != nil {
		return nil, errors.New("invalid role")
	}

	user.Username = req.Username
	user.Email = req.Email
	user.RoleID = req.RoleID

	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.Password = string(hashedPassword)
	}

	err = s.userRepo.Update(user)
	if err != nil {
		return nil, err
	}

	// Get updated user with role
	updatedUser, err := s.userRepo.GetByID(user.ID)
	if err != nil {
		return nil, err
	}

	return s.mapToUserResponse(updatedUser), nil
}

func (s *userService) DeleteUser(id string) error {
	return s.userRepo.Delete(id)
}

func (s *userService) mapToUserResponse(user *model.User) *model.UserResponse {
	return &model.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role: model.RoleResponse{
			ID:          user.Role.ID,
			Name:        user.Role.Name,
			Description: user.Role.Description,
		},
	}
}