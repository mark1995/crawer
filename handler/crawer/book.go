package crawer

import (
	"bookcrawer/handler"
	"bookcrawer/parser"
	"github.com/gin-gonic/gin"

	"net/http"
)

func Book(c *gin.Context) {
	url := c.PostForm("url")
	go func() {
		bookp, err := parser.NewParser("booktxt")
		if err != nil {
			handler.SendResponse(c, http.ErrServerClosed, nil)
		}
		bookp.ParserUrl(url)
	}()
	c.String(http.StatusOK, "ok", "")
}
