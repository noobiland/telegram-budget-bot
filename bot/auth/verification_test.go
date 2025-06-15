package auth

import (
	"context"
	"testing"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// GetUserName test data preparation
type getUserNameTestCase struct {
	name     string
	chatId   int64
	expected string
	found    bool
}

var getUserNameTestCases = []getUserNameTestCase{
	{"existing user", 1, "Alice", true},
	{"non-existing user", 3, "", false},
}

func TestGetUserName(t *testing.T) {
	Ids = map[int]string{
		1: "Alice",
		2: "Bob",
	}

	for _, tc := range getUserNameTestCases {
		t.Run(tc.name, func(t *testing.T) {
			result, ok := GetUserName(tc.chatId)
			if result != tc.expected || ok != tc.found {
				t.Errorf("expected (%q, %v), got (%q, %v)", tc.expected, tc.found, result, ok)
			}
		})
	}
}

// SendMessageToUnregisteredUser test data preparation
type sendMessageTestCase struct {
	name         string
	inputChatID  int64
	expectedText string
}

var sendMessageTestCases = []sendMessageTestCase{
	{
		name:         "unregistered user with valid ID",
		inputChatID:  12345,
		expectedText: "This is a private bot. You are not authorized to use it",
	},
}

type mockMessenger struct {
	lastMesage *bot.SendMessageParams
}

func (m *mockMessenger) SendMessage(ctx context.Context, params *bot.SendMessageParams) (*models.Message, error) {
	m.lastMesage = params
	return &models.Message{}, nil
}

func TestSendMessageToUnregisteredUser(t *testing.T) {
	ctx := context.Background()

	for _, tc := range sendMessageTestCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &mockMessenger{}

			SendMessageToUnregisteredUser(ctx, mock, tc.inputChatID)

			if mock.lastMesage == nil {
				t.Fatal("expected message to be sent, but none was")
			}
			if mock.lastMesage.ChatID != tc.inputChatID {
				t.Errorf("expected ChatID %d, got %d", tc.inputChatID, mock.lastMesage.ChatID)
			}
			if mock.lastMesage.Text != tc.expectedText {
				t.Errorf("expected Text %q, got %q", tc.expectedText, mock.lastMesage.Text)
			}
		})
	}
}
