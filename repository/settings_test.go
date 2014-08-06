package repository

import (
  "fmt"
  "path/filepath"
  "os"
)

func ExampleSettings() {
  path, pathError := filepath.Abs("../settings.json")

  if pathError != nil {
    panic(pathError)
  }

  os.Remove(path)

  Settings().StorePath = path

  fmt.Println(Settings().Get("consumer_key") == "")

  Settings().Set("consumer_key", "123")
  fmt.Println(Settings().Get("consumer_key"))

  // Output:
  // true
  // 123
}
