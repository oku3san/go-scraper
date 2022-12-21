package main

import (
  "gorm.io/gorm"
  "path/filepath"
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
  Description         string
  ImageURL            string
  ImageLastModifiedAt time.Time
  ImageDownloadPath   string
  PdfURL              string
  PdfLastModifiedAt   time.Time
  PdfDownloadPath     string
}

func (ItemMaster) TableName() string {
  return "item_master"
}

func (i ItemMaster) equals(target ItemMaster) bool {
  return i.Description == target.Description
}

func (i ItemMaster) equals(target ItemMaster) bool {
  return i.Description == target.Description &&
    i.ImageURL == target.ImageURL &&
    i.ImageLastModifiedAt == target.ImageLastModifiedAt &&
    i.ImageDownloadPath == target.ImageDownloadPath &&
    i.PdfURL == target.PdfURL &&
    i.PdfLastModifiedAt == target.PdfLastModifiedAt &&
    i.PdfDownloadPath == target.PdfDownloadPath
}

func (i ItemMaster) ImageFileName() string {
  return filepath.Base(i.ImageURL)
}

func (i ItemMaster) PdfFileName() string {
  return filepath.Base(i.PdfURL)
}
