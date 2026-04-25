package lox

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/islml/neptune/internal/scanner"
)

const (
	exitCodeUsage     = 64
	exitCodeDataError = 65
	exitCodeSoftware  = 70
)

var errScanFailed = errors.New("scan failed")

type Lox struct {
	in       io.Reader
	out      io.Writer
	errOut   io.Writer
	hadError bool
}

func New(in io.Reader, out io.Writer, errOut io.Writer) *Lox {
	return &Lox{
		in:     in,
		out:    out,
		errOut: errOut,
	}
}

func (l *Lox) Execute(args []string) (int, error) {
	switch {
	case len(args) > 1:
		return exitCodeUsage, fmt.Errorf("usage: neptune [script]")
	case len(args) == 1:
		err := l.runFile(args[0])
		switch {
		case err == nil:
			return 0, nil
		case errors.Is(err, errScanFailed):
			return exitCodeDataError, nil
		default:
			return exitCodeSoftware, err
		}
	default:
		if err := l.runPrompt(); err != nil {
			return exitCodeSoftware, err
		}
		return 0, nil
	}
}

func (l *Lox) runFile(path string) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("unable to read file %q: %w", path, err)
	}

	if err := l.run(string(bytes)); err != nil {
		return err
	}

	if l.hadError {
		return errScanFailed
	}

	return nil
}

func (l *Lox) runPrompt() error {
	input := bufio.NewScanner(l.in)

	for {
		fmt.Fprint(l.out, "> ")
		if ok := input.Scan(); !ok {
			break
		}

		line := input.Text()
		if err := l.run(line); err != nil && !errors.Is(err, errScanFailed) {
			return err
		}

		l.hadError = false
	}

	if err := input.Err(); err != nil {
		return fmt.Errorf("failed reading input: %w", err)
	}

	return nil
}

func (l *Lox) run(source string) error {
	s := scanner.New(source)
	tokens, scanErrors := s.ScanTokens()

	if len(scanErrors) > 0 {
		for _, scanErr := range scanErrors {
			l.report(scanErr.Line, "", scanErr.Message)
		}
		return errScanFailed
	}

	for _, t := range tokens {
		fmt.Fprintln(l.out, t.String())
	}

	return nil
}

func (l *Lox) report(line int, where string, message string) {
	if where != "" {
		where = " " + where
	}

	fmt.Fprintf(l.errOut, "[line %d] Error%s: %s\n", line, where, message)
	l.hadError = true
}