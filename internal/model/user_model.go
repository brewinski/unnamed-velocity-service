package model

import (
	"strings"

	"github.com/brewinski/unnamed-fiber/pkg/config"
	"github.com/brewinski/unnamed-fiber/pkg/envelope"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model               // Adds some metadata fields to the table
	ID                string `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"` // Explicitly specify the type to be uuid
	Visitor_UUID      string
	Unsubscribe_Key   string
	Ciphertext        string
	Credit_Ciphertext string
	Credit_Data       string
	User_Data         string
	Test_Field        string
	First_Name        string
	First_Name_2      string
	Last_Name         string
}

// setup gorm object lifecycle hooks
func (user *User) AfterFind(tx *gorm.DB) (err error) {
	// get the DEK from the KMS
	decryptedDek, err := envelope.ReadEncryptedDEK(config.Config("MASTER_KEY_USER_ENCRYPT_NAME"), user.Ciphertext)
	if err != nil {
		return err
	}

	userData, err := envelope.ReadEncryptedDataWithDEK(user.User_Data, string(decryptedDek.Plaintext))
	if err != nil {
		return err
	}

	if user.Test_Field != "" {
		testField, err := envelope.ReadEncryptedDataWithDEK(user.Test_Field, string(decryptedDek.Plaintext))
		if err != nil {
			return err
		}

		user.Test_Field = testField
	}

	if user.First_Name != "" {
		firstNameField, err := envelope.ReadEncryptedDataWithDEK(user.First_Name, string(decryptedDek.Plaintext))
		if err != nil {
			return err
		}

		user.First_Name = firstNameField
	}

	if user.First_Name_2 != "" {
		firstName2Field, err := envelope.ReadEncryptedDataWithDEK(user.First_Name, string(decryptedDek.Plaintext))
		if err != nil {
			return err
		}

		user.First_Name_2 = firstName2Field
	}

	if user.Last_Name != "" {
		lastNameField, err := envelope.ReadEncryptedDataWithDEK(user.Last_Name, string(decryptedDek.Plaintext))
		if err != nil {
			return err
		}

		user.Last_Name = lastNameField
	}
	// replace the user data with the decrypted data
	user.User_Data = userData

	return nil
}

func (user *User) BeforeSave(tx *gorm.DB) (err error) {
	// encrypt the user data
	dek, err := envelope.CreateNewDEK(strings.Split(config.Config("MASTER_KEY_USER_ENCRYPT_NAME"), "/keyRings/")[0])
	if err != nil {
		return err
	}

	encryptedDek, err := envelope.EncryptDEK(dek.Data, config.Config("MASTER_KEY_USER_ENCRYPT_NAME"))
	if err != nil {
		return err
	}

	encryptedUserData, err := envelope.EncryptDataWithDEK(dek.Data, user.User_Data)
	if err != nil {
		return err
	}

	encryptedTestField, err := envelope.EncryptDataWithDEK(dek.Data, user.Test_Field)
	if err != nil {
		return err
	}

	encryptedFirstName, err := envelope.EncryptDataWithDEK(dek.Data, user.First_Name)
	if err != nil {
		return err
	}

	encryptedFirstName2, err := envelope.EncryptDataWithDEK(dek.Data, user.First_Name)
	if err != nil {
		return err
	}

	encryptedLastName, err := envelope.EncryptDataWithDEK(dek.Data, user.Last_Name)
	if err != nil {
		return err
	}

	// set the DEK and encrypted user data
	user.Ciphertext = encryptedDek
	user.User_Data = encryptedUserData
	user.Test_Field = encryptedTestField
	user.First_Name = encryptedFirstName
	user.First_Name_2 = encryptedFirstName2
	user.Last_Name = encryptedLastName

	return nil
}

// This is the user data object is the response object for the user
type UserResponse struct {
	ID                 string `json:"id"`
	First_Name         string `json:"first_name"`
	Last_Name          string `json:"last_name"`
	Nick_Name          string `json:"nickname"`
	Provider           string `json:"provider"`
	Signed_Up_From     string `json:"signed_up_from"`
	Visitor_UUID       string `json:"visitor_uuid"`
	Username           string `json:"username"`
	Unsubscribe_Key    string `json:"unsubscribe_key"`
	Created_Date       string `json:"created_date"`
	Last_Modified_Date string `json:"last_modified_date"`
	Last_Login_Date    string `json:"last_login_date"`
	Accepted_Timestamp string `json:"accepted_timestamp"`
}
