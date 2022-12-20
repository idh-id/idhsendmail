package main

import (
	"fmt"
	"kotretan/idhsendmail"
	"log"
)

const CONFIG_SMTP_HOST = "mail.testing.id"
const CONFIG_SMTP_PORT = 12345
const CONFIG_SENDER_NAME = "Test <test@test.id>"
const CONFIG_AUTH_EMAIL = "test@test.id"
const CONFIG_AUTH_PASSWORD = "1testingPWD23"

func main() {

	auth := new(idhsendmail.AuthData)
	auth.CONFIG_SMTP_HOST = CONFIG_SMTP_HOST
	auth.CONFIG_SMTP_PORT = CONFIG_SMTP_PORT
	auth.CONFIG_SENDER_NAME = CONFIG_SENDER_NAME
	auth.CONFIG_AUTH_EMAIL = CONFIG_AUTH_EMAIL
	auth.CONFIG_AUTH_PASSWORD = CONFIG_AUTH_PASSWORD

	errSend := idhsendmail.IDHSend(*auth, "Testing plugin kirim email", []string{"test.to1@gmail.com", "test.to2@gmail.com"}, map[string]interface{}{
		"files": []string{"output.pdf", "446645171_082210234858.jpg"},
		"cc":    "test.cc@gmail.com",
	})
	if errSend != nil {
		fmt.Println("errSend: ", errSend)
		return
	}

	log.Println("Email Terkirim bos!")
}
