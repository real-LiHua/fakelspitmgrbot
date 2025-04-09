package main

import (
	"encoding/json"
	"encoding/pem"
	"io"
	"net/http"

	"golang.org/x/crypto/ssh"
	"pault.ag/go/sshsig"
)

func (bot *Bot) webappSubmit(w http.ResponseWriter, r *http.Request) {
	var requestBody map[string]interface{}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	username, ok := requestBody["username"].(string)
	if !ok || username == "" {
		http.Error(w, "Missing or invalid phone parameter", http.StatusBadRequest)
		return
	}

	signature, ok := requestBody["signature"].(string)
	if !ok || signature == "" {
		http.Error(w, "Missing or invalid signature parameter", http.StatusBadRequest)
		return
	}

	publicKey, err := getPublicKey(username)
	if err != nil {
		http.Error(w, "Failed to get public key", http.StatusInternalServerError)
		return
	}

	userID := bot.webappGetUserID(w, r)
	data := bot.db.GetChallengeCode(userID)
	if err := verifySignature(data, bot.namespace, signature, publicKey); err != nil {
		http.Error(w, "Invalid signature", http.StatusUnauthorized)
		return
	}

	if err := verifyQualification(username); err != nil {
		http.Error(w, "Invalid qualification", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"message": "Signature verified successfully",
	}
	json.NewEncoder(w).Encode(response)
}

func getPublicKey(username string) ([]byte, error) {
	req, err := http.NewRequest("GET", "https://github.com/"+username+"/keys", nil)
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func verifySignature(data, namespace, signature string, pubKey []byte) error {
	block, _ := pem.Decode([]byte(signature))
	sig, err := sshsig.ParseSignature(block.Bytes)
	if err != nil {
		return err
	}
	cHash, err := sig.HashAlgorithm.Hash()
	if err != nil {
		return err
	}
	h := cHash.New()
	h.Write([]byte(data))
	hash := h.Sum(nil)
	publicKey, err := ssh.ParsePublicKey(pubKey)
	if err != nil {
		return err
	}
	err = sshsig.Verify(publicKey, []byte(namespace), sig.HashAlgorithm, hash, sig)
	if err != nil {
		return err
	}
	return nil
}

func verifyQualification(username string) error {

	return nil
}
