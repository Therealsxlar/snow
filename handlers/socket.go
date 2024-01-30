package handlers

import (
	"github.com/ectrc/snow/aid"
	"github.com/ectrc/snow/socket"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func MiddlewareWebsocket(c *fiber.Ctx) error {
	if !websocket.IsWebSocketUpgrade(c) {
		return fiber.ErrUpgradeRequired
	}

	var protocol string

	switch c.Get("Sec-WebSocket-Protocol") {
	case "xmpp":
		protocol = "jabber"
	default:
		protocol = "matchmaking"
	}

	c.Locals("identifier", uuid.New().String())
	c.Locals("protocol", protocol)

	return c.Next()
}

func WebsocketConnection(c *websocket.Conn) {
	protocol := c.Locals("protocol").(string)
	identifier := c.Locals("identifier").(string)

	switch protocol {
	case "jabber":
		socket.JabberSockets.Set(identifier, socket.NewJabberSocket(c, identifier, socket.JabberData{}))
		socket.HandleNewJabberSocket(identifier)
	case "matchmaking":
		// socket.MatchmakerSockets.Set(identifier, socket.NewMatchmakerSocket(c, socket.MatchmakerData{}))
	default:
		aid.Print("Invalid protocol: " + protocol)
	}
}

func GetConnectedSockets(c *fiber.Ctx) error {
	jabber := []socket.Socket[socket.JabberData]{}
	socket.JabberSockets.Range(func(key string, value *socket.Socket[socket.JabberData]) bool {
		jabber = append(jabber, *value)
		return true
	})

	matchmaking := []socket.Socket[socket.MatchmakerData]{}
	socket.MatchmakerSockets.Range(func(key string, value *socket.Socket[socket.MatchmakerData]) bool {
		matchmaking = append(matchmaking, *value)
		return true
	})

	return c.Status(200).JSON(aid.JSON{
		"jabber": jabber,
		"matchmaking": matchmaking,
	})
}