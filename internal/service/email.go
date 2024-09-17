package service

import (
	"fmt"
	"reminders-api/internal/config"
	"reminders-api/internal/models"
	"time"

	"gopkg.in/gomail.v2"
)

type EmailService struct {
	config *config.Config
}

func NewEmailService(cfg *config.Config) *EmailService {
	return &EmailService{
		config: cfg,
	}
}

func (s *EmailService) SendReminderEmail(event models.Event) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.config.EmailFrom)
	m.SetHeader("To", s.config.EmailTo)

	subject := fmt.Sprintf("Reminder: %s", event.Name)
	body := s.generateEmailBody(event)

	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer("smtp.gmail.com", 587, s.config.EmailFrom, s.config.EmailPassword)

	return d.DialAndSend(m)
}

func (s *EmailService) generateEmailBody(event models.Event) string {
	now := time.Now()
	if event.Recurring {
		return fmt.Sprintf("This is a reminder for your recurring event: %s", event.Name)
	}

	daysUntil := int(event.Date.Sub(now).Hours() / 24)
	if daysUntil > 0 {
		return fmt.Sprintf("Reminder: %d days until %s", daysUntil, event.Name)
	} else if daysUntil == 0 {
		return fmt.Sprintf("Reminder: %s is today!", event.Name)
	} else {
		return fmt.Sprintf("Reminder: %s was %d days ago", event.Name, -daysUntil)
	}
}

func (s *EmailService) SendDailyReminders(events []models.Event) {
	now := time.Now()
	for _, event := range events {
		if event.Recurring && event.Date.Day() == now.Day() {
			err := s.SendReminderEmail(event)
			if err != nil {
				fmt.Printf("Error sending reminder for %s: %v\n", event.Name, err)
			}
		} else if !event.Recurring {
			daysUntil := int(event.Date.Sub(now).Hours() / 24)
			if daysUntil <= 7 && daysUntil >= 0 {
				err := s.SendReminderEmail(event)
				if err != nil {
					fmt.Printf("Error sending reminder for %s: %v\n", event.Name, err)
				}
			}
		}
	}
}
