package storage

import (
	"context"
	"io"
	"net/http"
	"time"
)

type Visibility string

const (
	Public        Visibility = "public"
	Private       Visibility = "private"
	Authenticated Visibility = "authenticated" // e.g., user must be logged in
)

type UploadOptions struct {
	ContentType string            // e.g., "image/png", "application/pdf"
	Folder      string            // Path/folder in the storage (e.g., "products/images")
	FileName    string            // Custom file name (without extension)
	Tags        []string          // Optional tags (used for Cloudinary, indexing, etc.)
	Metadata    map[string]string // Custom metadata (stored in DB or on storage)
	Encrypted   bool              // Whether to encrypt the content URL
	OwnerID     string            // Optional: Associate with user/vendor/store
	Visibility  Visibility        // Enum: Public, Private, Authenticated
	ContentHash string            // Optional checksum or hash for deduplication
}
type ContentMetadata struct {
	ID          string            // Unique content ID (could be UUID or storage-specific)
	StoragePath string            // Path or key in the storage bucket
	URL         string            // Accessible URL (could be signed, encrypted, or raw)
	ContentType string            // e.g., "image/jpeg", "application/pdf"
	Size        int64             // Size in bytes
	UploadedAt  time.Time         // Timestamp of upload
	Tags        []string          // Tags associated with content
	Metadata    map[string]string // Custom metadata (category, description, etc.)
	Visibility  Visibility        // Public, Private, or Authenticated
	Encrypted   bool              // Indicates if the URL or content is encrypted
	OwnerID     string            // Associated user/vendor/store ID
	ContentHash string            // Optional SHA-256 hash (for integrity/deduplication)
}

type StorageAdapter interface {
	Upload(ctx context.Context, file io.Reader, fileName string, opts UploadOptions) (ContentMetadata, error)
	Stream(ctx context.Context, encryptedURL string, w http.ResponseWriter) error
	Delete(ctx context.Context, encryptedURL string) error
	GetDecryptedURL(ctx context.Context, encryptedURL string) (string, error)
}
