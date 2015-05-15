package bebber

import (
  "github.com/gin-gonic/gin"
)

func main() {
  router := gin.Default()
  router.Run(":8080")
}

