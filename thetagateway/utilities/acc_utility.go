package utilities

import "crypto/rand"

func GenerateAccountNumber() (acc string, err error) {

	p, err := rand.Prime(rand.Reader, 64)
	if err != nil {

		return "", err
	}

	return p.String(), nil
}

func GenerateTempPin() (*int, error) {
	p, err := rand.Prime(rand.Reader, 6)
	if err != nil {

		return nil, err
	}
	value := int(p.Int64())
	return &value, nil
}
