package observer

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type SoulManager struct {
	Path string
}

func NewSoulManager(path string) *SoulManager {
	return &SoulManager{Path: path}
}

func (s *SoulManager) AppendEntry(speaker, target, content string) error {
	ts := time.Now().Format("15:04:05")
	entry := fmt.Sprintf("\n\n## [%s -> %s] (%s)\n%s\n", speaker, target, ts, content)
	
	f, err := os.OpenFile(s.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(entry)
	return err
}

func (s *SoulManager) GetRecentHistory(lineCount int) (string, error) {
	content, err := os.ReadFile(s.Path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}

	lines := strings.Split(string(content), "\n")
	if len(lines) > lineCount {
		lines = lines[len(lines)-lineCount:]
	}
	return strings.Join(lines, "\n"), nil
}

func (s *SoulManager) Initialize() error {
	if _, err := os.Stat(s.Path); os.IsNotExist(err) {
		return os.WriteFile(s.Path, []byte("# Project-Atra: Semantic Drift Log\n"), 0644)
	}
	return nil
}
