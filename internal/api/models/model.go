package models

import "time"

type FilePart struct {
	ID               uint      `gorm:"primaryKey"`                         // Auto-incrementing primary key
	FileID           string    `json:"file_id" gorm:"size:255;index"`      // Identifier for the whole file, indexed
	FileIndex        int       `json:"file_index" gorm:"file_index"`       // Part number or index
	OriginalFileName string    `json:"original_file_name" gorm:"size:255"` // Original file name (optional)
	FileSize         int       `json:"file_size"`                          // File size in bytes
	FileType         string    `json:"file_type"`                          // File type
	FileContent      []byte    `json:"file_content" gorm:"type:bytea"`     // File part content (stored as binary data)
	CreatedAt        time.Time `json:"created_at"`                         // Timestamp when the part was created
	UpdatedAt        time.Time `json:"updated_at"`                         // Timestamp for the last update
}
