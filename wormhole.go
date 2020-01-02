package main

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/cheggaaa/pb/v3"
	"github.com/psanford/wormhole-william/wormhole"
)

func guiDownload(code string) {
	msg, err := wormholeConnect(code)
	if err != nil {
		// log.Fatal(err)
		fmt.Printf("Could not connect, wrong code?")
		return
	}

	switch msg.Type {
	case wormhole.TransferText:
		wormholeTransferText(msg)
	case wormhole.TransferFile:
		wormholeTransferFile(msg)
	case wormhole.TransferDirectory:
		wormholeTransferDirectory(msg)
	}
}

func wormholeConnect(code string) (*wormhole.IncomingMessage, error) {
	var c wormhole.Client

	ctx := context.Background()
	msg, err := c.Receive(ctx, code)
	if err != nil {
		//log.Fatal("Receive error: ", err)
		return nil, err
	}

	return msg, err
}

func wormholeTransferFile(msg *wormhole.IncomingMessage) {
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
		bail("transfer rejected")
	} else {
		wd, err := os.Getwd()
		if err != nil {
			bail("Failed to get working directory: %s", err)
		}
		f, err := ioutil.TempFile(wd, fmt.Sprintf("%s.tmp", msg.Name))
		if err != nil {
			bail("Failed to create tempfile: %s", err)
		}

		proxyReader := pbProxyReader(msg, msg.TransferBytes)

		_, err = io.Copy(f, proxyReader)
		if err != nil {
			os.Remove(f.Name())
			bail("Receive file error: %s", err)
		}

		proxyReader.Close()

		tmpName := f.Name()
		f.Close()

		err = os.Rename(tmpName, msg.Name)
		if err != nil {
			bail("Rename %s to %s failed: %s", tmpName, msg.Name, err)
		}
	}
}

func wormholeTransferDirectory(msg *wormhole.IncomingMessage) {
	var acceptDir bool

	wd, err := os.Getwd()
	if err != nil {
		bail("Failed to get working directory: %s", err)
	}

	dirName := msg.Name
	dirName, err = filepath.Abs(dirName)
	if err != nil {
		bail("Failed to get abs directory: %s", err)
	}

	if filepath.Dir(dirName) != wd {
		bail("Bad Directory name %s", msg.Name)
	}

	if _, err := os.Stat(dirName); err == nil {
		errf("Error refusing to overwrite existing '%s'", msg.Name)
	} else if !os.IsNotExist(err) {
		errf("Error stat'ing existing '%s'\n", msg.Name)
	} else {
		fmt.Printf("Receiving directory (%s) into: %s\n", formatBytes(msg.TransferBytes), msg.Name)
		fmt.Printf("%d files, %s (uncompressed)\n", msg.FileCount, formatBytes(msg.UncompressedBytes))

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

		if !acceptDir {
			msg.Reject()
			bail("transfer rejected")
		} else {
			err = os.Mkdir(msg.Name, 0777)
			if err != nil {
				bail("Mkdir error for %s: %s\n", msg.Name, err)
			}

			tmpFile, err := ioutil.TempFile(wd, fmt.Sprintf("%s.zip.tmp", msg.Name))
			if err != nil {
				bail("Failed to create tempfile: %s", err)
			}

			defer tmpFile.Close()
			defer os.Remove(tmpFile.Name())

			proxyReader := pbProxyReader(msg, msg.TransferBytes)

			n, err := io.Copy(tmpFile, proxyReader)
			if err != nil {
				os.Remove(tmpFile.Name())
				bail("Receive file error: %s", err)
			}

			tmpFile.Seek(0, io.SeekStart)
			zr, err := zip.NewReader(tmpFile, int64(n))
			if err != nil {
				bail("Read zip error: %s", err)
			}

			for _, zf := range zr.File {
				p, err := filepath.Abs(filepath.Join(dirName, zf.Name))
				if err != nil {
					bail("Failes to calculate file path ABS: %s", err)
				}

				if !strings.HasPrefix(p, dirName) {
					bail("Dangerous filename detected: %s", zf.Name)
				}

				rc, err := zf.Open()
				if err != nil {
					bail("Failed to open file in zip: %s %s", zf.Name, err)
				}

				dir := filepath.Dir(p)
				err = os.MkdirAll(dir, 0777)
				if err != nil {
					bail("Failed to mkdirall %s: %s", dir, err)
				}

				f, err := os.Create(p)
				if err != nil {
					bail("Failed to open %s: %s", p, err)
				}

				_, err = io.Copy(f, rc)
				if err != nil {
					bail("Failed to write to %s: %s", p, err)
				}

				err = f.Close()
				if err != nil {
					bail("Error closing %s: %s", p, err)
				}

				rc.Close()
			}

			proxyReader.Close()

		}
	}
}

func wormholeTransferText(msg *wormhole.IncomingMessage) {
	body, err := ioutil.ReadAll(msg)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(body))
}

////

func sendText() {
	var c wormhole.Client

	msg := "Dillinger-entertainer"

	ctx := context.Background()

	code, status, err := c.SendText(ctx, msg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("On the other computer, please run: wormhole receive")
	fmt.Printf("Wormhole code is: %s\n", code)

	s := <-status

	if s.OK {
		fmt.Println("OK!")
	} else {
		log.Fatalf("Send error: %s", s.Error)
	}
}

func recvText(code string) {
	var c wormhole.Client

	ctx := context.Background()
	msg, err := c.Receive(ctx, code)
	if err != nil {
		log.Fatal(err)
	}

	if msg.Type != wormhole.TransferText {
		log.Fatalf("Expected a text message but got type %s", msg.Type)
	}

	msgBody, err := ioutil.ReadAll(msg)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("got message:")
	fmt.Println(msgBody)
}

func errf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg, args...)
	if !strings.HasSuffix("\n", msg) {
		fmt.Fprint(os.Stderr, "\n")
	}
}

func bail(msg string, args ...interface{}) {
	errf(msg, args...)
	os.Exit(1)
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
