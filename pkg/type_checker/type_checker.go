package typechecker

import "strconv"

// Helper function to check if the input is a valid number
func IsNumber(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

// Helper function to check if the input is a boolean
func IsBool(s string) bool {
	_, err := strconv.ParseBool(s)
	return err == nil
}

// Helper function to check if the input is a string
func IsString(s string) bool {
	// Strings are already handled as they are enclosed in double quotes
	// So we check if they aren't numeric or boolean
	return !(IsNumber(s) || IsBool(s))
}

// Helper function to parse numbers (either int or float)
func ParseNumber(s string) interface{} {
	if i, err := strconv.Atoi(s); err == nil {
		return i
	} else if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	}
	return s // return as string if not number
}

// Helper function to parse booleans
func ParseBool(s string) bool {
	b, _ := strconv.ParseBool(s)
	return b
}
