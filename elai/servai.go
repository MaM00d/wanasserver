package elai

import (
	"bufio"
	"fmt"
	"log/slog"
	"net"

	// "github.com/daulet/tokenizers"
)

type Aiserver struct {
	conn    *net.Conn
	address string
}

// func callapi(elmessage string) string {
//
// }
func InitAiServer(address string) *Aiserver {
	slog.Info("Start Init Ai Server")
	ais := &Aiserver{
		address: address,
	}
	err := ais.connectToServer()
	if err != nil {
		slog.Error("Init Ai Server", "El Error:", err)
	}

	slog.Info("Done Init Ai Server")
	return ais
}

// Connect to the server and return the connection object
func (ais *Aiserver) connectToServer() error {
	conn, err := net.Dial("tcp", ais.address)
	if err != nil {
		return fmt.Errorf("error connecting to server: %w", err)
	}
	ais.conn = &conn
	// defer conn.Close()
	return nil
}

// Send a message to the server and receive the response
func (ais *Aiserver) SendMessage(message string) (string, error) {
	// Send the message to the server
	_, err := fmt.Fprintf(*ais.conn, message+"\n")
	if err != nil {
		return "", fmt.Errorf("error sending message: %w", err)
	}

	// Create a buffer to read data from the connection
	reader := bufio.NewReader(*ais.conn)

	// Read the server's response
	response, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("error reading response: %w", err)
	}
	slog.Info(response)

	return response, nil
}

// func main() {
//     // Connect to the server
//     conn, err := connectToServer("localhost:12345")
//     if err != nil {
//         fmt.Println(err)
//         return
//     }
//     defer conn.Close()
//
//     // Read the initial message from the server
//     initialMessage, err := sendMessage(conn, "")
//     if err != nil {
//         fmt.Println(err)
//         return
//     }
//     fmt.Println("Server:", initialMessage)
//
//     // Scanner to read user input from the command line
//     inputReader := bufio.NewReader(os.Stdin)
//
//     for {
//         // Read input from the user
//         fmt.Print("Enter message: ")
//         userInput, err := inputReader.ReadString('\n')
//         if err != nil {
//             fmt.Println("Error reading input:", err)
//             return
//         }
//
//         // Send the user input to the server and get the response
//         response, err := sendMessage(conn, userInput)
//         if err != nil {
//             fmt.Println(err)
//             return
//         }
//         fmt.Println("Server:", response)
//     }
// }

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
