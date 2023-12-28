package util

import (
	"fmt"
	"math/rand"
)

func GenerateOTP() string {

	arr := make([]int, 4)
	otpString := ""

	for range arr {
		otpString += fmt.Sprintf("%v", rand.Intn(10))
	}

	return otpString
}

func SendEmailOTP(email string, otp string) error{
	fmt.Printf("Sent otp (%v) to %v", otp, email)
	return nil
}
