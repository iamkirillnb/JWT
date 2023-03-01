package mysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"github.com/uptrace/bun/extra/bundebug"
	"log"
	"os"
	"runtime"
)

// setConnection хелпер устанавливает подключение к БД mysql
// для подключения необходимо прокинуть STORE_DSN в переменные окружения
// для включения режима отладки STORE_DEBUG
func setConnection(debug bool) (connection *bun.DB, err error) {
	if dsn, ok := os.LookupEnv("STORE_DSN"); ok {
		sqldb, err := sql.Open("mysql", dsn)
		if err != nil {
			log.Fatalln("Could not set connection to:", err)
			return nil, err
		}
		// Create a Bun db on top of it.
		connection = bun.NewDB(sqldb, mysqldialect.New())
		if err := connection.Ping(); err != nil {
			log.Fatalln("Error establish connection to store MySQL:", err.Error())
		}
		maxOpenConns := 4 * runtime.GOMAXPROCS(0)
		sqldb.SetMaxOpenConns(maxOpenConns)
		sqldb.SetMaxIdleConns(maxOpenConns)
	} else {
		log.Fatalln("Not found STORE_DSN in env")
		return nil, err
	}

	connection.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(debug)))
	return connection, nil
}
