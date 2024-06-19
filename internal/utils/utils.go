package utils

import (
	"errors"
)

var (
	ErrOrderNumAttachedToAnotherUser = errors.New("sorry, but order number is already attached to another user")
	ErrOrderNumIsAlreadyRegistered   = errors.New("sorry, but order number is already registered")
	ErrInvalidOrderNum               = errors.New("invalid order number")
	ErrInsufficientFunds             = errors.New("insufficient funds")
)

func VerifyLuhn(code string) bool {
	if len(code) < 2 {
		return false
	}
	i, err := GenerateLuhn(code[:len(code)-1])

	return err == nil && i == int(code[len(code)-1]-'0')
}

func GenerateLuhn(seed string) (int, error) {
	if seed == "" {
		return 0, errors.New("invalid Argument")
	}

	sum, parity := 0, (len(seed)+1)%2
	for i, n := range seed {
		if isNotNumber(n) {
			return 0, errors.New("invalid Argument")
		}
		d := int(n - '0')
		if i%2 == parity {
			d *= 2
			if d > 9 {
				d -= 9
			}
		}
		sum += d
	}

	return sum * 9 % 10, nil
}

func isNotNumber(n rune) bool {
	return n < '0' || '9' < n
}
