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

	return &Messenger{messages: messages}, nil
}

func (m *Messenger) GetMessage() string {
	return m.messages[rand.Intn(len(m.messages)-1)]
}
