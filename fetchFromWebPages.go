package main

import (
  "fmt"
  "github.com/PuerkitoBio/goquery"
  "net/http"
  "net/url"
  "strconv"
  "strings"
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

func fetchDetailPages(items []ItemMaster) ([]ItemMaster, error) {
  parsePage := func(response *http.Response, item ItemMaster) (ItemMaster, error) {
    body := response.Body
    doc, err := goquery.NewDocumentFromReader(body)
    if err != nil {
      return ItemMaster{}, fmt.Errorf("get detail pages document body error %w", err)
    }

    item.Description = doc.Find("table tr:nth-of-type(2) td:nth-of-type(2)").Text()
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
