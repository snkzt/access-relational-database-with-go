package main

import (
  "database/sql"
  "fmt"
  "log"
  "os"

  "github.com/go-sql-driver/mysql"
)

var db *sql.DB // Database handle. Make this global is not good in production.

type Album  struct {
  ID      int64
  Title   string
  Artist  string
  Price   float
}

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

// Retrieve album(s) by artist name
func albumsByArtist(name string) ([]Album, error) {
  var albums []Album // Type Album slice
 
  // Separate query and parameter enables the database/sql package to
  // send them separately that removes SQL injection risk compare to concatinate them.
  rows, err := db.Query("SELECT * FROM album WHERE artist = ?", name)
  if err != nil {
    return nil, fmt.Errorf("albumByArtist %q: %v", name, err)
  }
  // Defer close rows so that it will release the resourse it holds at the end.
  defer rows.Close()
 
  // Loop through rows and assign column data to struct fields with Scan.
  for rows.Next() {
    var album Album
    if err := rows.Scan(&album.ID, &album.Title, &album.Artist, &album.Price); err != nil {
      return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
    }
    albums = append(albums, album)
  }
  if err := rows.Err(); err != nil {
    return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
  }
  return albums, nil
}

