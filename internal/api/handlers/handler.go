package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rocky2015aaa/filestorageservice/internal/api/models"
	"github.com/rocky2015aaa/filestorageservice/internal/config"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Handler struct {
	Database *gorm.DB
}

func NewHandler(database *gorm.DB) *Handler {
	return &Handler{
		Database: database,
	}
}

func getResponse(success bool, data interface{}, err, description string) models.Resp {
	return models.Resp{
		Success:     success,
		Data:        data,
		Error:       err,
		Description: description,
	}
}

// PostSwapOrder godoc
// @Title        Ping
// @Description  Check Server Status
// @Tags         Ping
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.HealthResp
// @Router       /api/v1/ping [get]
func (h *Handler) Ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, getResponse(true, nil, "", "ok"))
}

// PostSwapOrder godoc
// @Title        Upload File
// @Description  Upload a file
// @Tags         File
// @Accept       json
// @Produce      json
// @Param        file formData file true "File metadata and content"
// @Success      201  {object}  models.UploadResp
// @Failure      400  {object}  models.ErrorResp
// @Failure      500  {object}  models.ErrorResp
// @Router       /api/v1/upload [post]
func (h *Handler) UploadHandler(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		log.Error("Error retrieving the file:", err)
		ctx.JSON(http.StatusBadRequest, getResponse(false, nil, err.Error(), "Invalid file"))
		return
	}

	openedFile, err := file.Open()
	if err != nil {
		log.Error("Error opening the file:", err)
		ctx.JSON(http.StatusInternalServerError, getResponse(false, nil, err.Error(), "File open error"))
		return
	}
	defer openedFile.Close()

	// Split file into parts
	flieSliceNumber, err := strconv.Atoi(os.Getenv(config.EnvSvrFileSliceNumber))
	if err != nil {
		log.Error("Error converting the number of file slice:", err)
		ctx.JSON(http.StatusInternalServerError, getResponse(false, nil, err.Error(), "File split error"))
		return
	}
	parts, err := splitFile(openedFile, flieSliceNumber) // Split into 5 parts
	if err != nil {
		log.Error("Error split the files:", err)
		ctx.JSON(http.StatusInternalServerError, getResponse(false, nil, err.Error(), "File split error"))
		return
	}

	// Store each part in the database in parallel
	var wg sync.WaitGroup
	fileID := uuid.New().String()
	wg.Add(len(parts))
	for index, part := range parts {
		go func(part []byte, fileIndex int) {
			defer wg.Done()
			// Insert multiple parts of a file
			filePart := models.FilePart{
				FileID:           fileID,
				FileIndex:        fileIndex,
				OriginalFileName: file.Filename,
				FileSize:         len(part),
				FileType:         file.Header.Get("Content-Type"),
				FileContent:      part,
			}
			err := h.Database.Debug().Create(&filePart).Error
			if err != nil {
				log.Error("Error insert the file data:", err)
				ctx.JSON(http.StatusInternalServerError,
					getResponse(false, nil, err.Error(), "Creating file data has failed"))
				return
			}
		}(part, index)
	}
	wg.Wait()

	// Respond with success
	ctx.JSON(http.StatusCreated, getResponse(true, fileID, "", "File uploaded and split"))
}

// PostSwapOrder godoc
// @Title      Get file data
// @Description  Get subfile data
// @Tags         File
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.GetFileDataResp
// @Failure      500  {object}  models.ErrorResp
// @Router       /api/v1/files-data [get]
func (h *Handler) GetFilesDataHandler(ctx *gin.Context) {
	// Retrieve and list all files
	var files []models.FilePart
	err := h.Database.Model(&models.FilePart{}).
		Order("original_file_name ASC, file_index ASC").Find(&files).Error
	if err != nil {
		log.Error("Error get the file data:", err)
		ctx.JSON(http.StatusInternalServerError, getResponse(false, nil, err.Error(), "Failed to retrieve files"))
		return
	}
	ctx.JSON(http.StatusOK, getResponse(true, files, "", "Getting file data has succeeded"))
}

// PostSwapOrder godoc
// @Title        Get File
// @Description  Get file
// @Tags         File
// @Accept       json
// @Produce      json
// @Param file_id path string true "File ID"
// @Success      200  {object}  models.DownloadResp
// @Failure      400  {object}  models.ErrorResp
// @Failure      500  {object}  models.ErrorResp
// @Router       /api/v1/download [get]
func (h *Handler) DownloadHandler(ctx *gin.Context) {
	fileID := ctx.Query("file_id")
	if fileID == "" {
		log.Error("Error get file_id:", fmt.Errorf("file_id is missing"))
		ctx.JSON(http.StatusBadRequest, getResponse(false, nil, "File ID is required", "File ID is required"))
		return
	}

	/*
		Normal process
	*/
	// var parts []models.FilePart
	// err := h.Database.Debug().Where("file_id = ?", fileID).Order("file_index").Find(&parts).Error
	// if err == gorm.ErrRecordNotFound {
	// 	ctx.JSON(http.StatusNotFound, getResponse(false, nil, err.Error(), "Failed to retrieve files"))
	// 	return
	// } else if err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, getResponse(false, nil, err.Error(), "Failed to retrieve files"))
	// 	return
	// }
	// var mergedFile bytes.Buffer
	// for _, part := range parts {
	// 	mergedFile.Write(part.FileContent)
	// }
	// ctx.Header("Content-Disposition", "attachment; filename="+parts[0].OriginalFileName)
	// ctx.Header("Content-Type", parts[0].FileType)
	// ctx.Data(http.StatusOK, parts[0].FileType, mergedFile.Bytes())

	/*
		The process in the condition:
		Get files parallel using threads and merge back to one file
	*/
	fileInfo := struct {
		MaxFileIndex     int
		OriginalFileName string
		FileType         string
	}{}
	err := h.Database.Raw(`
		SELECT
			file_id,
			file_type,
			file_index AS max_file_index,
			original_file_name
		FROM file_parts
		WHERE file_index = (
			SELECT MAX(file_index)
			FROM file_parts
			WHERE file_id = ?
		)
		AND file_id = ?
	`, fileID, fileID).Scan(&fileInfo).Error
	if err != nil {
		log.Error("Error get max file_index:", err)
		ctx.JSON(http.StatusInternalServerError, getResponse(false, nil, err.Error(), "Failed to retrieve max file index"))
		return
	}

	// Merge file parts
	var mergedFile bytes.Buffer
	var wg sync.WaitGroup
	fileContents := make([][]byte, fileInfo.MaxFileIndex+1)
	wg.Add(fileInfo.MaxFileIndex + 1)
	for i := 0; i <= fileInfo.MaxFileIndex; i++ {
		go func(index int, fileID string, fileContents [][]byte) {
			defer wg.Done()
			var part models.FilePart
			err := h.Database.Model(&models.FilePart{}).
				Where("file_id = ?", fileID).
				Where("file_index = ?", index).
				First(&part).Error
			if err != nil {
				log.Printf("Error getting file content for index %d: %v", index, err)
				return
			}
			fileContents[index] = part.FileContent
		}(i, fileID, fileContents)
	}
	wg.Wait()
	// Merging always need to be in order not to break the file
	for _, fileContent := range fileContents {
		mergedFile.Write(fileContent)
	}

	// Set response headers and write file content
	ctx.Header("Content-Disposition", "attachment; filename="+fileInfo.OriginalFileName)
	ctx.Header("Content-Type", fileInfo.FileType)
	ctx.Data(http.StatusOK, fileInfo.FileType, mergedFile.Bytes())
}
