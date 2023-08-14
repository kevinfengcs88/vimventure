package main

import (
	"fmt"
	"os"
	"os/exec"
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
		fmt.Println("this goroutine just quit and its IP was", ip)
		return
	case timeChan <- currentTime:
		fmt.Println("sent the currentTime to ch")
		quit <- true
	case ipChan <- ip:
		fmt.Println("sent the ip to ipChan")
		quit <- true
	}
}

type editorFinishedMsg struct{ err error }

func openEditor(filename string, start time.Time) tea.Cmd {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}
	c := exec.Command(editor, filename) //nolint:gosec
	return tea.ExecProcess(c, func(err error) tea.Msg {

		timeChan := make(chan time.Time)
		ipChan := make(chan string)
		quit := make(chan bool)

		for _, val := range servers {
			go queryNTP(val, timeChan, ipChan, quit)
		}

		firstResponse := <-timeChan
		firstResponseIP := <-ipChan

		fmt.Println("THIS IS END TIME - First response is", firstResponse, "and it came from", firstResponseIP)
		fmt.Println("THIS IS END TIME - First response is", firstResponse, "and it came from", firstResponseIP)
		fmt.Println("THIS IS END TIME - First response is", firstResponse, "and it came from", firstResponseIP)
		fmt.Println("THIS IS END TIME - First response is", firstResponse, "and it came from", firstResponseIP)
		fmt.Println("THIS IS END TIME - First response is", firstResponse, "and it came from", firstResponseIP)
		fmt.Println("THIS IS END TIME - First response is", firstResponse, "and it came from", firstResponseIP)

		d := firstResponse.Sub(start)
		score := calculateScore("test.txt", d)
		fmt.Println(score)
		fmt.Println("WOOOOOOOOOOOOOOOOOOOOOOOOOOO")
		fmt.Println(score)
		fmt.Println(score)
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
			timeChan := make(chan time.Time)
			ipChan := make(chan string)
			quit := make(chan bool)

			for _, val := range servers {
				go queryNTP(val, timeChan, ipChan, quit)
			}

			firstResponse := <-timeChan
			firstResponseIP := <-ipChan

			fmt.Println("First response is", firstResponse, "and it came from", firstResponseIP)
			fmt.Println("First response is", firstResponse, "and it came from", firstResponseIP)
			fmt.Println("First response is", firstResponse, "and it came from", firstResponseIP)
			fmt.Println("First response is", firstResponse, "and it came from", firstResponseIP)

			return m, openEditor("test.txt", firstResponse)
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
