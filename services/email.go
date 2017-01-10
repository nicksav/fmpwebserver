package services

import "gopkg.in/gomail.v2"
import "fmpwebserver/services/models"



func SendConfirmEmail(u *models.User,code string)(bool, error){

    m := gomail.NewMessage()
    m.SetHeader("From", "madrobots@mail.ru")
    m.SetHeader("To", u.UserName)
    //m.SetAddressHeader("Cc", "dan@example.com", "Dan")
    m.SetHeader("Subject", "Hello!")
    m.SetBody("text/html", "Hello <b>Bob</b> and <i>Cora. This is your code: "+code+"</i>!")

    d := gomail.NewDialer("smtp.mail.ru", 465, "madrobots@mail.ru", "hast1ng$")

    // Send the email to Bob, Cora and Dan.
    if err := d.DialAndSend(m); err != nil {
        return false, err
    }

    return true, nil
}

func SendConfirmEmailKiosk(k *models.Kiosk,code string)(bool, error){

    m := gomail.NewMessage()
    m.SetHeader("From", "madrobots@mail.ru")
    m.SetHeader("To", k.UserName)
    //m.SetAddressHeader("Cc", "dan@example.com", "Dan")
    m.SetHeader("Subject", "New Kiosk has been added!")
    m.SetBody("text/html", "New Kiosk device has been added. Please use this code: <b>"+code+"</b>")

    d := gomail.NewDialer("smtp.mail.ru", 465, "madrobots@mail.ru", "hast1ng$")

    // Send the email to Bob, Cora and Dan.
    if err := d.DialAndSend(m); err != nil {
        return false, err
    }

    return true, nil
}
