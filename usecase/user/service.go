package user

import (
	"go_rest_api_with_mysql/config"
	"go_rest_api_with_mysql/entity"
	"go_rest_api_with_mysql/infrastructure/repository"
	logger "go_rest_api_with_mysql/pkg/log"

	svcErr "go_rest_api_with_mysql/usecase/error"

	"go.uber.org/zap"
)

const (
	DEFAULTCOST int = 10
)

type Service struct {
	appConfig      *config.AppConfig
	userRepository *repository.UserMySQL
}

var log *zap.SugaredLogger = logger.GetLogger().Sugar()

// new service
func NewService(appConfig *config.AppConfig, userRepository *repository.UserMySQL) *Service {
	return &Service{
		appConfig:      appConfig,
		userRepository: userRepository,
	}
}

// createUser create an user
func (service *Service) CreateUser(userConfigs *entity.User) (entity.ID, error) {
	userData, err := entity.NewUser(userConfigs)
	if err != nil {
		return userData.ID, &svcErr.ServiceError{
			Err:  err,
			Code: svcErr.PROCESSING_ERROR,
		}
	}
	return service.userRepository.CreateUser(userData)
}

// get users
func (service *Service) GetUsers() ([]*entity.User, error) {
	userList, err := service.userRepository.GetUserList()
	if err != nil {
		log.Errorf("Error getting user list. %v", err)
		return nil, &svcErr.ServiceError{
			Err:  err,
			Code: svcErr.PROCESSING_ERROR,
		}
	}
	return userList, nil
}

// update an user
func (service *Service) UpdateUser(user *entity.User, userId entity.ID) error {
	err := user.Validate()
	if err != nil {
		log.Errorf("Error updating user details. %v", err)
		return &svcErr.ServiceError{
			Err:  err,
			Code: svcErr.PROCESSING_ERROR,
		}
	}
	return service.userRepository.UpdateUser(user, userId)
}

// get a user
func (service *Service) GetUser(id entity.ID) (*entity.User, error) {
	user, err := service.userRepository.GetUser(id)
	if err != nil {
		log.Errorf("Error getting user details. %v", err)
		return nil, &svcErr.ServiceError{
			Err:  err,
			Code: svcErr.PROCESSING_ERROR,
		}
	}
	return user, nil
}
