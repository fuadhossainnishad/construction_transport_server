package notification

import (
	"encoding/json"
	"log"
	"net/smtp"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

type EmailEvent struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
}

func StartEmailConsumer(conn *amqp.Connection) {
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to open channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare("user.registered", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("failed to declare queue: %v", err)
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("failed to consume: %v", err)
	}

	for d := range msgs {
		var event EmailEvent
		if err := json.Unmarshal(d.Body, &event); err != nil {
			log.Printf("invalid email event: %v", err)
			continue
		}
		go sendEmail(event.Email, event.OTP)
	}
}

func sendEmail(to, otp string) {
	from := os.Getenv("SMTP_FROM")
	password := os.Getenv("SMTP_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	msg := []byte("To: " + to + "\r\n" +
		"Subject: Your OTP Code\r\n" +
		"\r\n" +
		"Your OTP is: " + otp + "\r\n")

	auth := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
	if err != nil {
		log.Printf("failed to send email to %s: %v", to, err)
	} else {
		log.Printf("OTP email sent to %s", to)
	}
}
