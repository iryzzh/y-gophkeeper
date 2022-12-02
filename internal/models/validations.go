//nolint:gomnd
package models

import (
	"strconv"

	"github.com/pkg/errors"
)

func IsValidLuhn(value interface{}) error {
	number, err := strconv.Atoi(value.(string))
	if err != nil {
		return err
	}
	if (number%10+checksumLuhn(number/10))%10 != 0 {
		return errors.New("invalid card number")
	}

	return nil
}

func checksumLuhn(number int) int {
	var luhn int

	for i := 0; number > 0; i++ {
		cur := number % 10

		if i%2 == 0 { // even
			cur *= 2
			if cur > 9 {
				cur = cur%10 + cur/10
			}
		}

		luhn += cur
		number /= 10
	}

	return luhn % 10
}
