package main

import (
  "fmt"
  "io"
  "net/http"
  "os"
  "path/filepath"
)

func downloadFile(url string, downloadPath string) (downloadedPath string, err error) {
  err = os.MkdirAll(filepath.Dir(downloadedPath), 0777)
  if err != nil {
    return "", fmt.Errorf("mkdir error during download file: %w", err)
  }

  out, err := os.Create(downloadPath)
  if err != nil {
    return "", fmt.Errorf("create file error during download file: %w", err)
  }
  defer out.Close()

  resp, err := http.Get(url)
  if err != nil {
    return "", fmt.Errorf("download file error: %w", err)
  } else {
    fmt.Println("Download File:", url)
  }
  defer resp.Body.Close()

  _, err = io.Copy(out, resp.Body)
  if err != nil {
    return "", fmt.Errorf("copy file error during download file: %w", err)
  }

  downloadPath = filepath.Join(downloadPath, filepath.Base(downloadPath))
  return downloadPath, nil
}
