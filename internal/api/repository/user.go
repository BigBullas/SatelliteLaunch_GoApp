package repository

import (
	"context"
	"RIP_lab1/internal/models"
)

func (r *Repository) SignUp(ctx context.Context, newUser models.User) (int, bool, error) {
	err := r.db.Create(&newUser).Error
	if err != nil {
		return 0, false, err
	}
	return newUser.UserId, newUser.IsAdmin, nil
}

func (r *Repository) ChangeProfile(ctx context.Context, changedUser models.User) (error) {
	var oldUser models.User
	result := r.db.First(&oldUser, "user_id =?", changedUser.UserId)
	if result.Error != nil {
		return result.Error
	}

	if changedUser.Login != "" {
		oldUser.Login = changedUser.Login
	}

	if changedUser.Email != "" {
		oldUser.Email = changedUser.Email
	}

	if changedUser.Password != "" {
		oldUser.Password = changedUser.Password
	}

	result = r.db.Save(oldUser)
	return result.Error
}

func (r *Repository) GetByCredentials(ctx context.Context, user models.User) (models.User, error) {
	err := r.db.First(&user, "login = ? AND password = ?", user.Login, user.Password).Error
	return user, err
}

func (r *Repository) GetUserInfo(ctx context.Context, user models.User) (models.User, error) {
	err := r.db.First(&user, "user_id = ?", user.UserId).Error
	return user, err
}
