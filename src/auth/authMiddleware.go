package auth

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

// Middleware to protect routes
func AuthMiddleware(next http.Handler, store *sessions.CookieStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := store.Get(r, "session-name")
		if err != nil {
			log.Printf("Error retrieving session: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Check if user is authenticated
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		userID, ok := session.Values["user_id"].(int)
		if !ok || userID == 0 {
			log.Println("Invalid or missing user ID in session")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userName, ok := session.Values["username"].(string)
		if !ok || userName == "" {
			log.Println("Invalid or missing username in session")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		ctx = context.WithValue(ctx, UserNameKey, userName)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
