package main

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

type Dal struct {
	db *sql.DB
}

type Key struct {
	Id      int
	Name    string
	Created *time.Time
	Used    *time.Time
	Counter int
	Session int
	Public  string
	Secret  string
}

type App struct {
	Id   int
	Name string
	Key  string
}

func newDAL() (*Dal, error) {
	d, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		return nil, err
	}

	ret := Dal{db: d}
	ret.init()

	return &ret, nil
}

func (d *Dal) CreateApp(app *App) (*App, error) {
	app.Key = Sign([]string{app.Name}, app.Key)

	stmt, err := d.db.Prepare(`insert into apps(name, key, created) values(?, ?, ?)`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	_, err = stmt.Exec(app.Name, app.Key, time.Now())
	if err != nil {
		return nil, err
	}

	stmt2, err := d.db.Prepare("select MAX(id) from apps LIMIT 1")
	if err != nil {
		return nil, err
	}
	defer stmt2.Close()
	err = stmt2.QueryRow().Scan(&app.Id)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (d *Dal) CreateKey(key *Key) error {
	if key.Name == "" {
		return errors.New("name need to be indicated")
	} else if key.Public == "" {
		return errors.New("pub need to be indicated")
	} else if key.Secret == "" {
		return errors.New("secret need to be indicated")
	} else {
		k, _ := d.GetKey(key.Public)
		if k != nil {
			return errors.New("public key: " + key.Public + " already exists")
		} else {
			stmt, err := d.db.Prepare(`insert into keys(name, created, counter, session, public, secret) values(?, ?, ?, ?, ?, ?)`)
			if err != nil {
				return err
			}
			defer stmt.Close()
			_, err = stmt.Exec(key.Name, time.Now(), 0, 0, key.Public, key.Secret)
			if err != nil {
				return err
			}
			return nil
		}
	}
	return nil
}

func (d *Dal) UpdateKey(key *Key) error {
	stmt, err := d.db.Prepare("update keys set counter = ?, session = ?, used = ?where public = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(key.Counter, key.Session, time.Now(), key.Public)
	if err != nil {
		return err
	}

	return nil
}

func (d *Dal) GetApp(id string) (*string, error) {
	stmt, err := d.db.Prepare("select key from apps where id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	key := ""
	err = stmt.QueryRow(id).Scan(&key)
	if err != nil {
		return nil, err
	}
	return &key, nil
}

func (d *Dal) GetKey(pub string) (*Key, error) {
	stmt, err := d.db.Prepare("select name, created, used, counter, session, public, secret from keys where public = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	user := Key{}
	err = stmt.QueryRow(pub).Scan(&user.Name, &user.Created, &user.Used, &user.Counter, &user.Session, &user.Public, &user.Secret)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *Dal) init() {
	sqlStmt := `
create table keys (
id integer not null primary key AUTOINCREMENT,
name text,
created datetime,
used datetime,
counter int,
session int,
public text,
secret text);
`
	d.db.Exec(sqlStmt)

	sqlStmt = `
create table apps (
id integer not null primary key AUTOINCREMENT,
name text,
created datetime,
key text
)`
	d.db.Exec(sqlStmt)
	return
}

func (d *Dal) Close() {
	d.db.Close()
}
