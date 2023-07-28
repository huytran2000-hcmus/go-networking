package ch04

import "io"

func proxy(from io.Reader, to io.Writer) error {
	fromWriter, writerOK := from.(io.Writer)
	toReader, readerOK := to.(io.Reader)

	if writerOK && readerOK {
		go func() {
			io.Copy(fromWriter, toReader)
		}()
	}

	_, err := io.Copy(to, from)

	return err
}
