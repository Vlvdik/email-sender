package config

import (
	"bytes"
	"encoding/json"
	"github.com/BurntSushi/toml"
	"html/template"
	"io/ioutil"
	"log"
	"net/smtp"
)

const mime = "MIME-version: 1.0;\nContent-type: text/html; charset=\"UTF-8\";\n\n"

type Person struct {
	Name     string `json:"Name"`
	LastName string `json:"lastName"`
	Birthday string `json:"birthday"`
}

type Receivers []struct {
	Email        []string `json:"email"`
	PersonalInfo Person   `json:"personalInfo"`
}

func (r *Receivers) NewReceivers() {
	data, err := ioutil.ReadFile("config/users.json")
	if err != nil {
		log.Fatal(err)
	}

	if err = json.Unmarshal(data, &r); err != nil {
		log.Fatal(err)
	}
}

type Sender struct {
	Email    string `toml:"sender_email"`
	Password string `toml:"password"`
	Host     string `toml:"host"`
	Port     string `toml:"port"`
}

func (s *Sender) NewSender() {
	_, err := toml.DecodeFile("config/config.toml", &s)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Sender) Auth() smtp.Auth {
	return smtp.PlainAuth("", s.Email, s.Password, s.Host)
}

type Data struct {
	subject string
	body    string
	message []byte
}

func (d *Data) NewData() {
	d.subject = ""
	d.body = ""
	d.message = []byte(d.subject + mime + d.body)
}

func (d *Data) GetTemplate(personInfo Person) string {
	var t *template.Template
	t, err := template.ParseFiles("templates/messageTemplate.html")
	if err != nil {
		log.Fatal(err)
	}

	body := new(bytes.Buffer)
	err = t.Execute(body, personInfo)
	if err != nil {
		log.Fatal(err)
	}

	return body.String()
}

func (d *Data) SetSubjectPersonInfo(name string, lastName string) {
	d.subject = name + " " + lastName
	d.message = []byte(d.subject + mime + d.body)
}

func (d *Data) SetBodyPersonInfo(personInfo Person) {
	d.body = d.GetTemplate(personInfo)
	d.message = []byte(d.subject + mime + d.body)
}

type SendData struct {
	From Sender
	To   Receivers
	Data Data
}

func (sd *SendData) NewSendData(from Sender, to Receivers, data Data) {
	sd.From = from
	sd.To = to
	sd.Data = data
}

func SendMails(sd SendData) {
	auth := sd.From.Auth()
	address := sd.From.Host + ":" + sd.From.Port

	for _, value := range sd.To {
		sd.Data.SetSubjectPersonInfo(value.PersonalInfo.Name, value.PersonalInfo.LastName)
		sd.Data.SetBodyPersonInfo(value.PersonalInfo)

		err := smtp.SendMail(address, auth, sd.From.Email, value.Email, sd.Data.message)
		if err != nil {
			log.Fatal(err)
		}
	}
}
