package utils

import (
	"fmt"
	"strconv"
)

func GenerateBarcode(ticketNumber int64) (string, error) {
	// формируем 12-значную основу: "22001" + ticketNumber с ведущими нулями до 7 символов
	base := fmt.Sprintf("22001%07d", ticketNumber)

	if len(base) != 12 {
		return "", fmt.Errorf("invalid base length: %s", base)
	}

	var sumOdd, sumEven int
	for i, r := range base {
		digit, err := strconv.Atoi(string(r))
		if err != nil {
			return "", err
		}
		if i%2 == 0 {
			sumOdd += digit // чётные индексы
		} else {
			sumEven += digit // нечётные индексы
		}
	}

	checkDigit := 10 - (sumEven*3+sumOdd)%10
	if checkDigit == 10 {
		checkDigit = 0
	}

	return base + strconv.Itoa(checkDigit), nil
}
