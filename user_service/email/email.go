package email

import (
	"fmt"
	"log"
	"net/smtp"
)

type EmailService struct {
	from     string
	password string
	smtpHost string
	smtpPort string
}

func NewEmailService() *EmailService {
	return &EmailService{
		from:     "dungcongnghiep4@gmail.com", // Email của bạn
		password: "vclrrtxkmkgcccit",
		smtpHost: "smtp.gmail.com",
		smtpPort: "587",
	}
}

// SendOrderCancelledEmail gửi email thông báo hủy đơn hàng
func (e *EmailService) SendOrderCancelledEmail(toEmail string, orderID uint32) error {
	subject := "Đơn hàng đã được hủy"
	body := fmt.Sprintf(`
		Xin chào!
		
		Đơn hàng #%d của bạn đã được hủy thành công.
		
		Chúng tôi sẽ hoàn tiền trong vòng 3-5 ngày làm việc.
		
		Cảm ơn bạn đã sử dụng dịch vụ của chúng tôi!
		
		Trân trọng,
		E-commerce Team
	`, orderID)

	return e.sendEmail(toEmail, subject, body)
}

// sendEmail gửi email
func (e *EmailService) sendEmail(to, subject, body string) error {
	// Tạo message
	message := fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", to, subject, body)

	// Auth
	auth := smtp.PlainAuth("", e.from, e.password, e.smtpHost)

	// Gửi email
	err := smtp.SendMail(e.smtpHost+":"+e.smtpPort, auth, e.from, []string{to}, []byte(message))
	if err != nil {
		log.Printf("Failed to send email: %v", err)
		return err
	}

	log.Printf("Email sent successfully to %s", to)
	return nil
}
