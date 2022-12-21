package main

import (
  "os"
  "path/filepath"
)

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

  //for _, item := range items {
  //  fmt.Println(item)
  //}

  err = registerCurrentData(items, db)
  if err != nil {
    panic(err)
  }

  err = updateItemMaster(db)
  if err != nil {
    panic(err)
  }

  var updateChkItems []ItemMaster
  updateChkItems, err = getItemMasters(db)
  if err != nil {
    panic(err)
  }

  var updatedItems []ItemMaster
  currentDirectory, _ := os.Getwd()
  downloadBasePath := filepath.Join(currentDirectory, "work", "downloadFiles")

  updatedItems, err = fetchDetailPages(updateChkItems, downloadBasePath)
  if err != nil {
    panic(err)
  }
  if err = registerDetails(db, updatedItems); err != nil {
    panic(err)
  }
}
