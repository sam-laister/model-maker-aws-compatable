package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"

	"github.com/Soup666/diss-api/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRecorder() (*httptest.ResponseRecorder, *gin.Context) {
	recorder := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(recorder)
	c.Set("user", &model.User{Model: gorm.Model{ID: 1}})

	return recorder, c
}

func MockJsonPost(c *gin.Context /* the test context */, content interface{}) {
	jsonbytes, err := json.Marshal(content)
	if err != nil {
		panic(err)
	}

	// the request body must be an io.ReadCloser
	// the bytes buffer though doesn't implement io.Closer,
	// so you wrap it in a no-op closer
	c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonbytes))
}
