package main

import (
	"mail-Sender/config"
)

func main() {
	var s config.Sender
	s.NewSender()

	var r config.Receivers
	r.NewReceivers()

	var d config.Data
	d.NewData()

	var sd config.SendData
	sd.NewSendData(s, r, d)

	config.SendMails(sd)
}
