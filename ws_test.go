package go_coinbasev3

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

var upgrader = websocket.Upgrader{}

func TestClient_Close(t *testing.T) {
	ha := handler{}

	s, ws := newWSServer(t, ha)
	defer s.Close()
	defer func(ws *websocket.Conn) {
		_ = ws.Close()
	}(ws)

	chRead := make(chan []byte)

	onRecon := func() {
		t.Fatal("should not reconnect")
	}

	cl, err := NewWsClient(WsClientConfig{
		Url:              "https://badurl.com",
		ReadChannel:      chRead,
		SubscriptionData: []byte{},
		OnConnect:        func() {},
		OnDisconnect:     func() {},
		OnReconnect:      onRecon,
	})
	if err != nil {
		t.Fatalf("NewWsClient: %v", err)
	}

	ws, err = cl.ConnectWithUrl(makeWsProto(s.URL))
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}

	cl.Shutdown()
	time.Sleep(1 * time.Second)

	err = ws.NetConn().Close()
	if !strings.HasSuffix(err.Error(), "use of closed network connection") {
		t.Fatalf("Shutdown: %v", err)
	}
}

func TestClient_Reconnecting(t *testing.T) {
	ha := handler{}
	s, ws := newWSServer(t, ha)
	defer s.Close()
	defer func(ws *websocket.Conn) {
		err := ws.Close()
		if err != nil {
			t.Fatalf("Shutdown: %v", err)
		}
	}(ws)

	count := 0
	want := 5 // original connection + 4 underlying connection closes in loop
	reconCount := 0
	reconWant := 4 // 4 underlying connection closes in loop

	onConn := func() {
		count++
	}
	onRecon := func() {
		reconCount++
	}

	chRead := make(chan []byte)
	cl, err := NewWsClient(WsClientConfig{
		Url:              makeWsProto(s.URL),
		ReadChannel:      chRead,
		SubscriptionData: []byte{},
		OnConnect:        onConn,
		OnDisconnect:     func() {},
		OnReconnect:      onRecon,
	})
	if err != nil {
		t.Fatalf("NewWsClient: %v", err)
	}

	_, err = cl.Connect()
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}

	for i := 1; i <= want-1; i++ {
		time.Sleep(100 * time.Millisecond)
		_ = cl.conn.Close()
	}
	time.Sleep(100 * time.Millisecond)
	cl.Shutdown()

	if count != want {
		t.Errorf("count = %d; want %d", count, want)
	}
	if reconCount != reconWant {
		t.Errorf("reconCount = %d; want %d", reconCount, reconWant)
	}
}

func TestWsClient_Backoff_ShouldNotTriggerRecon(t *testing.T) {
	ha := handler{}
	s, ws := newWSServer(t, ha)
	defer s.Close()
	defer func(ws *websocket.Conn) {
		err := ws.Close()
		if err != nil {
			t.Fatalf("Shutdown: %v", err)
		}
	}(ws)

	reconCount := 0
	reconWant := 0

	chRead := make(chan []byte)
	cl, err := NewWsClient(WsClientConfig{
		Url:              makeWsProto(s.URL),
		ReadChannel:      chRead,
		SubscriptionData: []byte{},
		OnConnect:        func() {},
		OnDisconnect:     func() {},
		OnReconnect: func() {
			reconCount++
		},
		UseBackoff: true,
	})
	if err != nil {
		t.Fatalf("NewWsClient: %v", err)
	}

	_, err = cl.Connect()
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}

	for i := 1; i <= 2; i++ {
		time.Sleep(400 * time.Millisecond)
		_ = cl.conn.Close()
	}
	time.Sleep(100 * time.Millisecond)
	cl.Shutdown()

	if reconCount != reconWant {
		t.Errorf("reconCount = %d; want %d", reconCount, reconWant)
	}
}

func TestWsClient_Backoff_ShouldTriggerRecon(t *testing.T) {
	ha := handler{}
	s, ws := newWSServer(t, ha)
	defer s.Close()
	defer func(ws *websocket.Conn) {
		err := ws.Close()
		if err != nil {
			t.Fatalf("Shutdown: %v", err)
		}
	}(ws)

	reconCount := 0
	reconWant := 1

	onConn := func() {}
	onRecon := func() {
		reconCount++
	}
	onDisc := func() {}

	chRead := make(chan []byte)
	cl, err := NewWsClient(WsClientConfig{
		Url:              makeWsProto(s.URL),
		ReadChannel:      chRead,
		SubscriptionData: []byte{},
		OnConnect:        onConn,
		OnDisconnect:     onDisc,
		OnReconnect:      onRecon,
		UseBackoff:       true,
	})
	if err != nil {
		t.Fatalf("NewWsClient: %v", err)
	}

	_, err = cl.Connect()
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}

	for i := 1; i <= 2; i++ {
		time.Sleep(2000 * time.Millisecond)
		_ = cl.conn.Close()
	}
	time.Sleep(100 * time.Millisecond)
	cl.Shutdown()

	if reconCount != reconWant {
		t.Errorf("reconCount = %d; want %d", reconCount, reconWant)
	}
}

func makeWsProto(s string) string {
	return "ws" + strings.TrimPrefix(s, "http")
}

type handler struct {
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	//defer ws.Shutdown()
}
func newWSServer(t *testing.T, h http.Handler) (*httptest.Server, *websocket.Conn) {
	t.Helper()
	s := httptest.NewServer(h)
	wsURL := makeWsProto(s.URL)

	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatal(err)
	}

	return s, ws
}
