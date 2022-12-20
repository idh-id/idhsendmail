package idhsendmail

//	Notes:
//	Development of the gomail plugin so that users can send lots of files.
//	Please check the usage example in the "example" folder. A.. (QD)

import (
	"errors"
	"log"
	"reflect"

	"gopkg.in/gomail.v2"
)

type (
	AuthData struct { //Struct authentication for sender mail
		CONFIG_SMTP_HOST     string //host web server mail
		CONFIG_SMTP_PORT     int    //port web server mail
		CONFIG_SENDER_NAME   string //sender name <email>
		CONFIG_AUTH_EMAIL    string //sender email
		CONFIG_AUTH_PASSWORD string //password sender email
	}
	mailerData struct {
		From        string
		To          []string
		CC          string
		BCC         string
		Subject     string
		BodyMessage string
		Files       []string
	}
)

func idhMail() *gomail.Message {
	return gomail.NewMessage()
}

// The primary function for sending email data
// If you want to upload a file or multiple files,
// please use the "files" parameter with the string array data type
// and input file path on the parameter.

// Desc :
// files	= to send file (array string)
// cc		= to cc on another email (string)
// bcc		= to bcc on another email (string)
// message = to write a body message in the email (text/ HTML) (string)
func IDHSend(auth AuthData, Subject string, To []string, others ...map[string]interface{}) error {
	// Sender email authentication validation
	errAuth := validateAuth(auth)
	if errAuth != nil {
		return errAuth
	}

	dialer := gomail.NewDialer(
		auth.CONFIG_SMTP_HOST,
		auth.CONFIG_SMTP_PORT,
		auth.CONFIG_AUTH_EMAIL,
		auth.CONFIG_AUTH_PASSWORD,
	)

	var othersMap map[string]interface{}
	if others != nil && len(others) > 0 {
		othersMap = others[0]
	}

	// Parsing data on struct
	setupdata := parsData(auth.CONFIG_SENDER_NAME, Subject, To, othersMap)

	// Parsing data struct to gomail
	message := sendMessage(setupdata)

	// Process sending data
	eSend := sendDialer(dialer, message)
	if eSend != nil {
		return eSend
	}

	// Reset data on struct
	resetStruct(&auth)
	resetStruct(&setupdata)

	return nil
}

// Process parsing email data to struct mailerData
func parsData(Sender, Subject string, To []string, others ...map[string]interface{}) mailerData {
	var othersMap map[string]interface{}
	if others != nil && len(others) > 0 {
		othersMap = others[0]
	}
	mailerdata := new(mailerData)

	mailerdata.From = Sender
	mailerdata.To = To

	if othersMap["cc"] != nil && othersMap["cc"].(string) != "" {
		mailerdata.CC = othersMap["cc"].(string)
	} else {
		mailerdata.CC = ""
	}

	if othersMap["bcc"] != nil && othersMap["bcc"].(string) != "" {
		mailerdata.BCC = othersMap["bcc"].(string)
	} else {
		mailerdata.BCC = ""
	}

	if othersMap["message"] != nil && othersMap["message"].(string) != "" {
		mailerdata.BodyMessage = othersMap["message"].(string)
	} else {
		mailerdata.BodyMessage = ""
	}

	if Subject != "" {
		mailerdata.Subject = Subject
	} else {
		mailerdata.Subject = ""
	}

	var files []string
	if othersMap["files"] != nil && len(othersMap["files"].([]string)) > 0 {
		mailerdata.Files = othersMap["files"].([]string)
	} else {
		mailerdata.Files = files
	}

	return *mailerdata
}

// Process sending email
func sendDialer(d *gomail.Dialer, m *gomail.Message) error {
	err := d.DialAndSend(m)
	if err != nil {
		log.Fatal(err.Error())
	}

	return nil
}

func sendMessage(ed mailerData) *gomail.Message {
	idhMail := idhMail()
	idhMail.SetHeader("From", ed.From)
	idhMail.SetHeader("To", ed.To...)
	if ed.CC != "" {
		idhMail.SetAddressHeader("Cc", ed.CC, "")
	}
	if ed.BCC != "" {
		idhMail.SetAddressHeader("BCC", ed.BCC, "")
	}
	if ed.Subject != "" {
		idhMail.SetHeader("Subject", ed.Subject)
	}
	if ed.BodyMessage != "" {
		idhMail.SetBody("text/html", ed.BodyMessage)
	}
	if len(ed.Files) > 0 {
		for i, _ := range ed.Files {
			idhMail.Attach(ed.Files[i])
		}
		idhMail.SetBody("text/html", ed.BodyMessage)
	}

	return idhMail
}

func validateAuth(auth AuthData) error {
	if auth.CONFIG_SMTP_HOST == "" {
		return errors.New("Missing email host")
	}
	if auth.CONFIG_SMTP_PORT < 1 {
		return errors.New("Missing port")
	}
	if auth.CONFIG_AUTH_EMAIL == "" {
		return errors.New("Missing email sender")
	}
	if auth.CONFIG_AUTH_PASSWORD == "" {
		return errors.New("Missing password sender")
	}

	return nil
}

func resetStruct(v interface{}) {
	p := reflect.ValueOf(v).Elem()
	p.Set(reflect.Zero(p.Type()))
}
