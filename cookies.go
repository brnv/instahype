package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"database/sql"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/pbkdf2"
)

func getSessionFromCookies(cookiesPath string) (string, error) {
	var (
		encryptedSessionID string
	)

	db, err := sql.Open("sqlite3", cookiesPath)
	if err != nil {
		return "", err
	}
	defer db.Close()

	row := db.QueryRow(
		"select encrypted_value from cookies where host_key=? and name=?",
		".instagram.com",
		"sessionid",
	)
	if err != nil {
		return "", err
	}

	row.Scan(&encryptedSessionID)

	sessionID, err := decryptCookie(encryptedSessionID)
	if err != nil {
		return "", err
	}

	return sessionID, nil
}

func decryptCookie(value string) (string, error) {
	var (
		iterations     = 1
		length         = 16
		iv             = []byte("                ")
		decryptedValue = make([]byte, 1024)
		result         string
	)

	key := pbkdf2.Key(
		[]byte("peanuts"),
		[]byte("saltysalt"),
		iterations, length, sha1.New,
	)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	//drop "v10" prefix
	value = value[3:]

	decrypter := cipher.NewCBCDecrypter(block, iv)

	decrypter.CryptBlocks(decryptedValue, []byte(value))

	result = strings.Trim(string(decryptedValue), "\x00")
	result = strings.Trim(result, "\x10")
	result = strings.Trim(result, "\x01")

	return result, nil
}
