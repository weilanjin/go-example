package main

import "github.com/gin-gonic/gin"

func upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": "File upload failed"})
		return
	}

	err = c.SaveUploadedFile(file, "./"+file.Filename)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to save file"})
		return
	}

	c.JSON(200, gin.H{"message": "File uploaded successfully", "filename": file.Filename})
}
