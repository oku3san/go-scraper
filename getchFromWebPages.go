package main

import (
  "fmt"
  "net/http"
)

func getResponse(url string) (*http.Response, error) {
  response, err := http.Get(url)
  if err != nil {
    return nil, fmt.Errorf("http get request error: %w", err)
  }
  return response, nil
}
