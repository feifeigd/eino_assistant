package mem

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/cloudwego/eino/schema"
)

type (
	Conversation struct {
		mu sync.Mutex

		ID            string            `json:"id"`
		Messages      []*schema.Message `json:"messages"`
		filePath      string
		maxWindowSize int
	}

	// simple memory can stroe messages o each onversation
	SimpleMemory struct {
		mu            sync.Mutex
		dir           string
		maxWindowSize int
		conversations map[string]*Conversation
	}

	SimpleMemoryConfig struct {
		Dir           string
		MaxWindowSize int
	}
)

func GetDefaultMemory() *SimpleMemory {
	return NewSimpleMemory(SimpleMemoryConfig{
		Dir:           "data/memory",
		MaxWindowSize: 6,
	})
}

func NewSimpleMemory(config SimpleMemoryConfig) *SimpleMemory {
	if config.Dir == "" {
		config.Dir = "/tmp/eino/memory"
	}
	if err := os.MkdirAll(config.Dir, 0755); err != nil {
		return nil
	}

	return &SimpleMemory{
		dir:           config.Dir,
		maxWindowSize: config.MaxWindowSize,
		conversations: make(map[string]*Conversation),
	}
}

func (m *SimpleMemory) GetConveration(id string, createIfNotExist bool) *Conversation {
	m.mu.Lock()
	defer m.mu.Unlock()

	if con, ok := m.conversations[id]; ok {
		return con
	}

	filePath := filepath.Join(m.dir, id, ".jsonl")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {

		if createIfNotExist {
			if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
				return nil
			}
			if _, err := os.Create(filePath); err != nil {
				return nil
			}
		} else {
			return nil
		}
	}

	con := &Conversation{
		ID:            id,
		Messages:      make([]*schema.Message, 0),
		filePath:      filePath,
		maxWindowSize: m.maxWindowSize,
	}
	con.load()
	m.conversations[id] = con
	return con
}

func (m *SimpleMemory) ListConversations() []string {
	m.mu.Lock()
	defer m.mu.Unlock()

	files, err := os.ReadDir(m.dir)
	if err != nil {
		return nil
	}
	ids := make([]string, 0, len(files))
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		ids = append(ids, strings.TrimSuffix(file.Name(), "jsonl"))
	}
	return ids
}

func (m *SimpleMemory) DeleteConversation(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	filePath := filepath.Join(m.dir, id+".jsonl")
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	delete(m.conversations, id)
	return nil
}

func (c *Conversation) Append(msg *schema.Message) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Messages = append(c.Messages, msg)
	c.save(msg)
}

func (c *Conversation) GetFullMessages() []*schema.Message {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.Messages
}

func (c *Conversation) GetMessages() []*schema.Message {
	c.mu.Lock()
	defer c.mu.Unlock()

	if l := len(c.Messages); l > c.maxWindowSize {
		return c.Messages[l-c.maxWindowSize:]
	}
	return c.Messages
}

func (c *Conversation) load() error {
	reader, err := os.Open(c.filePath)
	if err != nil {
		return fmt.Errorf("failed to open conversation file: %v", err)
	}
	defer reader.Close()

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		var msg schema.Message
		if err := json.Unmarshal([]byte(line), &msg); err != nil {
			return fmt.Errorf("failed to unmarshal message: %v", err)
		}
		c.Messages = append(c.Messages, &msg)
	}

	return scanner.Err()
}

func (c *Conversation) save(msg *schema.Message) {
	str, err := json.Marshal(msg)
	if err != nil {
		return
	}
	f, err := os.OpenFile(c.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()

	f.Write(str)
	f.WriteString("\n")
}
