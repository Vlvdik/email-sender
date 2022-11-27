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

type Receiver struct {
	Email        string
	PersonalInfo Person
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

	err = json.Unmarshal(data, &r)
	if err != nil {
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
	Subject string
	Body    string
	Message []byte
}

func (d *Data) NewData() {
	d.Subject = ""
	d.Body = ""
	d.Message = []byte(d.Subject + mime + d.Body)
}

func (d *Data) GetTemplate(personInfo Person, email string) string {
	var info = Receiver{Email: email, PersonalInfo: personInfo}

	var t *template.Template
	t = template.Must(template.ParseFiles("templates/messageTemplate.html"))

	body := new(bytes.Buffer)
	err := t.Execute(body, info)
	if err != nil {
		log.Fatal(err)
	}

	return body.String()
}

func (d *Data) SetSubjectPersonInfo(name string, lastName string) {
	d.Subject = name + " " + lastName
	d.Message = []byte(d.Subject + mime + d.Body)
}

func (d *Data) SetBodyPersonInfo(personInfo Person, email string) {
	d.Body = d.GetTemplate(personInfo, email)
	d.Message = []byte(d.Subject + mime + d.Body)
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
