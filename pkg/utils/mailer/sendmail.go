package mailer

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendOnboardingMail(mailId string,token string)  error{
        from := os.Getenv("GMAIL_USERNAME")
        pass := os.Getenv("GMAIL_PASSWORD")
        to := mailId
    
        msg := "From: " + from + "\n" +
            "To: " + to + "\n" +
            "Subject: Verify your Identity\n\n" +
            "Here is your magic link: https://audio-library-frontend.vercel.app/login/verify?token="+token+"\n\nDo not Share this link with anyone\n\n\n\n" 
    
        err := smtp.SendMail("smtp.gmail.com:587",
            smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
            from, []string{to}, []byte(msg))
        
        fmt.Printf("Error sending mail: %v", err)
    
        return err
}