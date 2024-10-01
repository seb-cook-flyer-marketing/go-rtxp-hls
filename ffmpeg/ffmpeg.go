package ffmpeg

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

type RunningPid struct {
	StreamID string
	PID      int
}

var (
	runningStreams []RunningPid
	mu             sync.Mutex
	config         = getConfig()
)

type FFmpegConfig struct {
	Ffmpeg string
}

func getConfig() FFmpegConfig {
	// Initialize your configuration here
	return FFmpegConfig{
		Ffmpeg: "ffmpeg", // Adjust the path if necessary
	}
}

// StopStream terminates the FFmpeg process for a given stream ID.
func StopStream(streamID string) (bool, error) {
	mu.Lock()
	defer mu.Unlock()

	for i, stream := range runningStreams {
		if stream.StreamID == streamID {
			log.Printf("Killing previous process PID: %d", stream.PID)
			process, err := os.FindProcess(stream.PID)
			if err != nil {
				log.Printf("Cannot find PID: %d", stream.PID)
				return false, err
			}

			err = process.Kill()
			if err != nil {
				log.Printf("Cannot kill PID: %d", stream.PID)
				return false, err
			}

			// Remove the stream from runningStreams
			runningStreams = append(runningStreams[:i], runningStreams[i+1:]...)
			return true, nil
		}
	}

	// If the stream was not found, consider it already stopped
	return false, nil
}

// ConvertStream starts the FFmpeg process to convert a stream.
func ConvertStream(streamURL, streamID string) (string, error) {
	mu.Lock()
	// Check if the stream is already running
	for _, stream := range runningStreams {
		if stream.StreamID == streamID {
			mu.Unlock()
			return getRelativePath(getOutputPath(streamID)), nil
		}
	}
	mu.Unlock()

	outputPath := getOutputPath(streamID)
	fullPath := filepath.Join(outputPath, "stream.m3u8")

	// Define FFmpeg command parameters
	cmdParams := []string{
		"-i", streamURL,
		"-rtsp_flags", "prefer_tcp",
		"-c:v", "libx264",
		"-crf", "21",
		"-preset", "ultrafast",
		"-sc_threshold", "0",
		"-c:a", "aac",
		"-b:a", "128k",
		"-ac", "2",
		"-f", "hls",
		"-hls_time", "2",
		"-hls_delete_threshold", "1",
		"-hls_flags", "delete_segments",
		"-hls_list_size", "30",
		"-g", "48", // GOP
		fullPath,
	}

	cmd := exec.Command(config.Ffmpeg, cmdParams...)
	log.Printf("FFMPEG Command: %s %v", config.Ffmpeg, cmdParams)

	// Start the FFmpeg process
	if err := cmd.Start(); err != nil {
		log.Printf("Error starting FFmpeg: %v", err)
		return "", err
	}

	mu.Lock()
	runningStreams = append(runningStreams, RunningPid{
		StreamID: streamID,
		PID:      cmd.Process.Pid,
	})
	mu.Unlock()

	// Handle process output
	go func() {
		stdout, err := cmd.StdoutPipe()
		if err == nil {
			go func() {
				buf := make([]byte, 1024)
				for {
					n, err := stdout.Read(buf)
					if n > 0 {
						log.Printf("STDOUT: %s", string(buf[:n]))
					}
					if err != nil {
						break
					}
				}
			}()
		}

		stderr, err := cmd.StderrPipe()
		if err == nil {
			go func() {
				buf := make([]byte, 1024)
				for {
					n, err := stderr.Read(buf)
					if n > 0 {
						log.Printf("STDERR: %s", string(buf[:n]))
					}
					if err != nil {
						break
					}
				}
			}()
		}

		// Wait for the process to finish
		err = cmd.Wait()
		mu.Lock()
		// Remove the stream from runningStreams
		for i, stream := range runningStreams {
			if stream.StreamID == streamID {
				runningStreams = append(runningStreams[:i], runningStreams[i+1:]...)
				break
			}
		}
		mu.Unlock()

		if err != nil {
			log.Printf("Process exited with error: %v", err)
			if exitError, ok := err.(*exec.ExitError); ok {
				if exitError.ExitCode() != 255 {
					time.Sleep(1 * time.Second)
					_, _ = ConvertStream(streamURL, streamID)
				}
			}
		} else {
			log.Printf("Stream %s has ended.", streamID)
		}
	}()

	// Wait a second before returning
	time.Sleep(1 * time.Second)
	return getRelativePath(fullPath), nil
}

// getOutputPath constructs and ensures the output directory exists.
func getOutputPath(streamID string) string {
	outputPath := filepath.Join("public", streamID)

	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		err := os.MkdirAll(outputPath, os.ModePerm)
		if err != nil {
			log.Fatalf("Failed to create directory %s: %v", outputPath, err)
		}
	}

	return outputPath
}

// getRelativePath returns the path relative to the public directory.
func getRelativePath(fullPath string) string {
	relPath, err := filepath.Rel("public", fullPath)
	if err != nil {
		log.Printf("Error getting relative path: %v", err)
		return fullPath
	}
	return relPath
}
