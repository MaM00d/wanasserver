package persona

import (
	api "Server/elapi"
	user "Server/user"

	//"Server/eldb"
	db "Server/eldb"
	//"context"
	//"encoding/json"
	//"fmt"
	"log/slog"
	"net/http"
	"strconv"
)

type recent struct {
	ap     *api.ElApi
	db     *db.Storage
	elchat *int
}

func (s *ElChat) getrecentchats(w http.ResponseWriter, r *http.Request) error {
	slog.Info("Handling getrecentchats")
	eluserid := user.Getidfromheader(r)
	if eluserid < 0 {
		s.ap.WriteError(w, http.StatusUnauthorized, "invalid token")
	}

	// Assuming the persona ID is passed in as a query parameter named "persona_id"
	personaid := r.URL.Query().Get("persona_id")

	elpersonaid, err := strconv.Atoi(personaid)
	if err != nil {
		return s.ap.WriteError(w, http.StatusBadRequest, "invalid persona ID")
	}

	elchat, err := s.GetRecentChatByUserId(elpersonaid, eluserid)
	if err != nil {
		if err == s.db.NotFound {
			return s.ap.WriteJSON(w, http.StatusOK, map[string]interface{}{"chat": "", "persona": ""})
		}
		if err == s.PersonaDoesnotExsist {
			return s.ap.WriteError(w, http.StatusNotFound, "Persona doesn't exist")
		}
		return err
	}

	// Assuming GetPersonaByID retrieves the persona details
	persona, err := s.GetPersonaByID(elpersonaid)
	if err != nil {
		return err
	}

	response := map[string]interface{}{
		"chat":    elchat,
		"persona": persona,
	}

	return s.ap.WriteJSON(w, http.StatusOK, response)
}

// GetRecentChatByUserId returns the most recent chat for a given user and persona
func (s *ElChat) GetRecentChatByUserId(personaID, userID int) (Chat, error) {
	chats, err := s.GetChatsByUserId(personaID, userID)
	if err != nil {
		return Chat{}, err
	}
	if len(chats) == 0 {
		return Chat{}, errors.New("no chats found")
	}

	// Sort chats by timestamp in descending order
	sort.Slice(chats, func(i, j int) bool {
		return chats[i].Timestamp > chats[j].Timestamp
	})

	// Return the most recent chat
	return chats[0], nil
}

/*
	func (s *Storage) GetRecentChatMessage() (string, error) {
	    var message string
	    query := "SELECT message FROM chat_messages ORDER BY created_at DESC LIMIT 1"

	    conn, err := s.pool.Acquire(context.Background())
	    if err != nil {
	        return "", err
	    }
	    defer conn.Release()

	    row := conn.QueryRow(context.Background(), query)
	    err = row.Scan(&message)
	    if err != nil {
	        return "", err
	    }

	    return message, nil
	}

	func HandleRequests() {
	    // Set up the PostgreSQL database connection
	    config, err := pgxpool.ParseConfig("postgresql://username:password@localhost:5432/database_name")
	    if err != nil {
	        log.Fatal(err)
	    }
	    pool, err := pgxpool.ConnectConfig(context.Background(), config)
	    if err != nil {
	        log.Fatal(err)
	    }
	    defer pool.Close()

	    storage := NewStorage(pool)

	    // Define an HTTP handler to retrieve the recent chat message
	    http.HandleFunc("/recent-message", func(w http.ResponseWriter, r *http.Request) {
	        message, err := storage.GetRecentChatMessage()
	        if err != nil {
	            w.WriteHeader(http.StatusInternalServerError)
	            fmt.Fprintf(w, "Error: %v", err)
	            return
	        }

	        response := struct {
	            Message string `json:"message"`
	        }{Message: message}

	        w.Header().Set("Content-Type", "application/json")
	        json.NewEncoder(w).Encode(response)
	    })

	    log.Println("recent message ")

}

	func NewStorage(pool *pgxpool.Pool) {
		panic("unimplemented")
	}
*/
func (elchat *recent) RoutingFetch() {
	slog.Info("RecentChat Routes")
	elchat.ap.Route(elchat.recent, "POST")
	elchat.ap.Route(elchat.recent, "GET")
}
