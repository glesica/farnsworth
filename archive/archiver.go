package zip

import "io"

// Archiver abstracts the concept of an archive to support multiple
// archive formats and facilitate testing.
type Archiver interface {
	AddFile(path string, content io.Reader) error
	Write(path string) error
}

// TODO: ZipArchiver
// TODO: MockArchiver
