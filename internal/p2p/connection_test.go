package p2p

import (
	"testing"
)

func TestConnection(t *testing.T) {
	t.Run("bootstrap connection", func(t *testing.T) {
		bootstrapNode := NewNode("bootstrap", bootstrapAddress)
		bootstrapNode.Serve()

		newNode := NewNode("new", "localhost:9191")
		newNode.Serve()
		newNode.connectToBootstrapNode(bootstrapNode.Address)
		bootstrapNode.Peers[newNode.ID] = &Peer{ID: newNode.ID, Address: newNode.Address}
	})
}
