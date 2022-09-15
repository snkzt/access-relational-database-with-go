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
  Price   float32
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

  albums, err := albumsByArtist("Nina Simone")
  if err != nil {
    log.Fatal(err)
  }
  fmt.Printf("Albums found: %v\n", albums)

  album, err := albumById(4)
  if err != nil {
    log.Fatal(err)
  }
  fmt.Printf("Album found: %v\n", album)

  albumID, err := addAlbum(Album{
    Title:  "The Modern Sound of Betty Carter",
    Artist: "Betty Carter",
    Price: 49.99, 
  })
  if err != nil {
    log.Fatal(err)
  }
  fmt.Printf("ID of added album: %v\n", albumID)
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

// albumById queries for the album with the specified ID.
func albumById(id int64) (Album, error) {
  // An album to hold data from the returned row.
  var album Album

  row := db.QueryRow("SELECT * FROM album WHERE id = ?", id)
  if err := row.Scan(&album.ID, &album.Title, &album.Artist, &album.Price); err != nil {
    if err == sql.ErrNoRows {
      return album, fmt.Errorf("albumById %d: no such album", id)
    }
    return album, fmt.Errorf("albumById %d: %v", id, err)
  }
  return album, nil
}

// addAlbum adds the specified album to the record database,
// returning the album ID of the new entry.
func addAlbum(album Album) (int64, error) {
  result, err := db.Exec("INSERT INTO album (title, artist, price) VALUES (?, ?, ?)", album.Title, album.Artist, album.Price)
  if err != nil {
    return 0, fmt.Errorf("addAlbum: %v", err)
  }
  id, err := result.LastInsertId()
  if err != nil {
    return 0, fmt.Errorf("addAlbum: %v", err)
  }
  return id, nil
}
