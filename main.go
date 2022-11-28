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
	go app.FinishJobs(sd)
	ctx := context.Background()
	go func() {
		err := es.Start()
		if err != nil {
			log.Fatal(err)
		}
	}()
	defer app.Close(ctx, &es)

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
			app.SendMails(sd)
		case "2":
			var duration time.Duration
			fmt.Println("\nSpecify the delay in seconds: ")
			_, err = fmt.Scan(&duration)
			if err != nil && err != io.EOF {
				log.Fatal(err)
			}

			go app.SendMailsWithDuration(sd, duration*time.Second)
			fmt.Printf("\nThe task will be completed in %v seconds\n", duration)
		case "3":
			app.GetJobs()
		case "0":
			return
		default:
			fmt.Println("\nUnknown command")
		}
	}
}
