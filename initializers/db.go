package go_server

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var (
	tenant_naem string
)

func GetDb(TenantName string, ApiKey string, CheckTenant bool) (*sql.DB, error) {
	var (
		host     = os.Getenv("DB_HOST")
		port     = 5432
		user     = os.Getenv("DB_USER")
		password = os.Getenv("DB_PWD")
		dbname   = TenantName
		tenants  = "Tenants"
	)

	if CheckTenant {
		tenant_naem = ""
		psqlInfo_tenant := fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable",
			host, port, user, password, tenants)
		db, err := sql.Open("postgres", psqlInfo_tenant)
		if err != nil {
			return db, err
		}
		query := fmt.Sprintf("SELECT tenant_name FROM newtable where tenant_name ='%s' and api_key = '%s' ", TenantName, ApiKey)
		rows, err := db.Query(query)
		if err != nil {

			return db, err
		}

		for rows.Next() {
			err := rows.Scan(
				&tenant_naem,
			)

			if err != nil {
				return db, err
			}
		}
		if tenant_naem == "" {

			return db, errors.New("Wrong api_key or tenant_name")
		}
		defer db.Close()
		fmt.Println(tenant_naem, "oooo")

		return db, err
	} else {
		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname)
		db, err := sql.Open("postgres", psqlInfo)
		// if err != nil {
		// 	panic(err)
		// }
		err = db.Ping()
		// if err != nil {
		// 	panic(err)
		// }
		return db, err
	}
}
