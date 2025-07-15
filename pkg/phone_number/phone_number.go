package phone_number

func IsValidPhoneNumber(phoneNumber string) bool {

	// TODO refactor to REGX

	if len(phoneNumber) != 11 {
		return false
	}

	if phoneNumber[0:2] != "09" {
		return false
	}

	if phoneNumber == "" {
		return false
	}

	return true
}
