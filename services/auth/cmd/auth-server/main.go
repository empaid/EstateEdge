package main

import (
	"github.com/empaid/estateedge/pkg/config"
)

func main() {

	cfg := cfg{
		addr: config.GetString("ADDR", ":3000"),
	}

	app := &application{
		config: cfg,
	}

	app.Run()

}
