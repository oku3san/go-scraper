package main

import "fmt"

func main() {

  db, err := gormConnect()
  if err != nil {
    panic(err)
  }

  err = dbMigration(db)
  if err != nil {
    panic(err)
  }

  baseUrl := "http://localhost:8080/"
  response, err := getResponse(baseUrl)
  if err != nil {
    panic(err)
  }

  items, err := getList(response)
  if err != nil {
    panic(err)
  }

  for _, item := range items {
    fmt.Println(item)
  }
}
