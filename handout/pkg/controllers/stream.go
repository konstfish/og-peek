package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
)

func StreamObject(c *gin.Context, minioObj *minio.Object) {
	c.Header("Content-Type", "image/png")

	objInfo, err := minioObj.Stat()
	log.Println(objInfo)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.DataFromReader(
		http.StatusOK,
		objInfo.Size,
		objInfo.ContentType,
		minioObj,
		map[string]string{},
	)
}
