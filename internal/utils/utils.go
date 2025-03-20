package utils

import (
	"os"
	"io"
	"fmt"
	"bytes"
	"sort"
	"strings"
	"io/ioutil"
	"path/filepath"
	"netiso-go/pkg/logger"
)

var Log *logger.Logger

type OffsetEntry struct {
	Offset uint64
	Format  string
}

func NewUtils(log *logger.Logger) {
	Log = log
}

// Function to retrieve Xbox ISO files in a given directory
func GetXboxIsoFiles(directoryPath string, silentMode bool) []string {
	if !silentMode {
		Log.Info("Getting ISO files from path: %s", directoryPath)
	}
	var isoFiles []string

	files, err := ioutil.ReadDir(directoryPath)
	if err != nil {
		Log.Warning("Error reading directory: %s", err)
	}

	// Store only XGD formatted ISO files
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".iso") {
			isoPath := filepath.Join(directoryPath, file.Name())
			if IsoIsXboxFormat(isoPath) {
				isoFiles = append(isoFiles, isoPath)
			}
		}
	}

	if len(isoFiles) > 0 {
		sort.Strings(isoFiles)

		if !silentMode {
			Log.Info("Found %d Xbox ISO files", len(isoFiles))
			for index, file := range isoFiles {
				Log.Info("%d - %s", index, filepath.Base(file))
			}
		}
	}

	return isoFiles
}

// Checks ISO file is XGD formatted
func IsoIsXboxFormat(isoFile string) bool {
	xgdId := "MICROSOFT*XBOX*MEDIA"
	buf := make([]byte, len(xgdId))
	offsets := []OffsetEntry {
		{0x10000, "XGD1"},
		{0xfda0000, "XGD2"},
		{0x2090000, "XGD3"},
	}

	file, err := os.Open(isoFile)
	if err != nil {
		Log.Warning("Error opening file: ", err)
		return false
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		Log.Warning("Error getting file info: ", err)
		return false
	}

	for _, element := range offsets {
		// Check if the file length is equal to id length
		if fileInfo.Size() >= int64(element.Offset+uint64(len(xgdId))) {
			
			// Seek to the offset in the file
			_, err := file.Seek(int64(element.Offset), io.SeekStart)
			if err != nil {
				Log.Warning("Error seeking file: ", err)
				return false
			}

			// Read into the buffer
			_, err = file.Read(buf)
			if err != nil {
				Log.Warning("Error reading file: ", err)
				return false
			}

			// Check if the read data matches the magic value
			if bytes.Equal(buf, []byte(xgdId)) {
				Log.Debug("%s ISO: %s", element.Format, isoFile)
				return true
			}
		}
	}
	return false
}

// Helper function to convert a byte slice to a hex string with leading zeros
func ByteSliceToHex(byteSlice []byte) string {
	
	var hexString []string
	for _, b := range byteSlice {
		// Convert each byte to a two-character hex string and append to the slice
		hexString = append(hexString, fmt.Sprintf("%02x", b))
	}
	// Join the slice into a single string separated by spaces
	return strings.Join(hexString, " ")
}

// Function to find the index of a string in a slice
func FindSliceIndex(slice []string, target string) int {
	for i, s := range slice {
		if strings.Contains(s, target) {
			return i // Return the index if found
		}
	}
	return -1 // Return -1 if the target is not found
}