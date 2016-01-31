package main

import (
	"crypto/rand"
	"encoding/base64"
	"flag"
	"fmt"
)

func main() {
	serverMode := flag.Bool("s", false, "server mode")
	name := flag.String("name", "", "name")
	delete := flag.String("delete", "", "key to delete")
	pub := flag.String("pub", "", "public identity")
	secret := flag.String("secret", "", "secret key")
	app := flag.String("app", "", "application name")
	port := flag.String("p", "4242", "server port")
	host := flag.String("host", "127.0.0.1", "server addr")
	db := flag.String("db", "database.db", "database file")
	flag.Parse()

	dal, err := newDAL(*db)
	if err != nil {
		fmt.Println(err)
	}

	if *serverMode {
		runAPI(dal, *host, *port)
	} else {
		if *app != "" {
			randomkey := make([]byte, 256)
			_, err := rand.Read(randomkey)
			if err != nil {
				fmt.Println("error getting random data:", err)
			} else {
				app, err := dal.CreateApp(&App{Name: *app, Key: randomkey})
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println("app created, id:", app.Id, "key:", base64.StdEncoding.EncodeToString(app.Key))
				}
			}
		} else {
			if *delete != "" {

				err := dal.DeleteKey(&Key{Name: *delete})
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println("key deleted: OK")
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
}
