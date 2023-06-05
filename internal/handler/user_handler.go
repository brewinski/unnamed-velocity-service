package handler

import (
	"encoding/json"

	"github.com/brewinski/unnamed-fiber/internal/model"
	"github.com/brewinski/unnamed-fiber/platform/database"
	"github.com/gofiber/fiber/v2"
)

// This handler should abstract saving a note from the fiber fameowrk implementation details.
// Keep logic portable by extracting the request values we need and passing them to the worker functions.

func ListUsersHandler(c *fiber.Ctx) error {
	users, err := ListUsers()
	if err != nil {
		return fiber.ErrInternalServerError
	}

	if len(users) < 1 {
		return fiber.ErrNotFound
	}

	usersResponse := []model.UserResponse{}

	for _, user := range users {
		data := model.UserResponse{}
		err = json.Unmarshal([]byte(user.User_Data), &data)
		if err != nil {
			return err
		}
		data.First_Name = "Chris"
		usersResponse = append(usersResponse, data)
	}

	if err != nil {
		return err
	}

	return c.JSON(usersResponse)
}

func UpdateUserDataHandler(c *fiber.Ctx) error {
	userRequest := new(model.UserResponse)
	err := c.BodyParser(userRequest)
	if err != nil {
		return err
	}

	user, err := GetUserByID(c.Params("uuid"))
	if err != nil {
		return err
	}

	err = UpdateUserData(*userRequest, *user)
	if err != nil {
		return err
	}

	return c.JSON(userRequest)
}

func ListUsers() ([]model.User, error) {
	db := database.DB
	users := []model.User{}

	err := db.Find(&users).Error

	if err != nil {
		return users, err
	}

	return users, nil
}

func GetUserByID(id string) (*model.User, error) {
	db := database.DB
	user := &model.User{}

	err := db.Where("id = ?", id).First(&user).Error

	if err != nil {
		return user, err
	}

	return user, nil
}

func UpdateUserData(userRequest model.UserResponse, user model.User) error {
	db := database.DB
	updatedUserString, err := json.Marshal(userRequest)
	if err != nil {
		return err
	}

	user.User_Data = string(updatedUserString)
	db.Save(&user)

	return nil
}
