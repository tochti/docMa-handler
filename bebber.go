package bebber

import (
  "os"
  "fmt"
  "github.com/gin-gonic/gin"
)

func main() {
  router := gin.Default()
  router.Run(":8080")
}

func GetSettings(k string) string {
  v := os.Getenv(k)
  if v == "" {
    fmt.Sprintf("%v env is missing", k)
    os.Exit(2)
  }

  return v
}

func SubList(a, b []string) []string {
  m := make(map[string]bool)
  for _, v := range b {
    m[v] = false
  }

  var res []string
  for _, v := range a {
    _, ok := m[v]
    if ok == false {
      res = append(res, v)
    }
  }

  return res
}
