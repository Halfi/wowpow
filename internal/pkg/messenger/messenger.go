package messenger

import (
	"bufio"
	"fmt"
	"io/fs"
	"math/rand"
	"strings"
)

type Messenger struct {
	messages []string
}

// New constructor of messenger
// Gets messages to memory from file system
func New(f fs.FS) (*Messenger, error) {
	file, err := f.Open("quotes.txt")
	if err != nil {
		return nil, fmt.Errorf("open quotes file error: %w", err)
	}

	reader := bufio.NewReader(file)
	messages := make([]string, 0)
	for {
		quote, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		messages = append(messages, strings.Trim(quote, "\n"))
	}

	// we don't need to catch read error
	return &Messenger{messages: messages}, nil //nolint:nilerr
}

// GetMessage returns message (quote)
func (m *Messenger) GetMessage() string {
	// this is not secure function, math/rand enough good
	return m.messages[rand.Intn(len(m.messages)-1)] //nolint:gosec
}
