package elai

import (
	"fmt"
	"log/slog"
	"net"
	"time"
)

type Aiserver struct {
	address string
	conn    *net.TCPConn
}

func InitAiServer() *Aiserver {
	address := "localhost:12345"
	slog.Info("Start Init Ai Server")
	ais := &Aiserver{
		address: address,
	}
	err := ais.connectToServer()
	for err != nil {
		slog.Error("Init Ai Server", "El Error:", err)
		slog.Warn("trying to re connect...")
		time.Sleep(2 * time.Second)
		err = ais.connectToServer()
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

// Send a message to the server and receive the response
func (ais *Aiserver) SendMessage(message string) (string, error) {
	_, err := ais.conn.Write([]byte(message))
	if err != nil {
		println("Write to server failed:", err.Error())
	}

	println("write to server = ", message)

	reply := make([]byte, 1024)

	_, err = ais.conn.Read(reply)
	if err != nil {
		println("Write to server failed:", err.Error())
	}

	println("reply from server=", string(reply))

	return string(reply), nil
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
