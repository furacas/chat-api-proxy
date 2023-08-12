package chatgpt

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

var accounts []Account

type Account struct {
	Email    string `json:"username"`
	Password string `json:"password"`
}

// Read accounts.txt and create a list of accounts
func readAccounts() {
	accounts = []Account{}
	// Read accounts.txt and create a list of accounts
	if _, err := os.Stat("config/accounts.txt"); err == nil {
		// Each line is a proxy, put in proxies array
		file, _ := os.Open("config/accounts.txt")
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			// Split by :
			line := strings.Split(scanner.Text(), ":")
			// Create an account
			account := Account{
				Email:    line[0],
				Password: line[1],
			}
			// Append to accounts
			accounts = append(accounts, account)
		}
	}
}

func updateToken() {
	readAccounts()

	token_list := []string{}
	// Loop through each account
	for _, account := range accounts {

		println("Updating access token for " + account.Email)
		authenticator := NewAuthenticator(account.Email, account.Password, "http://127.0.0.1:7890")
		err := authenticator.Begin()
		if err != nil {
			// println("Error: " + err.Details)
			println("Location: " + err.Location)
			println("Status code: " + fmt.Sprint(err.StatusCode))
			println("Details: " + err.Details)
			println("Embedded error: " + err.Error.Error())
			return
		}
		access_token := authenticator.GetAccessToken()
		token_list = append(token_list, access_token)
		println("Success!")
		// Write authenticated account to authenticated_accounts.txt
		f, go_err := os.OpenFile("authenticated_accounts.txt", os.O_APPEND|os.O_WRONLY, 0600)
		if go_err != nil {
			continue
		}
		defer f.Close()
		if _, go_err = f.WriteString(account.Email + ":" + account.Password + "\n"); go_err != nil {
			continue
		}
		// Remove accounts.txt
		os.Remove("accounts.txt")
		// Create accounts.txt
		f, go_err = os.Create("accounts.txt")
		if go_err != nil {
			continue
		}
		defer f.Close()
		// Remove account from accounts
		accounts = accounts[1:]
		// Write unauthenticated accounts to accounts.txt
		for _, acc := range accounts {
			// Check if account is authenticated
			if acc.Email == account.Email {
				continue
			}
			if _, go_err = f.WriteString(acc.Email + ":" + acc.Password + "\n"); go_err != nil {
				continue
			}
		}
	}

	tokenManager = initializeTokens(token_list)

	time.AfterFunc(1.2096e15, updateToken)
}

func init() {
	if os.Getenv("CHATGPT_ENABLED") != "false" {
		go updateToken()
	}

}
