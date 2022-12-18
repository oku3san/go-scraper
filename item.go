package main

import (
  "gorm.io/gorm"
  "time"
)

type Item struct {
  Name  string
  Price int
  Url   string
}

type LatestItem struct {
  Item
  CreatedAt time.Time
}

type ItemMaster struct {
  gorm.Model
  Item
  Description string
}

func (ItemMaster) TableName() string {
  return "item_master"
}

func (i ItemMaster) equals(target ItemMaster) bool {
  return i.Description == target.Description
}
