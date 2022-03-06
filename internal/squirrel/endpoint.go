package squirrel

import (
	"encoding/json"
	"es/internal/utils"
	"net/http"
)

func (s *Squirrel) flashReq(w http.ResponseWriter, r *http.Request) {
	// TODO flash here
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
	go s.flash()
}

func (s *Squirrel) knownSquirrelsReq(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	s.entangledSquirrelLock.RLock()
	err := json.NewEncoder(w).Encode(utils.GetKeys(s.entangledSquirrelURLs))
	s.entangledSquirrelLock.RUnlock()

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (s *Squirrel) addSquirrelsReq(w http.ResponseWriter, r *http.Request) {
	var squirrels []string
	err := json.NewDecoder(r.Body).Decode(&squirrels)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.addSquirrels(squirrels)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.entangledSquirrelURLs)
}
