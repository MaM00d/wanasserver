package elai

// import (
// 	"fmt"
// 	"log/slog"
// 	"net"
// 	"strconv"
// 	"time"
// )

// type Aiserver struct {
// 	address string
// 	conn    *net.TCPConn
// }
//
// func InitAiServer() *Aiserver {
// 	address := "localhost:12345"
// 	slog.Info("Start Init Ai Server")
// 	ais := &Aiserver{
// 		address: address,
// 	}
// 	err := ais.connectToServer()
// 	for err != nil {
// 		slog.Error("Init Ai Server", "El Error:", err)
// 		slog.Warn("trying to re connect...")
// 		time.Sleep(2 * time.Second)
// 		err = ais.connectToServer()
// 	}
//
// 	slog.Info("Done Init Ai Server")
// 	return ais
// }
//
// func (ais *Aiserver) connectToServer() error {
// 	slog.Info("connecting to ai server")
// 	tcpAddr, err := net.ResolveTCPAddr("tcp", ais.address)
// 	if err != nil {
// 		slog.Error("ResolveTCPAddr failed:", "Error:", err)
// 		return err
// 	}
// 	conn, err := net.DialTCP("tcp", nil, tcpAddr)
// 	if err != nil {
// 		slog.Error("connecting to ai server", "error: ", err)
// 		return fmt.Errorf("error connecting to server: %w", err)
// 	}
// 	ais.conn = conn
// 	// defer conn.Close()
//
// 	slog.Info("connected successfully to ai server")
// 	return nil
// }
//
// // Send a message to the server and receive the response
// func PadLeft(str string, length int) string {
// 	for len(str) < length {
// 		str = "0" + str
// 	}
// 	return str
// }
//
// func firstN(str string, n int) string {
// 	v := []rune(str)
// 	if n >= len(v) {
// 		return str
// 	}
// 	return string(v[:n])
// }
//
// func (ais *Aiserver) SendMessage(msg string, history string) (string, error) {
// 	message := history + "~~~" + msg + "!~!~!"
// 	msglen := PadLeft(strconv.Itoa(len(message)), 10)
//
// 	_, err := ais.conn.Write([]byte(msglen + message))
// 	if err != nil {
// 		slog.Error("Write to server failed:", "error", err.Error())
// 	}
//
// 	slog.Info("write to server = ", "message", message)
//
// 	l := make([]byte, 10)      // takes the first two bytes: the Length Message.
// 	i, err := ais.conn.Read(l) // read the bytes.
// 	if err != nil {
// 		slog.Error("Write to server failed:", "err", err.Error())
// 	}
// 	slog.Info("send message", "length of message", i)
// 	lm := string(firstN(string(l), 10)) // convert the bytes into int16.
// 	num, err := strconv.Atoi(lm)
// 	slog.Info("message", "number", lm)
// 	if err != nil {
// 		slog.Error("convert string to int header tcp ai", "error", err)
// 	}
//
// 	reply := make([]byte, num)
//
// 	_, err = ais.conn.Read(reply)
// 	if err != nil {
// 		slog.Error("Write to server failed:", "err", err.Error())
// 	}
//
// 	response := string(reply)
// 	slog.Info("reply from server=", "reply", response)
//
// 	return response, nil
// }
//
// // func Tokenize(text string) error {
// // 	tk, err := tokenizers.FromFile(
// // 		"/home/me/loc/.hfcache/hub/models--FreedomIntelligence--AceGPT-13B/snapshots/4ceb8c997a28c82f8cd3f55834d690f19456f5df/tokenizer.json",
// // 	)
// // 	if err != nil {
// // 		return err
// // 	}
// // 	fmt.Println("Vocab size:", tk.VocabSize())
// // 	// Vocab size: 30522
// // 	fmt.Println(tk.Encode(text, true))
// // 	defer tk.Close()
// // 	return nil
// // }
