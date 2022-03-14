package squirrel

import (
	"encoding/json"
	"es/internal/consts"
	"es/internal/utils"
	"fmt"
	"net/http"
)

func (s *Squirrel) flashReq(w http.ResponseWriter, r *http.Request) {
	fmt.Println("flash request received")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
	go s.flash(consts.NoOfFlashes)
}

func (s *Squirrel) knownSquirrelsReq(w http.ResponseWriter, r *http.Request) {
	fmt.Println("known squirrels request received")
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
	fmt.Println("add squirrel request received")
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
