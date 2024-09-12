package models

import "time"

type HealthResp struct {
	Success     bool   `json:"success" example:"true"`
	Data        string `json:"data" example:"null"`
	Error       string `json:"error" example:""`
	Description string `json:"description" example:"ok"`
}

type UploadResp struct {
	Success     bool   `json:"success" example:"true"`
	Data        string `json:"data" example:"57c18389-32cc-4248-9e74-47d27658456e"`
	Error       string `json:"error" example:""`
	Description string `json:"description" example:"File uploaded and split"`
}

type GetFileDataResp struct {
	Success     bool            `json:"success" example:"true"`
	Data        FilePartSwagger `json:"data"`
	Error       string          `json:"error" example:""`
	Description string          `json:"description" example:"Getting file data has succeeded"`
}

type FilePartSwagger struct {
	ID               uint      `gorm:"primaryKey" example:"31"`                                // Auto-incrementing primary key
	FileID           string    `json:"file_id" example:"57c18389-32cc-4248-9e74-47d27658456e"` // Identifier for the whole file, indexed
	FileIndex        int       `json:"file_index" example:"4"`                                 // Part number or index
	OriginalFileName string    `json:"original_file_name" example:"one.txt"`                   // Original file name (optional)
	FileSize         int       `json:"file_size" example:"154"`                                // File size in bytes
	FileType         string    `json:"file_type" example:"text/plain"`                         // File type
	FileContent      string    `json:"file_content" example:"dAojIGt1a3UKdGhpcyBpcyBhIGN1"`    // File part content (stored as binary data)
	CreatedAt        time.Time `json:"created_at" example:"2024-09-12T17:39:51.792128+09:00"`  // Timestamp when the part was created
	UpdatedAt        time.Time `json:"updated_at" example:"2024-09-12T17:39:51.792128+09:00"`  // Timestamp for the last update
}

type DownloadResp struct {
	Success     bool   `json:"success" example:"true"`
	Data        string `json:"data" example:"data file"`
	Error       string `json:"error" example:""`
	Description string `json:"description" example:"File uploaded and split"`
}

type ErrorResp struct {
	Success     bool   `json:"success" example:"false"`
	Data        string `json:"data" example:"null"`
	Error       string `json:"error" example:"Failed to retrieve files"`
	Description string `json:"description" example:"Failed to retrieve files"`
}
