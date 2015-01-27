package main

import (
	"flag"
	"fmt"
	"time"
)

func main() {
	serverMode := flag.Bool("s", false, "server mode")
	name := flag.String("name", "", "name")
	pub := flag.String("pub", "", "public identity")
	secret := flag.String("secret", "", "secret key")
	app := flag.String("app", "", "application name")
	port := flag.String("p", "4242", "server port")
	flag.Parse()

	dal, err := newDAL()
	if err != nil {
		fmt.Println(err)
	}

	if *serverMode {
		runAPI(dal, *port)
	} else {
		if *app != "" {
			app, err := dal.CreateApp(&App{Name: *app, Key: time.Now().Format(time.RFC3339Nano)})
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("app created, id:", app.Id, "key:", app.Key)
			}
		} else {
			err := dal.CreateKey(&Key{Name: *name, Public: *pub, Secret: *secret})
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("creation of the key: OK")
			}
		}
	}
}
