package go_coinbasev3

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
	ErrNoURL              = fmt.Errorf("no url provided")
	ErrInvalidReadChannel = fmt.Errorf("read channel is invalid")
	ErrNotConnected       = fmt.Errorf("not connected")
)

type WsClient struct {
	conn             *websocket.Conn
	url              string
	subscriptionData []byte
	chans            channels
	isShutdown       bool
	cbs              callbacks
	useBackoff       bool
	ctx              context.Context
	cancel           context.CancelFunc
	debug            bool
}

type channels struct {
	read  chan []byte
	recon chan bool
	done  chan bool
}

type callbacks struct {
	onConnect    func()
	onDisconnect func()
	onReconnect  func()
}

// WsClientConfig is the configuration for the websocket client.
type WsClientConfig struct {
	Url              string
	ReadChannel      chan []byte
	SubscriptionData []byte
	OnConnect        func()
	OnDisconnect     func()
	OnReconnect      func()
	UseBackoff       bool
	Debug            bool
}

// tryValidate validates the websocket client configuration.
func (c *WsClientConfig) tryValidate() error {
	if c.Url == "" {
		return ErrNoURL
	}
	if c.ReadChannel == nil {
		return ErrInvalidReadChannel
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
		chans: channels{
			read:  cfg.ReadChannel,
			recon: make(chan bool),
		},
		cbs: callbacks{
			onConnect:    cfg.OnConnect,
			onDisconnect: cfg.OnDisconnect,
			onReconnect:  cfg.OnReconnect,
		},
		subscriptionData: cfg.SubscriptionData,
		ctx:              ctx,
		cancel:           cancel,
		useBackoff:       cfg.UseBackoff,
		debug:            cfg.Debug,
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

	if len(c.subscriptionData) > 0 {
		err = c.Write(c.subscriptionData)
		if err != nil {
			_ = c.conn.Close()
			return nil, err
		}
	}

	c.cbs.onConnect()
	go c.read()
	return conn, nil
}

// Connect connects to the websocket server.
func (c *WsClient) Connect() (*websocket.Conn, error) {
	if c.url == "" {
		return nil, ErrNoURL
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

// Close closes the websocket connection. This will also close the read channel and stop the reconnect channel.
func (c *WsClient) Close() {
	c.cancel()

	c.isShutdown = true
	err := c.conn.Close()
	if err != nil {
		return
	}
}

// ReadChan returns the channel that receives messages from the websocket connection.
func (c *WsClient) ReadChan() chan []byte {
	return c.chans.read
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
		case <-c.chans.recon:
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
				c.chans.recon <- true
			}
			return
		}

		c.chans.read <- message
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
