package term

import (
	"fmt"
	"io"
	"strings"

	"errors"

	"github.com/chzyer/readline"
)

var (
	ErrInputInterrupted = errors.New("input interrupted")
	ErrInputKilled      = errors.New("input killed")
	ErrReadlineInit     = errors.New("error initializing readline")
)

func NewInputArea() (string, error) {
	rl, err := readline.NewEx(&readline.Config{
		Prompt:                 ">>> ",
		HistoryFile:            "/dev/null",
		InterruptPrompt:        "Interrupted. Quitting...",
		EOFPrompt:              "Killed. Quitting...",
		AutoComplete:           nil,
		DisableAutoSaveHistory: true,
	})
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrReadlineInit, err)
	}
	defer rl.Close()

	var lines []string
	var emptyLines int

	for {
		line, err := rl.Readline()
		if err != nil {
			if err == readline.ErrInterrupt {
				return "", ErrInputInterrupted
			}
			if err == io.EOF {
				return "", ErrInputKilled
			}
			return "", fmt.Errorf("error reading input: %w", err)
		}

		if strings.TrimSpace(line) == "" {
			emptyLines++
		} else {
			emptyLines = 0
		}

		lines = append(lines, line)

		if emptyLines >= 2 {
			break
		}
	}

	return strings.Join(lines, "\n"), nil
}
