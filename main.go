package main

import (
	// "bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/beevik/ntp"
	tea "github.com/charmbracelet/bubbletea"
)

const NTP_SERVER = "129.6.15.28"

type editorFinishedMsg struct{ err error }

func openEditor(filename string, start time.Time) tea.Cmd {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}
	c := exec.Command(editor, filename) //nolint:gosec
	return tea.ExecProcess(c, func(err error) tea.Msg {
		stop, err := ntp.Time(NTP_SERVER)
		if err != nil {
			log.Fatal(err)
		}
		d := stop.Sub(start)
		score := calculateScore("test.txt", d)
		fmt.Println(score)
		fmt.Println(score)
		fmt.Println(score)
		fmt.Println("WOOOOOOOOOOOOOOOOOOOOOOOOOOO")
		return editorFinishedMsg{err}
	})
}

func scaleDuration(d time.Duration) float64 {
	const (
		threshold        = 60 * time.Second
		maxValue         = 1.0
		minValue         = 0.0
		interval         = 5 * time.Second
		scalePerInterval = (maxValue - minValue) / float64(threshold/interval)
	)

	if d < interval {
		return maxValue
	}

	if d > threshold {
		return minValue
	}

	intervalsPassed := float64(d / interval)
	scaledValue := maxValue - (intervalsPassed * scalePerInterval)

	return scaledValue
}

func calculateScore(filename string, d time.Duration) int {
	// file, err := os.Open(filename)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer file.Close()
	//
	// scanner := bufio.NewScanner(file)
	//
	// for scanner.Scan() {
	// 	line := scanner.Text()
	// 	fmt.Println(line)
	// }
	//
	// if err := scanner.Err(); err != nil {
	// 	log.Fatal(err)
	// }

	// let's assume the file is modified correctly, and the player gets a score of 100
	// should also check for cheating
	// compare the top half to a solution set
	// then just compare the bottom half to the top half once confirmed that the solution was not modified
	fmt.Println(d)
	return int(100 * scaleDuration(d))
}

type model struct {
	err error
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "e":
			start, err := ntp.Time(NTP_SERVER)
			if err != nil {
				log.Fatal(err)
			}
			return m, openEditor("test.txt", start)
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case editorFinishedMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.err != nil {
		return "Error: " + m.err.Error() + "\n"
	}
	return "Press 'e' to play the Vim challenge! Change the content below the demarcation line to look like that which resides above the line.\nPress 'q' to quit.\n"
}

func main() {
	m := model{}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
