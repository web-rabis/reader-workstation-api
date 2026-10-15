package utils

import "testing"

// ean13Valid — независимая проверка контрольной цифры EAN-13.
func ean13Valid(code string) bool {
	if len(code) != 13 {
		return false
	}
	sum := 0
	for i := 0; i < 12; i++ {
		d := int(code[i] - '0')
		if i%2 == 1 {
			d *= 3
		}
		sum += d
	}
	return (10-sum%10)%10 == int(code[12]-'0')
}

func TestGenerateBarcode_KnownValues(t *testing.T) {
	cases := map[int64]string{
		0: "2200100000001",
		1: "2200100000018",
	}
	for ticket, want := range cases {
		got, err := GenerateBarcode(ticket)
		if err != nil {
			t.Fatalf("ticket %d: %v", ticket, err)
		}
		if got != want {
			t.Errorf("ticket %d: got %s, want %s", ticket, got, want)
		}
	}
}

func TestGenerateBarcode_ValidEAN13(t *testing.T) {
	for _, ticket := range []int64{1, 9, 10, 99, 12345, 130000, 999999, 9999999} {
		got, err := GenerateBarcode(ticket)
		if err != nil {
			t.Fatalf("ticket %d: %v", ticket, err)
		}
		if !ean13Valid(got) {
			t.Errorf("ticket %d: %s не проходит проверку EAN-13", ticket, got)
		}
	}
}

func TestGenerateBarcode_TooLongTicket(t *testing.T) {
	if _, err := GenerateBarcode(12345678); err == nil {
		t.Fatal("ожидалась ошибка для номера длиннее 7 цифр")
	}
}
