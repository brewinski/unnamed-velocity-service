package envelope

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"hash/crc32"
	"io"
	"strings"

	kms "cloud.google.com/go/kms/apiv1"
	"cloud.google.com/go/kms/apiv1/kmspb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func CreateNewDEK(location string) (*kmspb.GenerateRandomBytesResponse, error) {
	// Create the KMS client.
	ctx := context.Background()
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	req := &kmspb.GenerateRandomBytesRequest{
		Location:        location,
		LengthBytes:     32,
		ProtectionLevel: 2,
	}

	response, err := client.GenerateRandomBytes(ctx, req)
	if err != nil {
		return nil, err
	}

	println(response)

	return response, nil
}

func EncryptDataWithDEK(dek []byte, data string) (string, error) {
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	block, err := aes.NewCipher(dek)
	if err != nil {
		return "", err
	}

	plaintext := []byte(data)

	aead, err := cipher.NewGCMWithNonceSize(block, len(iv))
	if err != nil {
		return "", err
	}

	ciphertext := aead.Seal(nil, iv, plaintext, nil)
	tag := ciphertext[len(ciphertext)-aes.BlockSize:]
	ciphertext = ciphertext[:len(ciphertext)-aes.BlockSize]

	fmt.Println("ciphertext:", ciphertext)
	fmt.Println("tag:", tag)
	fmt.Print("iv:", iv)

	result := base64.StdEncoding.EncodeToString(iv) + ":" +
		base64.StdEncoding.EncodeToString(tag) + ":" +
		base64.StdEncoding.EncodeToString(ciphertext)

	return result, nil
}

func EncryptDEK(dek []byte, keyName string) (string, error) {
	// Create the client.
	ctx := context.Background()
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create kms client: %v", err)
	}
	defer client.Close()

	// base64-encode the dek
	encodedDEK := base64.StdEncoding.EncodeToString(dek)
	base64EncodedDek := []byte(encodedDEK)

	// Optional but recommended: Compute plaintext's CRC32C.
	crc32c := func(data []byte) uint32 {
		t := crc32.MakeTable(crc32.Castagnoli)
		return crc32.Checksum(data, t)
	}
	plaintextCRC32C := crc32c(base64EncodedDek)

	result, err := client.Encrypt(ctx, &kmspb.EncryptRequest{
		Name:            keyName,
		Plaintext:       base64EncodedDek,
		PlaintextCrc32C: wrapperspb.Int64(int64(plaintextCRC32C)),
	})
	if err != nil {
		return "", fmt.Errorf("failed to encrypt plaintext: %v", err)
	}

	// Optional, but recommended: perform integrity verification on result.
	// For more details on ensuring E2E in-transit integrity to and from Cloud KMS visit:
	// https://cloud.google.com/kms/docs/data-integrity-guidelines
	if !result.VerifiedPlaintextCrc32C {
		return "", fmt.Errorf("encrypt: request corrupted in-transit")
	}
	if int64(crc32c(result.Ciphertext)) != result.CiphertextCrc32C.Value {
		return "", fmt.Errorf("encrypt: response corrupted in-transit")
	}

	base64EncodedCiphertext := base64.StdEncoding.EncodeToString(result.Ciphertext)

	return base64EncodedCiphertext, nil
}

func ReadEncryptedDataWithDEK(data string, key string) (string, error) {
	encryptedParts := strings.Split(data, ":")
	encodedIv := encryptedParts[0]
	decodedIv, err := base64.StdEncoding.DecodeString(encodedIv)
	if err != nil {
		return "", err
	}
	encodedTag := encryptedParts[1]
	decodedTag, err := base64.StdEncoding.DecodeString(encodedTag)
	if err != nil {
		return "", err
	}
	encryptedText := encryptedParts[2]
	decodedEncryptedText, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", err
	}
	decodedKey, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(decodedKey)

	if err != nil {
		return "", err
	}

	aead, err := cipher.NewGCMWithNonceSize(block, len(decodedIv))
	if err != nil {
		return "", err
	}

	print(decodedTag)

	decrypted, err := aead.Open(nil, decodedIv, append(decodedEncryptedText[:], decodedTag[:]...), nil)
	if err != nil {
		return "", err
	}

	return string(decrypted), nil
}

func ReadEncryptedDEK(name string, encodedCiphertext string) (*kmspb.DecryptResponse, error) {
	decodedCiphertext, err := base64.StdEncoding.DecodeString(encodedCiphertext)
	if err != nil {
		return nil, fmt.Errorf("failed to base64 decode ciphertext: %v", err)
	}

	ciphertext := []byte(decodedCiphertext) // result of a symmetric encryption call

	// Create the client.
	ctx := context.Background()
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create kms client: %v", err)
	}
	defer client.Close()

	// Optional, but recommended: Compute ciphertext's CRC32C.
	crc32c := func(data []byte) uint32 {
		t := crc32.MakeTable(crc32.Castagnoli)
		return crc32.Checksum(data, t)
	}
	ciphertextCRC32C := crc32c(ciphertext)

	// Build the request.
	req := &kmspb.DecryptRequest{
		Name:             name,
		Ciphertext:       ciphertext,
		CiphertextCrc32C: wrapperspb.Int64(int64(ciphertextCRC32C)),
	}

	// Call the API.
	result, err := client.Decrypt(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt ciphertext: %v", err)
	}

	// Optional, but recommended: perform integrity verification on result.
	// For more details on ensuring E2E in-transit integrity to and from Cloud KMS visit:
	// https://cloud.google.com/kms/docs/data-integrity-guidelines
	if int64(crc32c(result.Plaintext)) != result.PlaintextCrc32C.Value {
		return nil, fmt.Errorf("decrypt: response corrupted in-transit")
	}

	fmt.Print("Decrypted plaintext: ", result.Plaintext)
	return result, nil
}
