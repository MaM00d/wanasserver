package elai

import (
	"bytes"
	"fmt"
	"log/slog"
	"net"
)

type Aiserver struct {
	address string
	conn    net.Conn
}

func InitAiServer() *Aiserver {
	slog.Info("Start Init Ai Server")

	address := "97.119.128.191:40378"
	conn, err := net.Dial("tcp", address)
	if err != nil {
		slog.Error("Error connecting to server:", err)
	}
	// defer conn.Close()

	ais := &Aiserver{
		address: address,
		conn:    conn,
	}

	slog.Info("Done Init Ai Server")
	return ais
}

func (ais *Aiserver) connectToServer() error {
	slog.Info("connecting to ai server")
	tcpAddr, err := net.ResolveTCPAddr("tcp", ais.address)
	if err != nil {
		slog.Error("ResolveTCPAddr failed:", "Error:", err)
		return err
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		slog.Error("connecting to ai server", "error: ", err)
		return fmt.Errorf("error connecting to server: %w", err)
	}
	ais.conn = conn
	// defer conn.Close()

	slog.Info("connected successfully to ai server")
	return nil
}

func (ais *Aiserver) SendMessage(message string, history string) (*string, error) {
	// Send message to the server
	_, err := ais.conn.Write([]byte(history + "~~~" + message + "!~!~!"))
	if err != nil {
		slog.Error("Error sending message:", err)
	}

	// Read response from the server
	var responseBuffer bytes.Buffer
	tmp := make([]byte, 1024)
	for {
		n, err := ais.conn.Read(tmp)
		if err != nil {
			fmt.Println("Error reading response:", err)
			return nil, err
		}
		responseBuffer.Write(tmp[:n])
		if bytes.Contains(responseBuffer.Bytes(), []byte("!~!~!")) {
			break
		}
	}

	// Get the complete response
	response := responseBuffer.String()
	response = response[:len(response)-5] // Remove the delimiter

	slog.Info("Server response:", "response", response)
	return &response, nil
}

// func Tokenize(text string) error {
// 	tk, err := tokenizers.FromFile(
// 		"/home/me/loc/.hfcache/hub/models--FreedomIntelligence--AceGPT-13B/snapshots/4ceb8c997a28c82f8cd3f55834d690f19456f5df/tokenizer.json",
// 	)
// 	if err != nil {
// 		return err
// 	}
// 	fmt.Println("Vocab size:", tk.VocabSize())
// 	// Vocab size: 30522
// 	fmt.Println(tk.Encode(text, true))
// 	defer tk.Close()
// 	return nil
// }
