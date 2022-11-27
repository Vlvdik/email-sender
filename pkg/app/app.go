package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mail-Sender/config"
	"mail-Sender/pkg/server"
	"net/smtp"
	"time"
)

type Jobs []struct {
	Date     time.Time     `json:"job_date"`
	Duration time.Duration `json:"duration"`
}

type Job struct {
	Date     time.Time     `json:"job_date"`
	Duration time.Duration `json:"duration"`
}

func writeJob(job Job) {
	dataIn, err := ioutil.ReadFile("pkg/app/config/jobs.json")
	if err != nil {
		log.Fatal(err)
	}

	var jobs Jobs
	err = json.Unmarshal(dataIn, &jobs)
	if err != nil {
		log.Fatal(err)
	}

	jobs = append(jobs, job)

	dataOut, err := json.MarshalIndent(&jobs, "", " ")
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile("pkg/app/config/jobs.json", dataOut, 0)
	if err != nil {
		log.Fatal(err)
	}
}

func deleteJob(job Job) {
	dataIn, err := ioutil.ReadFile("pkg/app/config/jobs.json")
	if err != nil {
		log.Fatal(err)
	}

	var jobs Jobs
	err = json.Unmarshal(dataIn, &jobs)
	if err != nil {
		log.Fatal(err)
	}

	if len(jobs) > 1 {
		for idx, value := range jobs {
			if value.Date.Format("02-Jan-2006 15:04:05") == job.Date.Format("02-Jan-2006 15:04:05") {
				jobs[idx] = jobs[len(jobs)-1]
				jobs = jobs[:len(jobs)-1]
			}
		}
	} else {
		jobs = Jobs{}
	}

	dataOut, err := json.MarshalIndent(&jobs, "", " ")
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile("pkg/app/config/jobs.json", dataOut, 0)
	if err != nil {
		log.Fatal(err)
	}
}

func FinishJobs(sd config.SendData) {
	dataIn, err := ioutil.ReadFile("pkg/app/config/jobs.json")
	if err != nil {
		log.Fatal(err)
	}

	var jobs Jobs
	err = json.Unmarshal(dataIn, &jobs)
	if err != nil {
		log.Fatal(err)
	}

	if len(jobs) > 0 {
		stamp := time.Now()
		for _, value := range jobs {
			if value.Date.Add(value.Duration).Before(stamp) {
				SendMails(sd)
				deleteJob(value)
			} else {
				currentDuration := value.Date.Add(value.Duration).Sub(stamp)

				SendMailsWithDuration(sd, currentDuration)
				deleteJob(value)

				fmt.Printf("Выполню чуть позже, через %v секунд", currentDuration)
			}
		}
	}
}

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
	job := Job{Date: stamp, Duration: duration}

	writeJob(job)
	time.Sleep(job.Duration)

	for _, value := range sd.To {
		sd.Data.SetSubjectPersonInfo(value.PersonalInfo.Name, value.PersonalInfo.LastName)
		sd.Data.SetBodyPersonInfo(value.PersonalInfo, value.Email[0])

		err := smtp.SendMail(address, auth, sd.From.Email, value.Email, sd.Data.Message)
		if err != nil {
			log.Fatal(err)
		}
	}
	deleteJob(job)
	log.Printf("Отложенная рассылка успешно выполнена!\n Дата и время задачи: %v\n", stamp.Format("02-Jan-2006 15:04:05"))
}

func Close(es *server.EmailServer, ctx context.Context) {
	log.Println("Приложение закрывается")

	err := es.Server.Shutdown(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
