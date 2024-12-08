package terminal

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type TerminalPrompter struct {
	reader io.Reader
	writer io.Writer
}

func NewTerminalPrompter(reader io.Reader, writer io.Writer) *TerminalPrompter {
	return &TerminalPrompter{
		reader: reader,
		writer: writer,
	}
}

func (tp *TerminalPrompter) Prompt(responseHandler func(string) bool, msg string) string {
	fmt.Fprint(tp.writer, msg)

	reader := bufio.NewReader(tp.reader)
	text, _ := reader.ReadString('\n')
	for !responseHandler(text) {
		fmt.Fprint(tp.writer, "Invalid response. Try again.\n")
		fmt.Fprint(tp.writer, msg)
		text, _ = reader.ReadString('\n')
	}

	return strings.TrimSpace(text)
}
