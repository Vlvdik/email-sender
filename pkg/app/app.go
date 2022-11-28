package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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

func finishJobs(sd config.SendData) {
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
		log.Println("\nExecuting unfinished jobs")
		stamp := time.Now()
		for _, value := range jobs {
			if value.Date.Add(value.Duration).Before(stamp) {
				sendMails(sd)
				deleteJob(value)
			} else {
				currentDuration := value.Date.Add(value.Duration).Sub(stamp)
				fmt.Printf("Execute in %v seconds", currentDuration)
				time.Sleep(currentDuration)

				deleteJob(value)
				sendMails(sd)
			}
		}
	}
}

func getJobs() {
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
		for _, value := range jobs {
			fmt.Printf("\nDate of creation: %v\nDelay before sending: %v",
				value.Date.Format("02-Jan-2006 15:04:05"),
				value.Duration)
		}
	} else {
		fmt.Print("\nNo postponed mailings!")
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

func sendMails(sd config.SendData) {
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

	log.Println("\nLetter successfully sent!")
}

func sendMailsWithDuration(sd config.SendData, duration time.Duration) {
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
	log.Printf("\nDeffered mailing successfully completed!\n Date and time of task: %v\n", stamp.Format("02-Jan-2006 15:04:05"))
}

func Start(sd *config.SendData, es *server.EmailServer) {
	go finishJobs(*sd)
	go func() {
		err := es.Start()
		if err != nil {
			log.Fatal(err)
		}
	}()

	fmt.Println("\nHi! To send a newsletter write 1, to send a deferred newsletter write 2, to see a list of deferred newsletters write 3, to exit write 0.")
	for {
		var input string
		fmt.Println("\nSelect the command: ")
		_, err := fmt.Scan(&input)
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}

		switch input {
		case "1":
			sendMails(*sd)
		case "2":
			var duration time.Duration
			fmt.Println("\nSpecify the delay in seconds: ")
			_, err = fmt.Scan(&duration)
			if err != nil && err != io.EOF {
				log.Fatal(err)
			}

			go sendMailsWithDuration(*sd, duration*time.Second)
			fmt.Printf("\nThe task will be completed in %v seconds\n", duration)
		case "3":
			getJobs()
		case "0":
			return
		default:
			fmt.Println("\nUnknown command")
		}
	}
}

func Close(ctx context.Context, es *server.EmailServer) {
	log.Println("\nThe app is closing")

	err := es.Server.Shutdown(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
