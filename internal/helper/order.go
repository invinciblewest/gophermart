package helper

func IsValidOrderNumber(number string) bool {
	var sum int
	double := false

	for i := len(number) - 1; i >= 0; i-- {
		r := number[i]

		if r < '0' || r > '9' {
			return false
		}

		digit := int(r - '0')

		if double {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
		double = !double
	}

	return sum%10 == 0
}
