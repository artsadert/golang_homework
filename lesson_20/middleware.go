package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type contextKey string

const (
	UserIDKey     contextKey = "userID"
	UserRoleKey   contextKey = "userRole"
	ResourceIDKey contextKey = "resourceID"
)

// AuthorizeUserAccess checks if the authenticated user has access to the requested resource
func AuthorizeUserAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get user from context (assuming you have authentication middleware before this)
		// This could come from a JWT token, session, etc.
		authenticatedUser, ok := r.Context().Value(UserIDKey).(*User)
		if !ok {
			// If no user in context, try to get from request headers/cookies
			// This is just an example - implement your actual auth logic
			user, err := authenticateUser(r)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			authenticatedUser = user
		}

		// Get the resource ID from the URL
		resourceID := chi.URLParam(r, "id")

		// For routes without ID parameter (like POST /users, GET /users)
		if resourceID == "" {
			// For collection endpoints, check if user has permission to access
			if !hasCollectionAccess(authenticatedUser, r.Method) {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
			return
		}

		// Check if the authenticated user owns this resource or has admin privileges
		if !canAccessResource(authenticatedUser, resourceID, r.Method) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// Add resource ID to context for handlers to use
		ctx := context.WithValue(r.Context(), ResourceIDKey, resourceID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// authenticateUser extracts user from request (implement your actual auth logic)
func authenticateUser(r *http.Request) (*User, error) {
	// Example: Get user from JWT token in Authorization header
	token := r.Header.Get("Authorization")
	if token == "" {
		return nil, http.ErrNoCookie
	}

	// Validate token and get user from database
	// This is where you'd implement your actual authentication logic
	// user := &User{
	// 	ID:   1,      // Example user ID
	// }
	//
	userCommand := LoginUserCommand{}

	err := json.NewDecoder(r.Body).Decode(&userCommand)
	if err != nil {
		return nil, http.ErrSchemeMismatch
	}

	user := NewUser(userCommand.Name)

	return user, nil
}

// hasCollectionAccess checks if user can access collection endpoints
func hasCollectionAccess(user *User, method string) bool {
	// Admins can access all collection endpoints
	// Regular users can only GET collections (list users) and POST (create new user)
	switch method {
	case http.MethodGet, http.MethodPost:
		return true
	default:
		return false
	}
}

// canAccessResource checks if user can access a specific resource
func canAccessResource(user *User, resourceID string, method string) bool {
	return user.Id == resourceID
}

// Helper function to get user ID from context
func GetUserIDFromContext(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value(UserIDKey).(int)
	return userID, ok
}

// Helper function to get resource ID from context
func GetResourceIDFromContext(ctx context.Context) (int, bool) {
	resourceID, ok := ctx.Value(ResourceIDKey).(int)
	return resourceID, ok
}
