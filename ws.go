package coinbasev3

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"time"
)

const (
	initialBackoff    = 1 * time.Second  // Initial backoff duration
	maxBackoff        = 1 * time.Minute  // Maximum backoff duration
	connectionTimeout = 15 * time.Second // Connection timeout duration
)

var (
	ErrInvalidReadChannel = fmt.Errorf("read channel is invalid")
	ErrNoApiKey           = fmt.Errorf("no api key provided")
	ErrNoSecretKey        = fmt.Errorf("no secret key provided")
	ErrNotConnected       = fmt.Errorf("not connected")
)

// WsClientConfig is the configuration struct for creating a new websocket client.
type WsClientConfig struct {
	Url          string      // optional. defaults to "wss://advanced-trade-ws.coinbase.com"
	ReadChannel  chan []byte // required for receiving messages from the websocket connection
	WsChannels   []WsChannel // required for subscribing to innerChannels on the websocket connection
	ApiKey       string      // required for signing websocket messages
	SecretKey    string      // required for signing websocket messages
	OnConnect    func()      // optional. called when the websocket connection is established
	OnDisconnect func()      // optional. called when the websocket connection is closed
	OnReconnect  func()      // optional. called when the websocket connection is re-established
	UseBackoff   bool        // optional. defaults to false. uses an exponential backoff strategy with jitter
	Debug        bool        // optional. defaults to false. prints debug messages
}

func NewWsClientConfig(apiKey, secretKey string, readCh chan []byte, wsChannels []WsChannel) WsClientConfig {
	return WsClientConfig{
		ApiKey:      apiKey,
		SecretKey:   secretKey,
		ReadChannel: readCh,
		WsChannels:  wsChannels,
		UseBackoff:  true,
		Debug:       true,
	}
}

// tryValidate validates the websocket client configuration.
func (c *WsClientConfig) tryValidate() error {
	if c.ApiKey == "" {
		return ErrNoApiKey
	}
	if c.SecretKey == "" {
		return ErrNoSecretKey
	}
	if c.ReadChannel == nil {
		return ErrInvalidReadChannel
	}

	if c.Url == "" {
		c.Url = "wss://advanced-trade-ws.coinbase.com"
	}
	if c.OnConnect == nil {
		c.OnConnect = func() {}
	}
	if c.OnDisconnect == nil {
		c.OnDisconnect = func() {}
	}
	if c.OnReconnect == nil {
		c.OnReconnect = func() {}
	}
	return nil
}

// WsClient is an automatically reconnecting websocket client.
type WsClient struct {
	conn          *websocket.Conn
	url           string
	wsChannels    []WsChannel
	innerChannels channels
	cbs           callbacks
	isShutdown    bool
	useBackoff    bool
	debug         bool
	apiKey        string
	secretKey     string
	ctx           context.Context
	cancel        context.CancelFunc
}

// channels contains the read channel given by the developer and the internal reconnection channel.
type channels struct {
	read  chan []byte
	recon chan bool
}

// callbacks contains the callbacks used by the developer to handle events.
type callbacks struct {
	onConnect    func()
	onDisconnect func()
	onReconnect  func()
}

// NewWsClient creates a new websocket client.
func NewWsClient(cfg WsClientConfig) (*WsClient, error) {
	err := cfg.tryValidate()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	c := &WsClient{
		conn: nil,
		url:  cfg.Url,
		innerChannels: channels{
			read:  cfg.ReadChannel,
			recon: make(chan bool),
		},
		cbs: callbacks{
			onConnect:    cfg.OnConnect,
			onDisconnect: cfg.OnDisconnect,
			onReconnect:  cfg.OnReconnect,
		},
		wsChannels: cfg.WsChannels,
		apiKey:     cfg.ApiKey,
		secretKey:  cfg.SecretKey,
		ctx:        ctx,
		cancel:     cancel,
		useBackoff: cfg.UseBackoff,
		debug:      cfg.Debug,
	}

	go c.initReconnectChannel()
	return c, nil
}

// ConnectWithUrl connects to the websocket server using the provided url.
func (c *WsClient) ConnectWithUrl(url string) (*websocket.Conn, error) {
	if c.url == "" {
		c.url = url
	}

	ctx, cancel := context.WithTimeout(context.Background(), connectionTimeout)
	defer cancel()

	conn, _, err := websocket.DefaultDialer.DialContext(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	c.conn = conn

	err = c.subscribeToChannels()
	if err != nil {
		_ = c.conn.Close()
		return nil, err
	}

	c.cbs.onConnect()
	go c.read()
	return conn, nil
}

// subscribeToChannels subscribes to the provided channels on the websocket connection.
func (c *WsClient) subscribeToChannels() error {
	if len(c.wsChannels) > 0 {
		for ch := range c.wsChannels {
			err := c.Write(c.wsChannels[ch].marshal(c.apiKey, c.secretKey))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Connect connects to the websocket server.
func (c *WsClient) Connect() (*websocket.Conn, error) {
	if c.url == "" {
		c.url = "wss://advanced-trade-ws.coinbase.com" // default url
	}
	return c.ConnectWithUrl(c.url)
}

// Write writes data to the websocket connection.
func (c *WsClient) Write(data []byte) error {
	if c.conn == nil {
		return ErrNotConnected
	}
	if c.conn.NetConn().LocalAddr() == nil {
		return ErrNotConnected
	}

	err := c.conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		return err
	}
	return nil
}

// Shutdown closes the websocket connection. This will also close the read channel and the underlying reconnect channel.
func (c *WsClient) Shutdown() {
	c.cancel()

	c.isShutdown = true
	close(c.innerChannels.read)
	close(c.innerChannels.recon)
	err := c.conn.Close()
	if err != nil {
		return
	}
}

// ReadChan returns the channel that receives messages from the websocket connection.
func (c *WsClient) ReadChan() chan []byte {
	return c.innerChannels.read
}

// initReconnectChannel will attempt to reconnect to the websocket server using an exponential backoff strategy with jitter.
func (c *WsClient) initReconnectChannel() {
	backoff := initialBackoff
	for {
		select {
		case <-c.ctx.Done():
			if c.debug {
				c.printf("Reconnect channel closed due to context cancellation.")
			}
			return
		case <-c.innerChannels.recon:
			_ = c.conn.Close()

			if c.useBackoff {
				jitter := time.Duration(rand.Intn(1000)) * time.Millisecond
				time.Sleep(backoff + jitter)
			}

			_, err := c.Connect()
			if err != nil {
				c.printf("Reconnection attempt failed: %s\n", err)
				backoff = calculateBackoff(backoff, maxBackoff)
				continue
			}

			backoff = initialBackoff
			c.cbs.onReconnect()
		}
	}
}

// read reads messages from the websocket connection. It will attempt to reconnect if the connection is closed.
func (c *WsClient) read() {
	defer func() {
		c.cbs.onDisconnect()
		if err := c.conn.Close(); err != nil {
			c.printf("Error closing WebSocket connection: %s\n", err)
		}
	}()

	for {
		if c.isShutdown {
			return
		}

		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.printf("WebSocket closed unexpectedly: %v\n", err)
			} else {
				c.printf("WebSocket read error: %v\n", err)
			}
			if !c.isShutdown {
				c.innerChannels.recon <- true
			}
			return
		}

		c.innerChannels.read <- message
	}
}

func (c *WsClient) printf(format string, a ...any) {
	if c.debug {
		log.Printf(format, a...)
	}
}

func (c *WsClient) println(a ...any) {
	if c.debug {
		log.Println(a...)
	}
}

// calculateBackoff calculates the next backoff duration.
// It doubles the current backoff, not exceeding the maxBackoff.
func calculateBackoff(currentBackoff, maxBackoff time.Duration) time.Duration {
	nextBackoff := currentBackoff * 2
	if nextBackoff > maxBackoff {
		return maxBackoff
	}
	return nextBackoff
}
