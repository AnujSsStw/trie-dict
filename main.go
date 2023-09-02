package main

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/AnujSsStw/goooooooooo/trie"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/google/uuid"
)

type Connected struct {
	ClientId string `json:"clientId"`
}
type Cmd struct {
	Method   string `json:"method"`
	ClientId string `json:"clientId"`
	Value    string `json:"value"`
}
type Res struct {
	Ans []string `json:"ans"`
}

var (
	clients   = make(map[string]*websocket.Conn)
	clientsMu sync.Mutex
)

func main() {
	//init the trie of dict
	a := trie.Trie{RootNode: &trie.Node{}}
	for i := 0; i < len(trie.Char); i++ {
		a.InsertText(trie.Char[i])
	}

	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Use("/ws", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		clientID := uuid.New().String()

		// Add the client to the map
		clientsMu.Lock()
		clients[clientID] = c
		payload := Connected{ClientId: clientID}
		// Notify connected clients of a new user
		clients[clientID].WriteJSON(payload)
		clientsMu.Unlock()

		// Handle incoming WebSocket messages
		for {
			var (
				msg []byte
				err error
			)

			// Read the message from the client
			if _, msg, err = c.ReadMessage(); err != nil {
				log.Println("read:", err)
				break
			}

			var cmd Cmd
			if err = json.Unmarshal(msg, &cmd); err != nil {
				log.Fatalln("some diff method prob.", err)
			}

			res := Res{Ans: trie.Trieee(cmd.Value, &a)}
			// Broadcast the message to all connected clients
			clients[clientID].WriteJSON(&res)
		}

		// When the client disconnects
		clientsMu.Lock()
		// Notify remaining clients of the user leaving ? TODO
		delete(clients, clientID)
		clientsMu.Unlock()

		// Close the WebSocket connection
		c.Close()
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		// Render index template
		return c.Render("index", fiber.Map{
			"Title": "Hello, World!",
		})
	})

	// Start the Fiber app
	app.Listen(":3000")
}
