package handler

import (
	"encoding/json"

	"github.com/brewinski/unnamed-fiber/data"
	"github.com/brewinski/unnamed-fiber/db"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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

	// decrypt all users in parallel
	userChannel := make(chan data.User, len(users))
	limit := make(chan struct{}, 10)
	for _, user := range users {
		limit <- struct{}{}
		go func(user data.User) {
			user.Decrypt()
			userChannel <- user
			<-limit
		}(user)
	}

	usersResponse := []data.UserResponse{}

	for user := range userChannel {
		data := data.UserResponse{}
		err = json.Unmarshal([]byte(user.User_Data), &data)
		if err != nil {
			return fiber.ErrInternalServerError
		}

		usersResponse = append(usersResponse, data)

		// err := UpdateUserData(data, user)
		// creditScore, err := CreateCreditScoreByUserID(user.ID)
		// err := DeleteCreditScoreByUserID(user.ID)
		// creditScore, err := GetCreditScoreByUserID(user.ID)
		// if err != nil {
		// 	return fiber.ErrInternalServerError
		// }

		// err = CreateUser()
		// if err != nil {
		// 	return fiber.ErrInternalServerError
		// }
	}

	return c.JSON(usersResponse)
}

func GetCreditScoreByUserID(userID string) (*data.Credit, error) {
	db := db.DB
	credit := &data.Credit{}

	err := db.Joins("User").Find(credit, "user_id = ?", userID).Error
	if err != nil {
		return nil, err
	}

	return credit, nil
}

func DeleteCreditScoreByUserID(userID string) error {
	db := db.DB
	credit := &data.Credit{}

	err := db.Delete(credit, "user_id = ?", userID).Error
	if err != nil {
		return err
	}

	return nil
}

func CreateCreditScoreByUserID(userID string) (*data.Credit, error) {
	db := db.DB
	err := db.Create(&data.Credit{Score: "100", User: data.User{ID: userID}}).Error
	if err != nil {
		return nil, err
	}

	credit := &data.Credit{}
	db.Find(credit, "user_id = ?", userID)

	return credit, nil
}

func UpdateUserDataHandler(c *fiber.Ctx) error {
	userRequest := &data.UserResponse{}
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

func ListUsers() ([]data.User, error) {
	db := db.DB
	users := []data.User{}

	err := db.Find(&users).Error
	if err != nil {
		return users, err
	}

	return users, nil
}

func GetUserByID(id string) (*data.User, error) {
	db := db.DB
	user := &data.User{}

	err := db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

func UpdateUserData(userRequest data.UserResponse, user data.User) error {
	db := db.DB
	updatedUserString, err := json.Marshal(userRequest)
	if err != nil {
		return err
	}

	user.User_Data = string(updatedUserString)
	db.Save(&user)

	return nil
}

func CreateUser() error {
	db := db.DB

	id, err := uuid.NewUUID()
	if err != nil {
		return err
	}

	userData := data.UserResponse{
		ID:                 id.String(),
		First_Name:         "Chris",
		Last_Name:          "string `json:\"last_name\"`",
		Nick_Name:          "string `json:\"nickname\"`",
		Provider:           "string `json:\"provider\"`",
		Signed_Up_From:     "string `json:\"signed_up_from\"`",
		Visitor_UUID:       "string `json:\"visitor_uuid\"`",
		Username:           "string `json:\"username\"`",
		Unsubscribe_Key:    "string `json:\"unsubscribe_key\"`",
		Created_Date:       "string `json:\"created_date\"`",
		Last_Modified_Date: "string `json:\"last_modified_date\"`",
		Last_Login_Date:    "string `json:\"last_login_date\"`",
		Accepted_Timestamp: "string `json:\"accepted_timestamp\"`",
	}

	jsonData, err := json.Marshal(userData)
	if err != nil {
		return err
	}

	err = db.Create(&data.User{
		ID:              id.String(),
		User_Data:       string(jsonData),
		Visitor_UUID:    id.String(),
		Unsubscribe_Key: id.String(),
	}).Error

	if err != nil {
		return err
	}

	return nil
}
