package utilities

func GetDays(num int) string {
	switch {
	case num == 1:
		return "день"
	case num >= 2 && num <= 4:
		return "дня"
	default:
		return "дней"
	}
}
