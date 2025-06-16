package comm_test

import (
	"image/color"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bmcszk/gptrts/pkg/comm"
	"github.com/bmcszk/gptrts/pkg/game"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

// Helper function to dial WebSocket and avoid bodyclose linter false positive
func dialWebSocket(url string) (*websocket.Conn, error) {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	return conn, err
}

func TestNewClient(t *testing.T) {
	// Create a test WebSocket connection
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatalf("Failed to upgrade connection: %v", err)
		}
		defer ws.Close()
	}))
	defer server.Close()

	// Convert http://... to ws://...
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/"

	ws, err := dialWebSocket(wsURL)
	if err != nil {
		t.Fatalf("Failed to dial WebSocket: %v", err)
	}
	defer ws.Close()

	client := comm.NewClient(ws)

	if client == nil {
		t.Fatal("expected client, got nil")
	}

	if !client.Connected {
		t.Error("expected client to be connected")
	}

	// WebSocket connection is private and set correctly if Connected is true
}

func TestClient_Send_Success(t *testing.T) {
	server := createSendTestServer(t)
	defer server.Close()

	client := setupTestClient(t, server)
	action := createTestPlayerJoinAction()

	err := client.Send(action)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func createSendTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatalf("Failed to upgrade connection: %v", err)
		}
		defer ws.Close()

		validateReceivedAction(t, ws)
	}))
}

func validateReceivedAction(t *testing.T, ws *websocket.Conn) {
	t.Helper()
	var action game.PlayerJoinAction
	err := ws.ReadJSON(&action)
	if err != nil {
		t.Errorf("Failed to read JSON: %v", err)
		return
	}

	if action.Type != game.PlayerJoinActionType {
		t.Errorf("expected action type %s, got %s", game.PlayerJoinActionType, action.Type)
	}

	if action.Payload.Name != "TestPlayer" {
		t.Errorf("expected player name TestPlayer, got %s", action.Payload.Name)
	}
}

func setupTestClient(t *testing.T, server *httptest.Server) *comm.Client {
	t.Helper()
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/"

	ws, err := dialWebSocket(wsURL)
	if err != nil {
		t.Fatalf("Failed to dial WebSocket: %v", err)
	}
	t.Cleanup(func() { ws.Close() })

	client := comm.NewClient(ws)
	client.PlayerId = game.PlayerIdType(uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"))
	return client
}

func createTestPlayerJoinAction() game.PlayerJoinAction {
	return game.PlayerJoinAction{
		Type: game.PlayerJoinActionType,
		Payload: game.Player{
			Id:    game.PlayerIdType(uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")),
			Name:  "TestPlayer",
			Color: color.RGBA{255, 0, 0, 255},
		},
	}
}

func TestClient_Send_Disconnected(t *testing.T) {
	// Create a mock WebSocket connection (doesn't matter if it's valid since we'll mark as disconnected)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatalf("Failed to upgrade connection: %v", err)
		}
		defer ws.Close()
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/"

	ws, err := dialWebSocket(wsURL)
	if err != nil {
		t.Fatalf("Failed to dial WebSocket: %v", err)
	}
	defer ws.Close()

	client := comm.NewClient(ws)
	client.Connected = false // Mark as disconnected

	action := game.PlayerJoinAction{
		Type: game.PlayerJoinActionType,
		Payload: game.Player{
			Id:   game.PlayerIdType(uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")),
			Name: "TestPlayer",
		},
	}

	err = client.Send(action)
	if err != nil {
		t.Errorf("expected no error when disconnected, got %v", err)
	}
}

func TestClient_HandleInMessages_Success(t *testing.T) {
	server := createReceiveTestServer(t)
	defer server.Close()

	client := setupTestClient(t, server)
	action := receiveAndValidateAction(t, client)
	validatePlayerJoinAction(t, action)
}

func createReceiveTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatalf("Failed to upgrade connection: %v", err)
		}
		defer ws.Close()

		sendTestAction(t, ws)
	}))
}

func sendTestAction(t *testing.T, ws *websocket.Conn) {
	t.Helper()
	action := createTestPlayerJoinAction()
	err := ws.WriteJSON(action)
	if err != nil {
		t.Errorf("Failed to write JSON: %v", err)
	}
}

func receiveAndValidateAction(t *testing.T, client *comm.Client) game.Action {
	t.Helper()
	action, err := client.HandleInMessages()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if action.GetType() != game.PlayerJoinActionType {
		t.Errorf("expected action type %s, got %s", game.PlayerJoinActionType, action.GetType())
	}

	return action
}

func validatePlayerJoinAction(t *testing.T, action game.Action) {
	t.Helper()
	playerJoinAction, ok := action.(game.PlayerJoinAction)
	if !ok {
		t.Fatalf("expected PlayerJoinAction, got %T", action)
	}

	if playerJoinAction.Payload.Name != "TestPlayer" {
		t.Errorf("expected player name TestPlayer, got %s", playerJoinAction.Payload.Name)
	}
}

func TestClient_PlayerId(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatalf("Failed to upgrade connection: %v", err)
		}
		defer ws.Close()
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/"

	ws, err := dialWebSocket(wsURL)
	if err != nil {
		t.Fatalf("Failed to dial WebSocket: %v", err)
	}
	defer ws.Close()

	client := comm.NewClient(ws)

	testPlayerId := game.PlayerIdType(uuid.MustParse("550e8400-e29b-41d4-a716-446655440123"))
	client.PlayerId = testPlayerId

	if client.PlayerId != testPlayerId {
		t.Errorf("expected player ID %s, got %s", testPlayerId, client.PlayerId)
	}
}

func TestClient_Connected(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatalf("Failed to upgrade connection: %v", err)
		}
		defer ws.Close()
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/"

	ws, err := dialWebSocket(wsURL)
	if err != nil {
		t.Fatalf("Failed to dial WebSocket: %v", err)
	}
	defer ws.Close()

	client := comm.NewClient(ws)

	// Should start connected
	if !client.Connected {
		t.Error("expected client to start connected")
	}

	// Test marking as disconnected
	client.Connected = false
	if client.Connected {
		t.Error("expected client to be marked as disconnected")
	}
}
