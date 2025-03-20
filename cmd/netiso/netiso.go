package main

import (
	"fmt"
	"net"
	"os"
	"io"
	"time"
	"flag"
	"strings"
	"path/filepath"
	"netiso-go/pkg/logger"
	"netiso-go/internal/utils"
)

var Log *logger.Logger
var logLevel int
var listenAddr string
var isoPath string
var logSpeed bool

func initLogger(logLevel int) {
	Log = logger.NewLogger(logLevel)
}

func initUtils(log *logger.Logger) {
	utils.NewUtils(Log)
}

func parseParams() {
	// Define flags
	flag.StringVar(&isoPath, "f", "", "Path to the Xbox ISO files.")
	flag.IntVar(&logLevel, "l", int(logger.InfoLevel), "Set log level (0: Error, 1: Warning, 2: Info, 3: Debug)")
	flag.BoolVar(&logSpeed, "s", false, "Logs network speed (mbps).")
	debugFlag := flag.Bool("d", false, "Enable debug mode.")
	
	// Parse the command-line flags
	flag.Parse()

	if isoPath == "" {
		Log.Warning("-f flag is not set, using working directory as ISO path.")
	}

	// Enable Debug mode if -d flag is passed
	if *debugFlag {
		Log.SetLevel(logger.DebugLevel)
		Log.Debug("Debug mode is enabled.")
	} else {
		Log.SetLevel(logLevel) // Set log level from the flag (-l)
	}
}

func main() {
	serverAddress := "0.0.0.0:4323"

	initLogger(logLevel)
	parseParams()
	initUtils(Log)

	// Verify ISO path
	isoPath = filepath.Clean(isoPath) // Normalize string (windows/linux format)
	info, err := os.Stat(isoPath)
	if err != nil {
		Log.Error("Error checking directory: %v\n", err)
	}
	if !info.IsDir() {
		Log.Error("%s is not a valid directory.\n", isoPath)
	}

	listen, err := net.Listen("tcp", serverAddress)
	if err != nil {
		Log.Error("Error starting server: %s", err)
	}
	defer listen.Close()

	Log.Info("Starting NetISO server on port 4323")
	for {
		conn, err := listen.Accept()
		if err != nil {
			Log.Warning("Error connecting client: %s", err)
			continue
		} else {
			Log.Info("Client connected: %s", conn.RemoteAddr().String())
		}

		go handleConnection(conn)
	}
}

// Function to display speed every second
func calculateTransferSpeed(ticker *time.Ticker, totalBytesSent *int64) {
	for {
		<-ticker.C
		speedMbps := float64(*totalBytesSent) / 1_000_000 * 8
		fmt.Printf("\r%s", strings.Repeat(" ", 80)) // Clear the line
		fmt.Printf("\rStreaming speed: %.2f Mbps", speedMbps)
		*totalBytesSent = 0
	}
}


// Handle an incoming connection
func handleConnection(conn net.Conn) {
	// Ensure the connection is closed after handling the request
	defer conn.Close()

	var totalBytesSent int64
	if logSpeed {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		go calculateTransferSpeed(ticker, &totalBytesSent)
	}

	isoFiles := utils.GetXboxIsoFiles(isoPath, false)
	isoFilesReloaded := false
	serverOkResponse := []byte("ISVRokOK")
	fileMounted := false
	fileIndex := 0
	selectedIso := ""
	selectedIsoPath := ""
	inactivityTimeout := 10 * time.Second
	
	conn.SetDeadline(time.Now().Add(inactivityTimeout))
	for {
		var byteResponse []byte
		clientRequest := make([]byte,264)

		_, err := conn.Read(clientRequest)
		if err != nil {
			Log.Warning("Client disconnected: %v", err)
			return
		}

		clientCommand := clientRequest[5]

		// File mount request, move to streaming ISO
		if fileMounted && clientRequest[0] == 0x5c {
			clientCommand = 5
		}

		// Debug client reponse
		if Log.GetLevel() == logger.DebugLevel {
			clientRequestHex := utils.ByteSliceToHex(clientRequest)
			if (clientCommand != 5) {
				clientRequestHex = clientRequestHex[:60]
			}
			Log.Debug("Client request (hex): %s\n", clientRequestHex)
		}

		// Client command switch
		switch clientCommand {
		case 0:
			Log.Debug("0: Server Ok request")
			byteResponse = serverOkResponse
		case 1:
			Log.Debug("1: Client signal ISO stream request")
			if fileMounted {
				byteResponse = []byte{0x00, 0x3a, 0x67, 0x20, 0x00, 0x00, 0x08, 0x00}
				fileMounted = false
				Log.Debug("Send game mount command")
			} else {
				byteResponse = make([]byte, 8)
				Log.Debug("Send 8 null byte command")
			}
		case 2:
			Log.Debug("2: 4 null byte request")
			byteResponse = make([]byte, 4)
		case 3:
			Log.Debug("3: ISO streaming request")
			
			if (selectedIsoPath != "") {
				bufferSize := int(clientRequest[18]) * 256
				hexOffset := int64(clientRequest[11])<<32 | int64(clientRequest[12])<<24 | int64(clientRequest[13])<<16 | int64(clientRequest[14])<<8 | int64(clientRequest[15])
				
				Log.Debug("Seek offset: %d | Seek buffer size: %d", hexOffset, bufferSize)

				file, err := os.Open(selectedIsoPath)
				if err != nil {
					Log.Warning("ISO open error: %s", err)
					byteResponse = make([]byte, 4)
				}
				defer file.Close()

				_, err = file.Seek(hexOffset, io.SeekStart)
				if err != nil {
					Log.Warning("ISO seek error: %s", err)
					byteResponse = make([]byte, 4)
				}

				byteResponse = make([]byte, bufferSize)
				_, err = file.Read(byteResponse)
				if err != nil && err != io.EOF {
					Log.Warning("ISO read error: %s", err)
					byteResponse = make([]byte, 4)
				}
			}
		case 4:
			Log.Debug("4: ISO List request")
			responseLength := 264

			// Only reload if changed
			if !isoFilesReloaded {
				newFiles := utils.GetXboxIsoFiles(isoPath, true)
				for i := range newFiles {
					if isoFiles[i] != newFiles[i] {
						isoFiles = utils.GetXboxIsoFiles(isoPath, false)
					}
				}
				isoFilesReloaded = true
			}
			
			if fileIndex < len(isoFiles) {
				// Build ISO string without extension
				_, fileName := filepath.Split(isoFiles[fileIndex])
				fileName = fileName[:len(fileName)-len(filepath.Ext(fileName))]
				byteResponse = []byte(fileName)
		
				// Pad out string to match required length
				if len(byteResponse) < responseLength {
					padding := responseLength - len(byteResponse)
					byteResponse = append(byteResponse, make([]byte, padding)...)
				}
				fileIndex++
			} else {
				byteResponse = make([]byte, responseLength)
				fileIndex = 0
				isoFilesReloaded = false
			}
		case 5:
			fileMounted = true
			selectedIso = string(clientRequest)
			Log.Debug("5: ISO mount request")

			if strings.Contains(selectedIso, "[Disable Current ISO]") {
				Log.Info("Unmounting ISO")
				byteResponse = make([]byte, 4)
				fileMounted = false
			} else if strings.Contains(selectedIso,"\\Mount\\") {
				Log.Debug("Client request ISO: %s", selectedIso)
				byteResponse = []byte{0x00, 0x00, 0x00, 0x01}
				selectedGame := strings.Split(selectedIso,"\\")[1]
				index := utils.FindSliceIndex(isoFiles, selectedGame)
				if index != -1 {
					Log.Info("Found game file, mounting: %s", selectedGame)
					selectedIsoPath = isoFiles[index]
				} else {
					Log.Warning("Game file '%s' not found in path!", selectedGame)
					byteResponse = make([]byte, 4)
				}
			}
		default:
			Log.Warning("Received unknown request: %x", clientRequest[4])
			byteResponse = make([]byte, 4)
		}

		_, err = conn.Write(byteResponse)
		if err != nil {
			Log.Warning("Error sending response to client: %s", err)
			return
		}

		if (Log.GetLevel() == logger.DebugLevel) {
			Log.Debug("Sent response to client (hex): %x", byteResponse)
		}

		if (logSpeed) {
			totalBytesSent += int64(len(byteResponse))
		}

		conn.SetDeadline(time.Now().Add(inactivityTimeout))
	}
}