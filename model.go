package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
	"time"
)

type errMsg error

type gameState int

const (
	gameStatePlaying gameState = iota
	gameStateWin
	gameStateLose
)

type model struct {
	wordList       *WordList
	textInput      textinput.Model
	enteredAnswers []string
	answerDay      int64
	correctAnswer  string
	err            error
	shareMessage   string
	maxTries       int
	gameState      gameState
}

func initialModel(wl *WordList, answerDay int64) model {
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 5
	ti.Width = 20

	var correctAnswer string

	if answerDay == -1 {
		answerDay, correctAnswer = wl.GetWordByDate(time.Now())
	} else {
		correctAnswer = wl.GetWord(answerDay)
	}

	return model{
		enteredAnswers: []string{},
		answerDay:      answerDay,
		correctAnswer:  correctAnswer,
		wordList:       wl,
		textInput:      ti,
		err:            nil,
		maxTries:       6,
		gameState:      gameStatePlaying,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, tea.EnterAltScreen)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			answer := strings.ToLower(strings.TrimSpace(m.textInput.Value()))
			if len(answer) != 5 {
				return m, nil
			}
			if !m.wordList.CheckWord(answer) {
				// not a valid word
				return m, nil
			}

			m.textInput.SetValue("")
			m.enteredAnswers = append(m.enteredAnswers, answer)
			m.shareMessage = m.RenderShareMessage()
			if answer == m.correctAnswer {
				m.gameState = gameStateWin
			}
			if len(m.enteredAnswers) == m.maxTries {
				m.gameState = gameStateLose
			}
			return m, nil
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	result := ""
	switch m.gameState {
	case gameStateWin:
		result = fmt.Sprintf("%sWell Done!\n\n%s", m.RenderAnswers(), m.RenderShareMessage())
	case gameStateLose:
		result = fmt.Sprintf("%sToday's word was:\n\n%s\n\n%s",
			m.RenderAnswers(),
			strings.ToUpper(m.correctAnswer),
			m.RenderShareMessage(),
		)
	case gameStatePlaying:
		result = fmt.Sprintf(
			"%s\nWhat’s your guess?\n\n%s",
			m.RenderAnswers(),
			m.textInput.View(),
		)
	}
	return lipgloss.NewStyle().
		Width(40).
		Height(16).
		PaddingLeft(1).
		PaddingRight(1).
		Render(result + "\n" + subtle("esc to quit") + "\n\n")
}

func subtle(msg string) string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("202")).Render(msg)
}

func getStyle(foregroundColor, backgroundColor string) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(foregroundColor)).Background(lipgloss.Color(backgroundColor))
}

var (
	styleCorrect    = getStyle("#CED2D4", "#437C3D")
	styleWrongPlace = getStyle("#CED2D4", "#A6902E")
	styleIncorrect  = getStyle("#CED2D4", "#2C2C2D")
)

var (
	runeIncorrect     = `⬛`
	runeWrongLocation = `🟨`
	runeCorrect       = `🟩`
)

func (m model) RenderAnswers() string {
	result := ""
	for row := 0; row < m.maxTries; row++ {
		if len(m.enteredAnswers) <= row {
			result += "\n"
			continue
		}

		answer := m.enteredAnswers[row]
		for i := range answer {
			if answer[i] == m.correctAnswer[i] {
				result += styleCorrect.Render(strings.ToUpper(" " + string(answer[i]) + " "))
			} else {
				if strings.Contains(m.correctAnswer, string(answer[i])) {
					result += styleWrongPlace.Render(strings.ToUpper(" " + string(answer[i]) + " "))
				} else {
					result += styleIncorrect.Render(strings.ToUpper(" " + string(answer[i]) + " "))
				}
			}
		}
		result += "\n"
	}
	return result
}

func (m model) RenderShareMessage() string {
	result := fmt.Sprintf("WordTerm %d %d/%d\n", m.answerDay, len(m.enteredAnswers), 6)
	for _, answer := range m.enteredAnswers {
		for i := range answer {
			if answer[i] == m.correctAnswer[i] {
				result += runeCorrect
			} else {
				if strings.Contains(m.correctAnswer, string(answer[i])) {
					result += runeWrongLocation
				} else {
					result += runeIncorrect
				}
			}
		}
		result += "\n"
	}
	return result
}
