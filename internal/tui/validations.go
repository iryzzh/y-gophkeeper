package tui

import (
	"errors"
	"fmt"
	"strconv"
)

func isValidLuhn(value interface{}) error {
	number, err := strconv.Atoi(value.(string))
	if err != nil {
		return err
	}
	if (number%10+checksumLuhn(number/10))%10 != 0 {
		return errors.New("invalid order number")
	}

	return nil
}

func isValidMonthYear(value interface{}) error {
	s, _ := value.(string)
	if len(s) != 2 {
		return fmt.Errorf("invalid month length: %v", len(s))
	}
	_, err := strconv.Atoi(value.(string))
	if err != nil {
		return fmt.Errorf("invalid format: %v", value)
	}

	return nil
}

func isValidCCV(value interface{}) error {
	s, _ := value.(string)
	if len(s) != 3 {
		return fmt.Errorf("invalid month length: %v", len(s))
	}
	_, err := strconv.Atoi(value.(string))
	if err != nil {
		return fmt.Errorf("invalid format: %v", value)
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
