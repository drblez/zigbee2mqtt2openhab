package utils

func BoolToString(b bool, trueStr string, falseStt string) string {
	if b {
		return trueStr
	} else {
		return falseStt
	}
}
