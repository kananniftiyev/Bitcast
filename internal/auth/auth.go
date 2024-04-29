package auth

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

// Build Google OAuth2.

var GoogleOauthConfig = &oauth2.Config{
	RedirectURL:  "http://localhost:8080/callback",
	ClientID:     "1004356047315-otq33mfc9opfgkcp7mjlv3lheluuuu3p.apps.googleusercontent.com",
	ClientSecret: "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",
	Scopes:       []string{"openid", "email", "profile"},
	Endpoint:     google.Endpoint,
}

type UserInfo struct {
	Email    string
	Name     string
	Picture  string
	Verified bool
}

func handleCallback(w http.ResponseWriter, r *http.Request, wg *sync.WaitGroup) {
	defer wg.Done()

	ctx := context.Background()

	// Exchange the authorization code for a token
	code := r.URL.Query().Get("code")
	token, err := GoogleOauthConfig.Exchange(ctx, code)
	if err != nil {
		http.Error(w, "Unable to retrieve token from web", http.StatusInternalServerError)
		log.Fatalf("Unable to retrieve token from web: %v", err)
		return
	}

	// Create a client with the retrieved token
	client := GoogleOauthConfig.Client(ctx, token)

	// Make a request using the authenticated client
	resp, err := client.Get("https://www.googleapis.com/oauth2/v1/userinfo")
	if err != nil {
		http.Error(w, "Unable to retrieve user info", http.StatusInternalServerError)
		log.Fatalf("Unable to retrieve user info: %v", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Unable to read response body", http.StatusInternalServerError)
		log.Fatalf("Unable to read response body: %v", err)
		return
	}

	fmt.Println("Response from userinfo endpoint:")
	fmt.Println(string(body))

	// Respond to the request
	fmt.Fprintf(w, "User info obtained successfully")
}

func LoginViaGoogle() {
	// Set up HTTP server to handle callback
	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Wait()

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		handleCallback(w, r, &wg)
		http.DefaultServeMux.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
			wg.Done()
		})
	})

	srv := &http.Server{Addr: ":8080"}
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// Get the authorization URL
	authURL := GoogleOauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	fmt.Printf("Go to the following link in your browser: \n%v\n", authURL)
}
