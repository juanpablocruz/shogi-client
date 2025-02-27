package config

import (
	"flag"
	"os"
)

type Config struct {
	GotePlayer  string `json:"gotePlayer"`
	SentePlayer string `json:"sentePlayer"`
	Port        int    `json:"port"`
}

func Init() Config {
	var config Config

	sente := flag.String("sente", "human", "sente(black) piece input")
	gote := flag.String("gote", "cpu", "gote(white) piece input")

	port := flag.Int("p", 8080, "server port to connect to")

	flag.Parse()

	if *sente != "human" && *gote != "human" {
		flag.PrintDefaults()
		os.Exit(0)
	}

	config.SentePlayer = *sente
	config.GotePlayer = *gote
	config.Port = *port

	return config
}
