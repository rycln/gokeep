package bin

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rycln/gokeep/internal/client/tui/shared/input"
	"github.com/rycln/gokeep/internal/client/tui/shared/messages"
	"github.com/rycln/gokeep/internal/shared/models"
)

const (
	nameField    = "Название"
	infoField    = "Информация"
	binPathField = "Путь к файлу"
)

type Model struct {
	input.Form
}

func InitialModel() Model {
	m := Model{
		input.Form{
			Inputs: make([]textinput.Model, 3),
		},
	}

	for i := range m.Inputs {
		t := textinput.New()
		t.Width = 30
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = nameField
			t.Focus()
		case 1:
			t.Placeholder = infoField
		case 2:
			t.Placeholder = binPathField
			t.CharLimit = 100
		}

		m.Inputs[i] = t
	}

	return m
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			return m, m.send()
		}
	}

	updated, cmd := m.Form.Update(msg)
	if inputForm, ok := updated.(input.Form); ok {
		m.Form = inputForm
	}
	return m, cmd
}

func (m Model) send() tea.Cmd {
	return func() tea.Msg {
		info := &models.ItemInfo{
			ItemType: models.TypeBinary,
			Name:     m.Inputs[0].Value(),
			Metadata: m.Inputs[1].Value(),
		}

		file, err := os.ReadFile(m.Inputs[2].Value())
		if err != nil {
			return messages.ErrMsg{Err: fmt.Errorf("ошибка чтения файла: %v", err)}
		}

		bin := &BinFile{
			Data: file,
		}

		content, err := json.Marshal(bin)
		if err != nil {
			return messages.ErrMsg{Err: err}
		}

		m.Inputs[0].Reset()
		m.Inputs[1].Reset()
		m.Inputs[2].Reset()

		return messages.ItemMsg{
			Info:    info,
			Content: content,
		}
	}
}

func (m *Model) SetStartData(info *models.ItemInfo, content []byte) error {
	var binary BinFile

	err := json.Unmarshal(content, &binary)
	if err != nil {
		return err
	}

	m.Inputs[0].SetValue(info.Name)
	m.Inputs[1].SetValue(info.Metadata)

	return nil
}
