// Mostly based on:
// https://github.com/psanford/wormhole-william/blob/master/cmd/send.go

package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/psanford/wormhole-william/wormhole"
)

var (
	codeLen  int
	codeFlag string
	verify   bool = false
)

/*
func sendCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "send [WHAT]",
		Short: "Send a text message, file, or directory...",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				sendText()
				return
			} else if len(args) > 1 {
				bail("Too many arguments")
			}

			stat, err := os.Stat(args[0])
			if err != nil {
				bail("Failed to read %s: %s", args[0], err)
			}

			if stat.IsDir() {
				sendDir(args[0])
				return
			} else {
				sendFile(args[0])
				return
			}
		},
	}

	cmd.Flags().BoolVarP(&verify, "verify", "v", false, "display verification string (and wait for approval)")
	cmd.Flags().IntVarP(&codeLen, "code-length", "c", 0, "length of code (in bytes/words)")
	cmd.Flags().StringVar(&codeFlag, "code", "", "human-generated code phrase")
	cmd.Flags().BoolVar(&hideProgressBar, "hide-progress", false, "supress progress-bar display")

	return &cmd
}
*/

func newClient() wormhole.Client {
	var serverAddress string = settingsBridge.ServerAddress()
	c := wormhole.Client{
		RendezvousURL: serverAddress,
	}

	if verify {
		c.VerifierOk = func(code string) bool {
			reader := bufio.NewReader(os.Stdin)
			fmt.Printf("Verifier %s. ok? (yes/no): ", code)

			yn, _ := reader.ReadString('\n')
			yn = strings.TrimSpace(yn)

			return yn == "yes"
		}
	}

	return c
}

func printInstructions(code string) {
	mwCmd := "wormhole receive"
	wwCmd := "wormhole-william recv"

	if verify {
		mwCmd = mwCmd + " --verify"
		wwCmd = wwCmd + " --verify"
	}

	fmt.Printf("On the other computer, please run: %s (or %s)\n", mwCmd, wwCmd)
	fmt.Printf("Wormhole code is: %s\n", code)
}

func sendFile(filename string, jobdone *int64, feedbackstr *string) (string, chan wormhole.SendResult, error) {
	f, err := os.Open(filename)
	if err != nil {
		//*feedbackstr = fmt.Sprintf("Failed to open %s: %s", filename, err)
		return "", nil, err
	}

	c := newClient()

	ctx := context.Background()

	//var bar *pb.ProgressBar

	args := []wormhole.SendOption{
		wormhole.WithCode(codeFlag),
	}

	if !hideProgressBar {
		args = append(args, wormhole.WithProgress(func(sentBytes int64, totalBytes int64) {
			fmt.Printf("Sent: %d", sentBytes)
			*jobdone = sentBytes
			/*if bar == nil {
				bar = pb.Full.Start64(totalBytes)
				bar.Set(pb.Bytes, true)
			}
			bar.SetCurrent(sentBytes)

			if sentBytes == totalBytes {
				bar.Finish()
			}*/
		}))
	}

	code, status, err := c.SendFile(ctx, filepath.Base(filename), f, args...)
	if err != nil {
		//*feedbackstr = fmt.Sprintf("Error sending message: %s", err)
		return "", nil, err
	}

	return code, status, err
}

func sendDir(dirpath string, jobdone *int64, feedbackstr *string) (string, chan wormhole.SendResult, error) {
	dirpath = strings.TrimSuffix(dirpath, "/")

	stat, err := os.Stat(dirpath)
	if err != nil {
		log.Fatal(err)
	}

	if !stat.IsDir() {
		//log.Fatalf("%s is not a directory", dirpath)
		*feedbackstr = fmt.Sprintf("%s is not a directory", dirpath)
		return "", nil, err
	}

	prefix, dirname := filepath.Split(dirpath)

	var entries []wormhole.DirectoryEntry

	filepath.Walk(dirpath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if !info.Mode().IsRegular() {
			return nil
		}

		relPath := strings.TrimPrefix(path, prefix)

		entries = append(entries, wormhole.DirectoryEntry{
			Path: relPath,
			Mode: info.Mode(),
			Reader: func() (io.ReadCloser, error) {
				return os.Open(path)
			},
		})

		return nil
	})

	c := newClient()

	ctx := context.Background()
	code, status, err := c.SendDirectory(ctx, dirname, entries, wormhole.WithCode(codeFlag))
	if err != nil {
		// *feedbackstr = fmt.Sprintf
		// log.Fatal(err)
		return "", nil, err
	}

	return code, status, err

	/*
		printInstructions(code)

		s := <-status

		if s.OK {
			fmt.Println("directory sent")
		} else {
			bail("Send error: %s", s.Error)
		}
	*/
}

func sendText() {
	c := newClient()

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Text to send: ")
	msg, _ := reader.ReadString('\n')

	msg = strings.TrimSpace(msg)

	ctx := context.Background()

	code, status, err := c.SendText(ctx, msg, wormhole.WithCode(codeFlag))
	if err != nil {
		log.Fatal(err)
	}

	printInstructions(code)

	s := <-status

	if s.Error != nil {
		log.Fatalf("Send error: %s", s.Error)
	} else if s.OK {
		fmt.Println("text message sent")
	} else {
		log.Fatalf("Hmm not ok but also not error")
	}
}
