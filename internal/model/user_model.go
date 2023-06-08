package model

import (
	"log"
	"strings"
	"time"

	"github.com/brewinski/unnamed-fiber/pkg/config"
	"github.com/brewinski/unnamed-fiber/pkg/envelope"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model               // Adds some metadata fields to the table
	ID                string `gorm:"primary_key;"` // Explicitly specify the type to be uuid
	Visitor_UUID      string
	Unsubscribe_Key   string
	Ciphertext        string `enc:"type:dek;key_name:USER_MASTER_DECRYPT_KEY"` // Encrypt this field with the specified algorithm
	Credit_Ciphertext string
	Credit_Data       string
	User_Data         string
	Test_Field        string
	First_Name        string `enc:"type:val;dek:Ciphertext"`
	First_Name_2      string `enc:"type:val;dek:Ciphertext"`
	Last_Name         string `enc:"type:val;dek:Ciphertext"`
}

// setup gorm object lifecycle hooks
func (user *User) Decrypt() (err error) {
	start := time.Now()
	// get the DEK from the KMS
	decryptedDek, err := envelope.ReadEncryptedDEK(config.Config("MASTER_KEY_USER_ENCRYPT_NAME"), user.Ciphertext)
	if err != nil {
		return err
	}

	userData, err := envelope.ReadEncryptedDataWithDEK(user.User_Data, string(decryptedDek.Plaintext))
	if err != nil {
		return err
	}

	// replace the user data with the decrypted data
	user.User_Data = userData

	elapsed := time.Since(start)
	log.Printf("User Decrypt took %s", elapsed)
	return nil
}

func (user *User) Encrypt() (err error) {
	start := time.Now()
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

	// set the DEK and encrypted user data
	user.Ciphertext = encryptedDek
	user.User_Data = encryptedUserData
	elapsed := time.Since(start)
	log.Printf("User Encrypt took %s", elapsed)
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

type Credit struct {
	gorm.Model
	Score      string
	Ciphertext string `enc:"type:dek;key_name:USER_MASTER_DECRYPT_KEY"`
	UserID     string
	User       User `gorm:"references:id"`
}

// setup gorm object lifecycle hooks
func (credit *Credit) AfterFind(tx *gorm.DB) (err error) {
	start := time.Now()
	// get the DEK from the KMS
	decryptedDek, err := envelope.ReadEncryptedDEK(config.Config("MASTER_KEY_USER_ENCRYPT_NAME"), credit.Ciphertext)
	if err != nil {
		return err
	}
	score, err := envelope.ReadEncryptedDataWithDEK(credit.Score, string(decryptedDek.Plaintext))
	if err != nil {
		return err
	}
	// replace the user data with the decrypted data
	credit.Score = score
	elapsed := time.Since(start)
	log.Printf("Credit Decrypt took %s", elapsed)
	return nil
}

func (credit *Credit) BeforeSave(tx *gorm.DB) (err error) {
	// encrypt the user data
	dek, err := envelope.CreateNewDEK(strings.Split(config.Config("MASTER_KEY_USER_ENCRYPT_NAME"), "/keyRings/")[0])
	if err != nil {
		return err
	}
	encryptedDek, err := envelope.EncryptDEK(dek.Data, config.Config("MASTER_KEY_USER_ENCRYPT_NAME"))
	if err != nil {
		return err
	}
	encryptedScore, err := envelope.EncryptDataWithDEK(dek.Data, credit.Score)
	if err != nil {
		return err
	}
	credit.Score = encryptedScore
	credit.Ciphertext = encryptedDek

	return nil
}
