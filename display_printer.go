package main

import (
	"bufio"
	"fmt"
	"io"
	"time"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
)

type DisplayPrinterTerminal Display

func (d DisplayPrinterTerminal) WriteTo(w io.Writer) (n int64, err error) {
	bw := bufio.NewWriter(w)
	defer errors.DeferredFlusher(&err, bw)

	displayPixels := Display(d)
	bw.WriteString("\033[1;1H")
	bw.WriteString("\033[H")

	for row := 0; row < Rows; row++ {
		for column := 0; column < Columns; column++ {
			fmt.Fprintf(bw, "\033[%d;%dH", row+1, column+1)
			loc := (row * Columns) + column
			pixel := displayPixels[loc]

			if pixel {
				bw.WriteRune('*')
			} else {
				bw.WriteRune(' ')
			}
		}
	}

	fmt.Fprintf(bw, "\033[%d;%dH%s", Rows+1, 1, time.Now())
	// fmt.Fprintf(bw, "\033[%d;%dH", Rows+2, 1)

	for row := 0; row < Rows; row++ {
		for column := 0; column < Columns; column++ {
			loc := (row * Columns) + column
			pixel := displayPixels[loc]

			if !pixel {
				continue
			}

			// fmt.Fprintf(bw, "%d:%d (%d)", row+1, column+1, loc)
		}
	}

	return
}

type DisplayPrinterOSCLED Display

func writeString(v string, w *bufio.Writer) (n int64, err error) {
	var n2 int

	n2, err = w.Write([]byte(v))
	n += int64(n2)

	// if err != nil {
	// 	err = errors.Wrap(err)
	// 	return
	// }

	// l := len(v)

	// zeroCount := ((l + 4) & ^0x03) - l

	// var n1 int64
	// n1, err = writeNullBytes(zeroCount, w)
	// n += n1

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func writeDisplay(d *DisplayPrinterOSCLED, w *bufio.Writer) (n int64, err error) {
	var n1 int

	for _, v := range *d {
		r := '0'
		// b := byte(48) // ASCII 0

		if v {
			r = '1'
			// b = 49 // ASCII 1
		}

		// err = w.WriteByte(b)
		n1, err = w.WriteRune(r)
		n += int64(n1)

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	n1, err = w.Write([]byte("\x00\x00\x00\x00"))
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}
	// n1, err = writeNullBytes(PixelCount, w)
	// n += n1

	// if err != nil {
	// 	err = errors.Wrap(err)
	// 	return
	// }

	return
}

func writeNullBytes(l int, w *bufio.Writer) (n int64, err error) {
	zeroCount := ((l + 4) & ^0x03) - l

	for i := 0; i < zeroCount; i++ {
		err = w.WriteByte('\x00')
		n += 1

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (d DisplayPrinterOSCLED) WriteTo(w io.Writer) (n int64, err error) {
	bw := bufio.NewWriter(w)
	defer errors.DeferredFlusher(&err, bw)

	var n1 int64
	n1, err = writeString("/\x00\x00\x00", bw)
	n += n1

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	n1, err = writeString(",s\x00\x00", bw)
	n += n1

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	n1, err = writeDisplay(&d, bw)
	n += n1

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
