package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"mail-Sender/pkg/app"
	"time"
)

func main() {
	sd, es := app.Init()
	app.FinishJobs(sd)
	ctx := context.Background()
	go func() {
		err := es.Start()
		if err != nil {
			log.Fatal(err)
		}
	}()
	defer app.Close(&es, ctx)

	fmt.Println("Привет! Чтобы отослать рассылку напиши 1, чтобы отослать рассылку через промежуток времени - 2, чтобы выйти - 0")
	for {
		var input string
		fmt.Println("Выберите команду: ")
		_, err := fmt.Scan(&input)
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}

		if input == "1" {
			app.SendMails(sd)

		} else if input == "2" {
			var duration time.Duration
			fmt.Println("Укажите задержку в секундах: ")
			_, err = fmt.Scan(&duration)
			if err != nil && err != io.EOF {
				log.Fatal(err)
			}

			go app.SendMailsWithDuration(sd, duration*time.Second)
			fmt.Printf("Задача будет выполнена через %v секунд\n", duration)
		} else if input == "0" {
			break
		} else {
			fmt.Println("Неизвестная команда")
		}
	}
}
