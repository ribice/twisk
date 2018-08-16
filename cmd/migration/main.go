package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/ribice/twisk/internal/iam/secure"
	"github.com/ribice/twisk/model"

	"github.com/go-pg/pg/orm"

	"github.com/go-pg/pg"
)

func main() {
	dbInsert := `INSERT INTO public.tenants VALUES (1, 'admin_tenant', true);
	INSERT INTO public.roles VALUES (1, 'SUPER_ADMIN');
	INSERT INTO public.roles VALUES (2, 'ADMIN');
	INSERT INTO public.roles VALUES (3, 'TENANT_ADMIN');
	INSERT INTO public.roles VALUES (4, 'USER');`

	// postgres://docvsjvu:hR7Zs5Yx58hKSoMOGeYjgJj3BbLZcQvM@horton.elephantsql.com:5432/docvsjvu
	// Provided database is a sample instance from elephantsql.com
	// Feel free to use it for testing. You can easily host one for yourself
	// and fill in the data using cmd/api/migration.go
	var psn = `postgres://postgres:postgres@localhost:5432/postgres`
	queries := strings.Split(dbInsert, ";")

	u, err := pg.ParseURL(psn)
	checkErr(err)
	db := pg.Connect(u)
	_, err = db.Exec("SELECT 1")
	checkErr(err)

	checkErr(db.RunInTransaction(func(tx *pg.Tx) error {
		createSchema(tx, &twisk.Tenant{}, &twisk.Role{}, &twisk.User{})
		for _, v := range queries[0 : len(queries)-1] {
			_, err := tx.Exec(v)
			if err != nil {
				return err
			}
		}
		return insertUser(tx)
	}))

}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func createSchema(db *pg.Tx, models ...interface{}) {
	for _, model := range models {
		checkErr(db.CreateTable(model, &orm.CreateTableOptions{
			FKConstraints: true,
			IfNotExists:   true,
		}))
	}
}

func insertUser(tx *pg.Tx) error {
	s := secure.Service{}
	userInsert := `INSERT INTO public.users VALUES (1, 'John', 'Doe', 'admin', '%s', 'johndoe@mail.com', NULL, NULL, true, null, 1, 1, now(), now(), NULL, NULL);`
	_, err := tx.Exec(fmt.Sprintf(userInsert, s.Hash("admin")))
	return err
}
