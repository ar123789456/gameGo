package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/mail"
	"net/smtp"
	"text/template"
)

var tpl *template.Template

type Form struct {
	name string
	tel  string
	mail string
}

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.html"))

}

func main() {
	http.HandleFunc("/", index)
	http.ListenAndServe(":8080", nil)
	// SendMes("ar123789456@mail.ru", "Arman", "4545445", "dasfghdj")
}

func index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		// http.Error()
		return
	}
	if r.Method == http.MethodPost {
		var form Form
		form.mail = r.FormValue("mail")
		form.name = r.FormValue("name")
		form.tel = r.FormValue("phone")
		SendMes("ar123789456@mail.ru", form.name, form.tel, form.mail)
	}
	err := tpl.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		// errorHandler(w, r, http.StatusInternalServerError)
		return
	}
}

func SendMes(tO, Name, tel, maIl string) {

	from := mail.Address{"", "ar123789456@mail.ru"}
	to := mail.Address{"", tO}
	subj := "Этот чувак заинтересовался: " + Name
	body := fmt.Sprintf("Name: %v \nNumber: %v\n mail: %v", Name, tel, maIl)

	// Setup headers
	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subj

	// Setup message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// Connect to the SMTP Server
	servername := "smtp.mail.ru:465"

	host, _, _ := net.SplitHostPort(servername)

	auth := smtp.PlainAuth("", "ar123789456@mail.ru", "hh3EAY8JuzK9XTrDVaNY", host)

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", servername, tlsconfig)
	if err != nil {
		log.Panic(err)
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		log.Panic(err)
	}

	// Auth
	if err = c.Auth(auth); err != nil {
		log.Panic(err)
	}

	// To && From
	if err = c.Mail(from.Address); err != nil {
		log.Panic(err)
	}

	if err = c.Rcpt(to.Address); err != nil {
		log.Panic(err)
	}

	// Data
	w, err := c.Data()
	if err != nil {
		log.Panic(err)
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		log.Panic(err)
	}

	err = w.Close()
	if err != nil {
		log.Panic(err)
	}

	c.Quit()
}
