package squirrel

import (
	"encoding/json"
	"es/internal/utils"
	"net/http"
)

func (s *Squirrel) update(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(utils.GetKeys(s.entangledSquirrelURLs))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (s *Squirrel) knownSquirrels(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(utils.GetKeys(s.entangledSquirrelURLs))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (s *Squirrel) addSquirrel(w http.ResponseWriter, r *http.Request) {
	var squirrels []string
	err := json.NewDecoder(r.Body).Decode(&squirrels)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, currentSquirrel := range squirrels {
		if found := s.entangledSquirrelURLs[currentSquirrel]; !found {
			s.entangledSquirrelURLs[currentSquirrel] = true
		}
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.entangledSquirrelURLs)
}
