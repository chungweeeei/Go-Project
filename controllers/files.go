package controllers

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type UploadFileManager struct {
	FileId     string `binding:"required"`
	Filename   string `binding:"required"`
	StartByte  int64  `binding:"required"`
	TotalBytes int64  `binding:"required"`
	Content    bytes.Buffer
	StartedAt  time.Time
	FinishedAt time.Time
}

func (ufm *UploadFileManager) uploadChunk(chunk []byte) error {

	_, err := ufm.Content.Write(chunk)
	return err
}

func (utf *UploadFileManager) deployResult() error {

	file, err := os.Create("./" + utf.Filename)
	if err != nil {
		return errors.New("Failed to create file")
	}

	defer file.Close()

	_, err = file.Write(utf.Content.Bytes())
	if err != nil {
		return errors.New("Failed to write content to file")
	}

	return nil
}

func Upload(context *gin.Context) {
	// Bind the query parameters to the UploadFileManager struct
	fileId := context.GetHeader("X-fileId")
	filename := context.GetHeader("X-filename")
	startByteStr := context.GetHeader("X-startByte")
	totalBytesStr := context.GetHeader("Content-Length")

	if fileId == "" || filename == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Missing required headers"})
		return
	}

	startByte, err := strconv.ParseInt(startByteStr, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "invalid start_byte header"})
		return
	}

	totalBytes, err := strconv.ParseInt(totalBytesStr, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "invalid total_bytes header"})
		return
	}

	uploadFile := UploadFileManager{
		FileId:     fileId,
		Filename:   filename,
		StartByte:  startByte,
		TotalBytes: totalBytes,
		Content:    bytes.Buffer{},
		StartedAt:  time.Now(),
	}

	defer context.Request.Body.Close()

	// Read streaming data in chunks
	chunk := make([]byte, 1024*1024) // 1 MB chunk size
	totalBytesReceived := int64(0)

	for {
		n, err := context.Request.Body.Read(chunk)
		if n > 0 {
			// Process the chunk
			chunkData := chunk[:n]
			err = uploadFile.uploadChunk(chunkData)
			if err != nil {
				context.JSON(http.StatusInternalServerError, gin.H{"message": "Error uploading chunk"})
				return
			}
			totalBytesReceived += int64(n)
		}

		// Check for end of stream
		if err == io.EOF {
			break
		}

		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"message": "Error reading streaming data"})
			return
		}
	}

	uploadFile.FinishedAt = time.Now()

	fmt.Println("Finished uploading:", uploadFile.Content.Len(), "bytes")

	uploadFile.deployResult()

	context.JSON(http.StatusOK, gin.H{"message": "File upload started"})
}
