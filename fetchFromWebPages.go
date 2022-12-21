package main

import (
  "fmt"
  "github.com/PuerkitoBio/goquery"
  "net/http"
  "net/url"
  "path/filepath"
  "strconv"
  "strings"
  "time"
)

func getResponse(url string) (*http.Response, error) {
  response, err := http.Get(url)
  if err != nil {
    return nil, fmt.Errorf("http get request error: %w", err)
  }
  return response, nil
}

func getList(response *http.Response) ([]Item, error) {
  body := response.Body

  requestURL := *response.Request.URL

  var items []Item

  doc, err := goquery.NewDocumentFromReader(body)
  if err != nil {
    return nil, fmt.Errorf("gt document error: %w", err)
  }

  tr := doc.Find("table tr")
  notFoundMessage := "ページが存在しません"
  if strings.Contains(doc.Text(), notFoundMessage) || tr.Size() == 0 {
    return nil, nil
  }

  tr.Each(func(_ int, s *goquery.Selection) {
    item := Item{}

    item.Name = s.Find("td:nth-of-type(2) a").Text()
    item.Price, _ = strconv.Atoi(strings.ReplaceAll(strings.ReplaceAll(s.Find("td:nth-of-type(3)").Text(), ",", ""), "円", ""))
    itemURL, exists := s.Find("td:nth-of-type(2) a").Attr("href")
    refURL, parseErr := url.Parse(itemURL)

    if exists && parseErr == nil {
      item.Url = (*requestURL.ResolveReference(refURL)).String()
    }

    if item.Name != "" {
      items = append(items, item)
    }
  })

  return items, nil
}

func fetchDetailPages(items []ItemMaster, downloadBasePath string) ([]ItemMaster, error) {
  parsePage := func(response *http.Response, item ItemMaster) (ItemMaster, error) {
    body := response.Body
    requestURL := *response.Request.URL
    doc, err := goquery.NewDocumentFromReader(body)
    if err != nil {
      return ItemMaster{}, fmt.Errorf("get detail pages document body error %w", err)
    }

    item.Description = doc.Find("table tr:nth-of-type(2) td:nth-of-type(2)").Text()

    href, exists := doc.Find("atble tr:nth-of-type(1) td:nth-of-type(1) img").Attr("src")
    refURL, parseErr := url.Parse(href)
    if exists && parseErr == nil {
      imageURL := (*requestURL.ResolveReference(refURL)).String()
      isUpdated, currentLastModified := checkFileUpdated(imageURL, item.ImageLastModifiedAt)
      if isUpdated {
        item.ImageURL = imageURL
        item.ImageLastModifiedAt = currentLastModified
        imageDownloadPath, err := downloadFile(imageURL, filepath.Join(downloadBasePath, "img", strconv.Itoa(int(item.ID)), item.ImageFileName()))
        if err != nil {
          return ItemMaster{}, fmt.Errorf("download image error: %w", err)
        }
        item.ImageDownloadPath = imageDownloadPath
      }
    }

    href, exists = doc.Find("table tr:nth-of-type(3) td:nth-of-type(2) a").Attr("href")
    refURL, parseErr = url.Parse(href)

    if exists && parseErr == nil {
      pdfURL := (*requestURL.ResolveReference(refURL)).String()
      isUpdated, currentLastModified := checkFileUpdated(pdfURL, item.ImageLastModifiedAt)
      if isUpdated {
        item.PdfURL = pdfURL
        item.PdfLastModifiedAt = currentLastModified
        pdfDownloadPath, err := downloadFile(pdfURL, filepath.Join(downloadBasePath, "pdf", strconv.Itoa(int(item.ID)), item.PdfFileName()))
        if err != nil {
          return ItemMaster{}, fmt.Errorf("download pdf error: %w", err)
        }
        item.PdfDownloadPath = pdfDownloadPath
      }
    }
    return item, nil
  }

  var updatedItems []ItemMaster

  for _, item := range items {
    response, err := getResponse(item.Url)
    if err != nil {
      return nil, fmt.Errorf("fetch detail page body error: %w", err)
    }

    currentIem, err := parsePage(response, item)
    if err != nil {
      return nil, fmt.Errorf("parse detail page content error: %w", err)
    }

    if !item.equals(currentIem) {
      updatedItems = append(updatedItems, currentIem)
    }
  }

  return updatedItems, nil
}

func checkFileUpdated(fileURL string, lastModified time.Time) (isUpdated bool, currentLastModified time.Time) {
  getLastModified := func(fileURL string) (time.Time, error) {
    res, err := http.Head(fileURL)
    if err != nil {
      return time.Unix(0, 0), err
    }
    lastModified, err := time.Parse("Mon, 02 Jan 2006 15:04:05 MST", res.Header.Get("Last-Modified"))
    return lastModified, err
  }

  currentLastModified, err := getLastModified(fileURL)
  if err != nil {
    return false, currentLastModified
  }

  if currentLastModified.After(lastModified) {
    return true, currentLastModified
  } else {
    return false, lastModified
  }
}
