package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"golang.org/x/term"
)

func main() {
	var err error

	var conn *net.UDPConn

	// "udp",
	// "10.100.2.87:4444",
	addr := &net.UDPAddr{
		// IP:   net.ParseIP("10.100.2.87"),
		// Port: 13000,
		IP:   net.ParseIP("10.100.7.28"),
		Port: 12000,
		// IP:   net.ParseIP("10.100.2.87"),
		// Port: 4444,
	}

	if conn, err = net.DialUDP("udp", nil, addr); err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}

	defer errors.DeferredCloser(&err, conn)

	d := &Display{}

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}

	defer term.Restore(int(os.Stdin.Fd()), oldState)

	// os.Stdout.WriteString("\033[?12l")
	// os.Stdout.WriteString("\033[?25h")
	os.Stdout.WriteString("\033[?1003h")
	os.Stdout.WriteString("\033[1;1'z") // mouse reporting

	erase := false

	go func() {
		var b bytes.Buffer
		br := bufio.NewReader(os.Stdin)

		for {
			var nextChar byte
			nextChar, err = br.ReadByte()
			if err != nil {
				panic(err)
			}

			if nextChar != '\033' {
				continue
			}

			var next [100]byte
			var n int
			n, err = br.Read(next[:])
			b1 := append(b.AvailableBuffer(), append([]byte{'\033'}, next[:n]...)...)
			b.Write(b1)

			els := bytes.Split(b.Bytes(), []byte{'\x1b'})

		ELS:
			for _, e := range els {
				b.Read(make([]byte, len(e)+1))
				if len(e) < 4 {
					continue
				}

				switch e[2] {
				case byte('e'):
				case byte(' '): // click down
					erase = true
					continue ELS
				case byte('#'): // click up
					erase = false
					continue ELS
				default:
				}

				// payload := e[2:]
				col := int(e[3]) - 32
				row := int(e[4]) - 32

				if row == 1 && col != 1 {
					continue
				}

				loc := (Columns * (row - 1)) + col - 1
				d[loc] = !erase
				// fmt.Printf("full: %s coord: %d,%d\n\r", e, col, row)
			}
		}
	}()

	for {
		_, err = (*DisplayPrinterOSCLED)(d).WriteTo(conn)

		if err != nil {
			if errors.Unwrap(err).Error() == "no buffer space available" {
				err = nil
				continue
			}

			fmt.Println(err)
			log.Fatal(err)
		}

		_, err = (*DisplayPrinterTerminal)(d).WriteTo(os.Stdout)

		if err != nil {
			if errors.Unwrap(err).Error() == "no buffer space available" {
				err = nil
				continue
			}

			fmt.Println(err)
			log.Fatal(err)
		}

		os.Stdout.WriteString("\033[2;1'z")
		time.Sleep(1e7)
	}
}
