package utils

import (
	"strconv"
	"strings"
)

func CleanCPF(cpf string) string {
	cpfNew := strings.ReplaceAll(cpf, ".", "")
	cpfNew = strings.ReplaceAll(cpfNew, "-", "")
	cpfNew = strings.ReplaceAll(cpfNew, " ", "")
	return cpfNew
}

// ValidateCPF verifica se um CPF é válido.
func ValidateCPF(cpf string) bool {

	if len(cpf) != 11 || InvalidCPF(cpf) {
		return false
	}

	digits := make([]int, 11)
	for i, c := range cpf {
		d, err := strconv.Atoi(string(c))
		if err != nil {
			return false
		}
		digits[i] = d
	}

	dv1 := CalculateCheckDigit(digits[:9])
	dv2 := CalculateCheckDigit(digits[:10])

	return dv1 == digits[9] && dv2 == digits[10]
}

// InvalidCPF verifica se o CPF está na lista de CPFs conhecidos por serem inválidos.
func InvalidCPF(cpf string) bool {
	invalids := []string{
		"00000000000", "11111111111", "22222222222", "33333333333",
		"44444444444", "55555555555", "66666666666", "77777777777",
		"88888888888", "99999999999",
	}
	for _, inv := range invalids {
		if cpf == inv {
			return true
		}
	}
	return false
}

func CalculateCheckDigit(nums []int) int {
	total := 0
	weight := len(nums) + 1
	for _, num := range nums {
		total += num * weight
		weight--
	}
	leftover := total % 11
	if leftover < 2 {
		return 0
	}
	return 11 - leftover
}

// Validar email
func ValidateEmail(email string) bool {
	if !strings.Contains(email, "@") {
		return false
	}
	if !strings.Contains(email, ".") {
		return false
	}
	return true
}
