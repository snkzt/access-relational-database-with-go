package main

import (
  "database/sql"
  "fmt"
  "log"
  "os"

  "github.com/go-sql-driver/mysql"
)
var db *sql.DB // Database handle. Make this global is not good in production.

func main() {
  // Connection properties
  cfg :=               mysql.Config{
    User:                 os.Getenv("DBUSER"),
    Passwd:               os.Getenv("DBPASS"),
    Net:                  "tcp",
    Addr:                 "127.0.0.1:3306",
    DBName:               "recordings",
    AllowNativePasswords: true,
  }

  // Get a database handle.
  var err error
  db, err = sql.Open("mysql", cfg.FormatDSN())
  // log.Fatal is to check the error in the console only for this, not in production
  if err != nil {
    log.Fatal(err)
  }

  // sql.Open might not immediately connect depends on the driver so db.Ping is to check
  // a connection to the db works.
  pingErr := db.Ping()
  if pingErr != nil {
    log.Fatal(pingErr)
  }
  fmt.Println("Connected!")
}
