package gosub

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Run boy run
func Run(path, language string) {

	dir, err := os.Open(path)
	defer dir.Close()
	panicOnError(err)

	fi, err := dir.Stat()
	panicOnError(err)

	var moviesFiles []os.FileInfo

	if fi.IsDir() {
		fmt.Printf("Checking %s: %s\n", "directory", filepath.Base(path))
		files, err := ioutil.ReadDir(path)
		panicOnError(err)
		for _, file := range files {
			if !file.IsDir() && strings.HasSuffix(file.Name(), ".mkv") {
				moviesFiles = append(moviesFiles, file)
			}
		}
	} else {
		fmt.Printf("Checking %s: %s\n", "file", filepath.Base(path))
		moviesFiles = append(moviesFiles, fi)
	}

	client := OpenSubtitle{}

	for _, fileInfo := range moviesFiles {

		fmt.Println("Searching for:", fileInfo.Name())

		// fmt.Println(dir + "/" + fileO.Name())
		file, err := os.Open(fileInfo.Name())
		defer file.Close()
		panicOnError(err)

		hash, _ := HashFile(file)

		subtitles, _ := client.Search(OpenSubtitleSearchParameters{
			moviebytesize: fileInfo.Size(),
			moviehash:     fmt.Sprintf("%x", hash),
			sublanguageid: language,
		})

		fmt.Printf("Found %d subtitle/s.\n", len(subtitles))

		// for _, sub := range subtitles {
		// fmt.Println(sub.DownloadLink)
		// client.Download(filePath+".srt", sub.DownloadLink)
		// }
	}

}

const (
	ChunkSize = 65536 // 64k
)

// Generate an OSDB hash for an *os.File.
func HashFile(file *os.File) (hash uint64, err error) {
	fi, err := file.Stat()
	if err != nil {
		return
	}
	if fi.Size() < ChunkSize {
		return 0, fmt.Errorf("File is too small")
	}

	// Read head and tail blocks.
	buf := make([]byte, ChunkSize*2)
	err = readChunk(file, 0, buf[:ChunkSize])
	if err != nil {
		return
	}
	err = readChunk(file, fi.Size()-ChunkSize, buf[ChunkSize:])
	if err != nil {
		return
	}

	// Convert to uint64, and sum.
	var nums [(ChunkSize * 2) / 8]uint64
	reader := bytes.NewReader(buf)
	err = binary.Read(reader, binary.LittleEndian, &nums)
	if err != nil {
		return 0, err
	}
	for _, num := range nums {
		hash += num
	}

	return hash + uint64(fi.Size()), nil
}

// Read a chunk of a file at `offset` so as to fill `buf`.
func readChunk(file *os.File, offset int64, buf []byte) (err error) {
	n, err := file.ReadAt(buf, offset)
	if err != nil {
		return
	}
	if n != ChunkSize {
		return fmt.Errorf("Invalid read %v", n)
	}
	return
}
