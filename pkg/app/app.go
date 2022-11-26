package app

import (
	"context"
	"log"
	"mail-Sender/config"
	"mail-Sender/server"
	"net/smtp"
	"time"
)

func Init() (config.SendData, server.EmailServer) {
	var s config.Sender
	s.NewSender()

	var r config.Receivers
	r.NewReceivers()

	var d config.Data
	d.NewData()

	var sd config.SendData
	sd.NewSendData(s, r, d)

	var es server.EmailServer
	es.SetEmailServerData(sd)

	return sd, es
}

func SendMails(sd config.SendData) {
	auth := sd.From.Auth()
	address := sd.From.Host + ":" + sd.From.Port

	for _, value := range sd.To {
		sd.Data.SetSubjectPersonInfo(value.PersonalInfo.Name, value.PersonalInfo.LastName)
		sd.Data.SetBodyPersonInfo(value.PersonalInfo, value.Email[0])

		err := smtp.SendMail(address, auth, sd.From.Email, value.Email, sd.Data.Message)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Println("Рассылка успешно выполнена!")
}

func SendMailsWithDuration(sd config.SendData, duration time.Duration) {
	auth := sd.From.Auth()
	address := sd.From.Host + ":" + sd.From.Port

	stamp := time.Now()
	time.Sleep(duration * time.Second)

	for _, value := range sd.To {
		sd.Data.SetSubjectPersonInfo(value.PersonalInfo.Name, value.PersonalInfo.LastName)
		sd.Data.SetBodyPersonInfo(value.PersonalInfo, value.Email[0])

		err := smtp.SendMail(address, auth, sd.From.Email, value.Email, sd.Data.Message)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Printf("Отложенная рассылка успешно выполнена!\n Дата и время задачи: %v\n", stamp.Format("02-Jan-2006 15:04:05"))
}

func Close(es *server.EmailServer, ctx context.Context) {
	log.Println("Приложение закрывается")

	err := es.Server.Shutdown(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
