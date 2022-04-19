package textutil

// IsLetter reports whether the rune is a letter
func IsLetter(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

// isDigit reports whether the rune is a decimal digit.
func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

// IsLetterDigitSymbol reports whether the rune is a decimal digit or any of passed symbols.
func IsLetterDigitSymbol(r rune, symbols ...rune) bool {
	if IsLetter(r) || isDigit(r) {
		return true
	}
	for _, symbol := range symbols {
		if r == symbol {
			return true
		}
	}
	return false
}
