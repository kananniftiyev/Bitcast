package auth

import (
	"context"
	"encoding/json"
	"fileguard/internal/db"
	"fileguard/utils"
	"fmt"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"log"
	"net/http"
	"os"
	"sync"
)

// Build Google OAuth2.

var GoogleOauthConfig *oauth2.Config

func init() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Initialize Google OAuth2 configuration
	GoogleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/callback",
		ClientID:     os.Getenv("GOOGLE_AUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_AUTH_SECRET_KEY"),
		Scopes:       []string{"openid", "email", "profile"},
		Endpoint:     google.Endpoint,
	}
}

type UserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
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

	var userInfo UserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		http.Error(w, "Unable to parse user info response", http.StatusInternalServerError)
		fmt.Printf("Unable to parse user info response: %v\n", err)
		return
	}

	db, err := db.NewDatabase()
	if err != nil {
		fmt.Println(err)
	}

	err = db.CreateNewUser(userInfo.GivenName, userInfo.Email, "123")
	if err != nil {
		log.Fatal(err)
		return
	}

	createdToken := utils.GenerateToken(userInfo.ID)
	err = utils.SaveToken(createdToken)
	if err != nil {
		log.Fatal(err)
	}

	// Respond to the request
	fmt.Fprintf(w, "User info obtained successfully")
}

// TODO: Fix This fully.
func LoginViaGoogle() {
	token, err := utils.LoadToken()
	if err == nil {
		log.Fatal("Already Logged in")
	}
	_ = token
	// Set up HTTP storage to handle callback
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
			log.Fatalf("Failed to start HTTP storage: %v", err)
		}
	}()

	// Get the authorization URL
	authURL := GoogleOauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	fmt.Printf("Go to the following link in your browser: \n%v\n", authURL)
}
