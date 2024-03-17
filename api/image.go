package api

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/erlanggawulung/shopifyx/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (server *Server) postImage(ctx *gin.Context) {
	// Parse the multipart form
	err := ctx.Request.ParseMultipartForm(2 << 20) // Max size 2MB
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse multipart form"})
		return
	}

	// Retrieve the file from the form
	file, fileHeader, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve file from form"})
		return
	}
	defer file.Close()

	// Validate file size
	fileSize := fileHeader.Size
	if fileSize < 10240 { // 10KB
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "File size too small, minimum size is 10KB"})
		return
	}

	// Validate file type
	fileType := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if fileType != ".jpg" && fileType != ".jpeg" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type, only JPG or JPEG allowed"})
		return
	}

	// Save the file to disk or process it as needed
	// Example: Save the file to a predefined directory
	// You need to adjust the directory path according to your setup
	// Make sure the directory exists and has appropriate permissions
	// Generate a new UUID
	fileName := uuid.New().String()

	// Uncomment this line if you want to save file in local storage
	// filePath := "images/" + fileName + fileType // Example path
	// err = ctx.SaveUploadedFile(fileHeader, filePath)

	s3Url, err := uploadToS3(file, fileName+fileType, server.config)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Respond with success
	ctx.JSON(http.StatusOK, gin.H{"imageUrl": s3Url})
}

func uploadToS3(file io.Reader, fileNameWithExt string, config util.Config) (string, error) {
	// Read the content into a buffer
	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, file); err != nil {
		return "", err
	}

	// Create AWS session with credentials
	creds := credentials.NewStaticCredentials(config.S3ID, config.S3SecretKey, "")
	_, err := creds.Get()
	if err != nil {
		return "", err
	}

	// Configure AWS session
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("ap-southeast-1"),
		Credentials: creds,
	})
	if err != nil {
		return "", err
	}

	// Create S3 client
	svc := s3.New(sess)

	// Prepare the parameters for the S3 upload
	params := &s3.PutObjectInput{
		Bucket: aws.String(config.S3BucketName),
		Key:    aws.String(fileNameWithExt),
		Body:   bytes.NewReader(buf.Bytes()),
		ACL:    aws.String("public-read"),
	}

	// Upload the file to S3
	_, err = svc.PutObject(params)
	if err != nil {
		return "", err
	}

	// Construct the URL to access the uploaded file
	url := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", config.S3BucketName, fileNameWithExt)

	return url, nil
}
