package main

import (
	"encoding/json"
	"encoding/pem"
	"errors"
	"io"
	"net/http"

	"golang.org/x/crypto/ssh"
	"pault.ag/go/sshsig"
)

func (bot *Bot) webappSubmit(w http.ResponseWriter, r *http.Request) {
	var requestBody map[string]interface{}
	response := map[string]string{}
	userID := bot.webappGetUserID(w, r)
	data := bot.db.GetChallengeCode(userID)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		response["message"] = "Invalid request body"
	} else if username, ok := requestBody["username"].(string); !ok || username == "" {
		response["message"] = "Missing or invalid username parameter\n用户名无效"
	} else if publicKey, err := getPublicKey(username); err != nil {
		response["message"] = "Failed to retrieve public key\n获取公钥失败"
	} else if signature, ok := requestBody["signature"].(string); !ok || signature == "" {
		response["message"] = "invalid signature\n签名无效"
	} else if err := verifySignature(data, bot.namespace, signature, publicKey); err != nil {
		response["message"] = "Invalid signature\n签名无效"
	} else if message, err := bot.verifyQualification(username); err != nil {
		response["message"] = message
	} else {
		response["ok"] = "200"
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

func (bot *Bot) verifyQualification(username string) (string, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/users/"+username, nil)
	if err != nil {
		return "", err
	}
	defer req.Body.Close()

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return "", err
	}
	var GitHubInfo map[string]interface{}
	if err := json.Unmarshal(body, &GitHubInfo); err != nil {
		return "", err
	}
	githubID := GitHubInfo["id"].(int64)
	user := bot.db.GetUserByGithubID(githubID)
	if user.TelegramID != 0 {
		return "Already verified\n重复验证", errors.New("Already verified")
	}
	if user.Flag&FlagBanned != 0 {
		return "You have been permanently banned\n你已被永久封禁", errors.New("You have been permanently banned")
	}
	return "", nil
}
