package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
)

type editorFinishedMsg struct{ err error }

func openEditor(filename string, m model) tea.Cmd {
	m.stopwatch.Start()
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}
	c := exec.Command(editor, filename) //nolint:gosec
	return tea.ExecProcess(c, func(err error) tea.Msg {
		// calculateScore(filename)
		fmt.Println(m.stopwatch.Elapsed())
		return editorFinishedMsg{err}
	})
}

func calculateScore(filename string) int {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return 3
}

type model struct {
	err       error
	stopwatch stopwatch.Model
}

func (m model) Init() tea.Cmd {
	return tea.Sequence(m.stopwatch.Init(), m.stopwatch.Stop())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "e":
			return m, tea.Sequence(m.stopwatch.Start(), openEditor("test.txt", m))
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case editorFinishedMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, tea.Quit
		}
	}
	// return m, nil
	var cmd tea.Cmd
	m.stopwatch, cmd = m.stopwatch.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.err != nil {
		return "Error: " + m.err.Error() + "\n"
	}
	s := m.stopwatch.View() + "\n"
	return s + "Press 'e' to play the Vim challenge! Change the content below the demarcation line to look like that which resides above the line.\nPress 'q' to quit.\n"
}

func main() {
	m := model{
		stopwatch: stopwatch.NewWithInterval(time.Millisecond),
	}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
