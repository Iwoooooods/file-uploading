package p2p

import (
	"encoding/json"
	"log"
	"net"
)

const (
	bootstrapAddress = "localhost:8080"
)

func (n *Node) Serve() {
	listener, err := net.Listen("tcp", n.Address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		panic(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("failed to accept connection: %v", err)
			continue
		}
		go n.handleConnection(conn)
	}
}

func (n *Node) handleConnection(conn net.Conn) {
	defer conn.Close()
	decoder := json.NewDecoder(conn)

	var msg map[string]any
	if err := decoder.Decode(&msg); err != nil {
		log.Fatalf("failed to decode message: %v", err)
		return
	}

	msgType := msg["type"].(string)
	switch msgType {
	case "ping":
		log.Println("Hello World!")
	case "connect":
		// add new peer to node list
		id := msg["id"].(string)
		address := msg["address"].(string)
		n.Peers[id] = &Peer{ID: id, Address: address}

		// expose peers to new peer
		resp := map[string]any{
			"type": "peers",
			"peers": n.peersAsList(),
		}
		encoder := json.NewEncoder(conn)
		encoder.Encode(resp)
	}
}

func (n *Node) peersAsList() []map[string]string {
	peers := []map[string]string{}
	for _, peer := range n.Peers {
		peers = append(peers, map[string]string{"id": peer.ID, "address": peer.Address})
	}
	return peers
}

func (n *Node) connectToBootstrapNode(bootstrapAddress string) {
	conn, err := net.Dial("tcp", bootstrapAddress)
	if err != nil {
		log.Fatalf("failed to connect to bootstrap node: %v", err)
		return
	}
	defer conn.Close()

	req := map[string]any{
		"type": "connect",
		"id":   n.ID,
		"address": n.Address,
	}

	encoder := json.NewEncoder(conn)
		encoder.Encode(req)

	decoder := json.NewDecoder(conn)
	var resp map[string]any
	if err := decoder.Decode(&resp); err != nil {
		log.Fatalf("failed to decode response: %v", err)
		return
	}

	if resp["type"] != "peers" {
		log.Fatalf("expected peers response, got %v", resp["type"])
		return
	}

	peers := resp["peers"].([]map[string]any)

	for _, peer := range peers {
		n.Peers[peer["id"].(string)] = &Peer{ID: peer["id"].(string), Address: peer["address"].(string)}
	}

	log.Printf("connected to %d peers", len(n.Peers))
	return
}
