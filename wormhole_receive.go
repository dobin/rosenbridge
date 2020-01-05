// Mostly based on:
// https://github.com/psanford/wormhole-william/blob/master/cmd/recv.go

package main

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/psanford/wormhole-william/wormhole"
)

func wormholeConnect(code string) (*wormhole.IncomingMessage, error) {
	//var c wormhole.Client
	var serverAddress string = settingsBridge.ServerAddress()
	c := wormhole.Client{
		RendezvousURL: serverAddress,
	}

	ctx := context.Background()
	msg, err := c.Receive(ctx, code)
	if err != nil {
		//log.Fatal("Receive error: ", err)
		return nil, err
	}

	return msg, err
}

// TODO handle this differently
func wormholeTransferText(msg *wormhole.IncomingMessage, jobtotal *int, jobdone *int, feedbackstr *string) {
	body, err := ioutil.ReadAll(msg)
	if err != nil {
		*feedbackstr = fmt.Sprintf("Error: %s", err.Error())
		return
	}

	fmt.Println(string(body))
}

func wormholeTransferFile(msg *wormhole.IncomingMessage, jobtotal *int, jobdone *int, feedbackstr *string) {
	fmt.Printf("Receiving file (%s) into: %s\n", formatBytes(msg.TransferBytes), msg.Name)
	/*
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("ok? (y/N):")

		line, err := reader.ReadString('\n')
		if err != nil {
			errf("Error reading from stdin: %s\n", err)
		}
		line = strings.TrimSpace(line)
		if line == "y" {
			acceptFile = true
		}
	*/

	var acceptFile = true

	if !acceptFile {
		msg.Reject()
		*feedbackstr = fmt.Sprintf("transfer rejected")
		return
	} else {
		wd := getWorkingDirectory()

		f, err := ioutil.TempFile(wd, fmt.Sprintf("%s.tmp", msg.Name))
		if err != nil {
			*feedbackstr = fmt.Sprintf("Failed to create tempfile: %s", err)
			return
		}

		buffer := make([]byte, BufferSize)

		for {
			bytesread, err := msg.Read(buffer)
			time.Sleep(1 * time.Second)

			if err != nil {
				if err != io.EOF {
					fmt.Println(err)
				}

				break
			}

			*jobdone += bytesread
			f.Write(buffer)
		}

		tmpName := f.Name()
		f.Close()

		err = os.Rename(tmpName, msg.Name)
		if err != nil {
			*feedbackstr = fmt.Sprintf("Rename %s to %s failed: %s", tmpName, msg.Name, err)
			return
		}

		*feedbackstr = "Done"
	}
}

func wormholeTransferDirectory(msg *wormhole.IncomingMessage, jobtotal *int, jobdone *int, feedbackstr *string) {
	wd, err := os.Getwd()
	if err != nil {
		*feedbackstr = fmt.Sprintf("Failed to get working directory: %s", err)
		return
	}

	dirName := msg.Name
	dirName, err = filepath.Abs(dirName)
	if err != nil {
		*feedbackstr = fmt.Sprintf("Failed to get abs directory: %s", err)
		return
	}

	if filepath.Dir(dirName) != wd {
		*feedbackstr = fmt.Sprintf("Bad Directory name %s", msg.Name)
		return
	}

	if _, err := os.Stat(dirName); err == nil {
		*feedbackstr = fmt.Sprintf("Error refusing to overwrite existing '%s'", msg.Name)
		return
	} else if !os.IsNotExist(err) {
		*feedbackstr = fmt.Sprintf("Error stat'ing existing '%s'\n", msg.Name)
		return
	} else {
		//fmt.Printf("Receiving directory (%s) into: %s\n", formatBytes(msg.TransferBytes), msg.Name)
		//fmt.Printf("%d files, %s (uncompressed)\n", msg.FileCount, formatBytes(msg.UncompressedBytes))

		/*
			reader := bufio.NewReader(os.Stdin)

			fmt.Print("ok? (y/N):")

			line, err := reader.ReadString('\n')
			if err != nil {
				errf("Error reading from stdin: %s\n", err)
			}
			line = strings.TrimSpace(line)
			if line == "y" {
				acceptDir = true
			}*/

		var acceptDir bool = true
		if !acceptDir {
			msg.Reject()
			*feedbackstr = fmt.Sprintf("transfer rejected")
			return
		} else {
			err = os.Mkdir(msg.Name, 0777)
			if err != nil {
				*feedbackstr = fmt.Sprintf("Mkdir error for %s: %s\n", msg.Name, err)
				return
			}

			tmpFile, err := ioutil.TempFile(wd, fmt.Sprintf("%s.zip.tmp", msg.Name))
			if err != nil {
				*feedbackstr = fmt.Sprintf("Failed to create tempfile: %s", err)
				return
			}

			defer tmpFile.Close()
			defer os.Remove(tmpFile.Name())

			proxyReader := pbProxyReader(msg, msg.TransferBytes)

			n, err := io.Copy(tmpFile, proxyReader)
			if err != nil {
				os.Remove(tmpFile.Name())
				*feedbackstr = fmt.Sprintf("Receive file error: %s", err)
				return
			}

			tmpFile.Seek(0, io.SeekStart)
			zr, err := zip.NewReader(tmpFile, int64(n))
			if err != nil {
				*feedbackstr = fmt.Sprintf("Read zip error: %s", err)
				return
			}

			for _, zf := range zr.File {
				p, err := filepath.Abs(filepath.Join(dirName, zf.Name))
				if err != nil {
					*feedbackstr = fmt.Sprintf("Failes to calculate file path ABS: %s", err)
					return
				}

				if !strings.HasPrefix(p, dirName) {
					*feedbackstr = fmt.Sprintf("Dangerous filename detected: %s", zf.Name)
					return
				}

				rc, err := zf.Open()
				if err != nil {
					*feedbackstr = fmt.Sprintf("Failed to open file in zip: %s %s", zf.Name, err)
					return
				}

				dir := filepath.Dir(p)
				err = os.MkdirAll(dir, 0777)
				if err != nil {
					*feedbackstr = fmt.Sprintf("Failed to mkdirall %s: %s", dir, err)
					return
				}

				f, err := os.Create(p)
				if err != nil {
					*feedbackstr = fmt.Sprintf("Failed to open %s: %s", p, err)
					return
				}

				_, err = io.Copy(f, rc)
				if err != nil {
					*feedbackstr = fmt.Sprintf("Failed to write to %s: %s", p, err)
					return
				}

				err = f.Close()
				if err != nil {
					*feedbackstr = fmt.Sprintf("Error closing %s: %s", p, err)
					return
				}

				rc.Close()
			}

			proxyReader.Close()

			*feedbackstr = "Done"
		}
	}
}

func errf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg, args...)
	if !strings.HasSuffix("\n", msg) {
		fmt.Fprint(os.Stderr, "\n")
	}
}

func formatBytes(b int) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}

type proxyReadCloser struct {
	*pb.Reader
	bar *pb.ProgressBar
}

func (p *proxyReadCloser) Close() error {
	p.bar.Finish()
	return nil
}

var hideProgressBar = false

func pbProxyReader(r io.Reader, size int) io.ReadCloser {
	if hideProgressBar {
		return ioutil.NopCloser(r)
	} else {
		progressBar := pb.Full.Start(size)
		proxyReader := progressBar.NewProxyReader(r)
		return &proxyReadCloser{
			Reader: proxyReader,
			bar:    progressBar,
		}
	}
}
