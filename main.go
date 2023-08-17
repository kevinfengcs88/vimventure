package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/beevik/ntp"
	tea "github.com/charmbracelet/bubbletea"
)

var servers = [4]string{
	"129.6.15.28",
	"129.6.15.29",
	"129.6.15.30",
	"129.6.15.27",
}

func queryNTP(ip string, timeChan chan<- time.Time, ipChan chan<- string, quit chan bool) {
	currentTime, err := ntp.Time(ip)
	if err != nil {
		fmt.Println("Some error")
		return
	}
	select {
	case <-quit:
		return
	case timeChan <- currentTime:
		quit <- true
	case ipChan <- ip:
		quit <- true
	}
}

func NTPQuerier() (time.Time, string) {
	timeChan := make(chan time.Time)
	ipChan := make(chan string)
	quit := make(chan bool)

	for _, val := range servers {
		go queryNTP(val, timeChan, ipChan, quit)
	}

	firstResponse := <-timeChan
	firstResponseIP := <-ipChan

	return firstResponse, firstResponseIP
}

type editorFinishedMsg struct{ err error }

func openEditor(filename string, path string, start time.Time) tea.Cmd {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory:", err)
		return nil
	}
	targetFilePath := filepath.Join(currentDir, path, filename)
	c := exec.Command(editor, targetFilePath) //nolint:gosec
	return tea.ExecProcess(c, func(err error) tea.Msg {
		firstResponse, _ := NTPQuerier()
		d := firstResponse.Sub(start)
		score := calculateScore(d)
		fmt.Println("Your score was", score)
		fmt.Println("dummy print")
		fmt.Println("dummy print")
		return editorFinishedMsg{err}
	})
}

func timeBenchmark(d time.Duration) float64 {
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

func accuracyBenchmark(filename string, path string) int {
	currentDir, err := os.Getwd()

	targetFilePath := filepath.Join(currentDir, path, filename)

	file, err := os.Open(targetFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	solutionStart := 2
	solutionEnd := 11
	solutionText := []string{}

	contentStart := 13
	contentEnd := 22
	contentText := []string{}

	scanner := bufio.NewScanner(file)

	for lineNum := 1; lineNum < solutionStart; lineNum++ {
		if !scanner.Scan() {
			fmt.Println("Start line number is beyond the end of the file.")
			return 0
		}
	}

	for lineNum := solutionStart; lineNum <= solutionEnd; lineNum++ {
		if !scanner.Scan() {
			fmt.Println("End line number is beyond the end of the file.")
			return 0
		}
		// integrate REST API here
		solutionText = append(solutionText, scanner.Text())
	}

	for lineNum := solutionEnd + 1; lineNum < contentStart; lineNum++ {
		if !scanner.Scan() {
			fmt.Println("Start line number is beyond the end of the file.")
			return 0
		}
	}

	for lineNum := contentStart; lineNum <= contentEnd; lineNum++ {
		if !scanner.Scan() {
			fmt.Println("End line number is beyond the end of the file.")
			return 0
		}
		contentText = append(contentText, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error scanning file:", err)
	}

	fmt.Println("THIS IS THE SOLUTION")
	fmt.Println(solutionText)
	fmt.Println("THIS IS THE CONTENT")
	fmt.Println(contentText)

	accuracy := 0
	for i := 0; i < 10; i++ {
		if solutionText[i] == contentText[i] {
			accuracy += 10
		}
	}
	return accuracy
}

func calculateScore(d time.Duration) int {
	accuracy := accuracyBenchmark("example1.txt", "challenges")
	// let's assume the file is modified correctly, and the player gets a score of 100
	// should also check for cheating
	// compare the top half to a solution set
	// then just compare the bottom half to the top half once confirmed that the solution was not modified
	fmt.Println(d)
	return int(float64(accuracy) * timeBenchmark(d))
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
			firstResponse, _ := NTPQuerier()
			return m, openEditor("example1.txt", "challenges", firstResponse)
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
	return "Press 'e' to play the Vim challenge! Change the CONTENT to match the SOLUTION.\nPress 'q' to quit.\n"
}

func main() {
	m := model{}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
