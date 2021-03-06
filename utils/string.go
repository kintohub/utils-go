package utils

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"strings"
)

/**
 * Idea from here: https://stackoverflow.com/questions/2253693/shortening-rehashing-uuids
 * TL;DR:
 *  This will preserve the distribution characteristics.
 * Like any solution that shortens the output, it will increase the possibility of collision due to the birthday paradox
 */
func ShortenUUID16(uuid string) (string, error) {
	u := strings.Replace(uuid, "-", "", -1)
	if len(u) != 32 {
		return "", fmt.Errorf("invalid uuid")
	}

	bytes, err := hex.DecodeString(u)
	if err != nil {
		return "", err
	}

	// uuid: 128bit = 32 hex character = 16bytes
	// final output should by 8 bytes, 16 characters
	const n = 8
	returnBytes := make([]byte, n)

	for i := 0; i < n; i++ {
		returnBytes[i] = bytes[i] ^ bytes[i+n]
	}

	return hex.EncodeToString(returnBytes), nil
}

/**
This method will further shorten the uuid to 8 character i.e. 4 bytes
Take cautions of collision before using this function
*/
func ShortenUUID8(uuid string) (string, error) {
	u := strings.Replace(uuid, "-", "", -1)
	if len(u) != 32 {
		return "", fmt.Errorf("invalid uuid")
	}

	bytes, err := hex.DecodeString(u)
	if err != nil {
		return "", err
	}

	// uuid: 128bit = 32 hex character = 16bytes
	// final output should by 4 bytes, 8 characters
	const n = 4
	returnBytes := make([]byte, n)

	for i := 0; i < n; i++ {
		returnBytes[i] = bytes[i] ^ bytes[i+n] ^ bytes[i+n*2] ^ bytes[i+n*3]
	}

	return hex.EncodeToString(returnBytes), nil
}

/**
This will shorten the input string to len [outputLen], with bitwise xor the bytes
The length of the input string/outputLen must be both even number,
And outputLen must be a positive and even divisor of the input length, and cannot be 1 or input length
i.e. if input length is 24, outputLen can be 12,8,6,4,2. but cannot be 24, 3, 1

This function will try to ensure the distribution of the randomness from the given input.
But again, like the function ShortenUUID16 above, this will increase the possibility on collision, especially when the outputLen is too short
*/
func ShortenHexString(input string, outputLen int) (string, error) {
	if len(input)%2 != 0 || outputLen%2 != 0 {
		return "", fmt.Errorf("input/output must have even number characters")
	}

	// outputLen must be a valid dividend of inputLen
	if len(input)%outputLen != 0 || outputLen == len(input) || outputLen == 1 {
		return "", fmt.Errorf("outputLen must be a valid divident")
	}

	bytes, err := hex.DecodeString(input)
	if err != nil {
		return "", err
	}

	n := outputLen / 2
	step := len(input) / outputLen
	returnBytes := make([]byte, n)

	for i := 0; i < n; i++ {
		returnBytes[i] = 0
		for j := 0; j < step; j++ {
			returnBytes[i] ^= bytes[i+j*n]
		}
	}

	return hex.EncodeToString(returnBytes), nil
}

// Random string function from stackoverflow
// not the best one but it is simple
// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
