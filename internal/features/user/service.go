package user

type UserService struct {
	userRepo UserRepository
}

func NewUserService(userRepo UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetCurrentUserInformation(userID uint) (CurrentUserDataResponse, error) {
	user, err := s.userRepo.GetUserByUserID(userID)
	if err != nil {
		return CurrentUserDataResponse{}, err
	}

	roles, err := s.userRepo.GetUserRoleNamesByUserID(userID)
	if err != nil {
		return CurrentUserDataResponse{}, err
	}

	permissions, err := s.userRepo.GetUserPermissionNamesByUserID(userID)
	if err != nil {
		return CurrentUserDataResponse{}, err
	}

	response := CurrentUserDataResponse{
		ID:          user.ID,
		Email:       user.Email,
		FirstName:   user.Profile.FirstName,
		LastName:    user.Profile.LastName,
		Phone:       user.Profile.Phone,
		Gender:      user.Profile.Gender,
		Roles:       roles,
		Permissions: permissions,
	}

	return response, nil
}
