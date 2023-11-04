package dataprocesslayer

import "github.com/thewhiterabbit1994/hangman-api/models"

func ConvertToUserDto(u *models.User) models.UserDto {
	return models.UserDto{ID: u.ID, Name: u.Name, Email: u.Email, CreatedAt: u.CreatedAt}
}
