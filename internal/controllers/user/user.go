package user

import (
	"context"
	"encoding/json"
	"fmt"
	"icealpha/internal/controllers/auth"
	"icealpha/internal/router"

	"bytes"
	"io"
	"net/http"
)

const QueryBoilerplate = "Give the result of the following problem: %s\n Give the result in the first line and the explanation in the following lines"

func validateImageFile(r *http.Request) (bool, error) {
	// Read first 512 bytes to determine the content type
	buffer := make([]byte, 512)
	_, err := r.Body.Read(buffer)
	if err != nil && err != io.EOF {
		return false, err
	}

	// Reset the body to be read again
	r.Body = io.NopCloser(bytes.NewBuffer(append(buffer, make([]byte, 0)...)))

	// Check the content type
	contentType := http.DetectContentType(buffer)

	validImageTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}
	// List of allowed image MIME types

	return validImageTypes[contentType], nil
}

// POST(problem: multipart[image]) -> Json(content: string)
func HandleSolveInputImage(rtr *router.Router) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		userSession, err := rtr.S.CookieStore.Get(
			r, "usersession",
		)

		if err != nil {

			http.Error(w, "User session expired", http.StatusBadRequest)
			return

		}

		userEmail := userSession.Values["email"].(string)

		credits, ok := rtr.S.CreditCache.Load(userEmail)
		if !ok {

			http.Error(w, "internal server error occureed", http.StatusInternalServerError)
			return

		}

		creditsUint64 := credits.(uint64)

		// handle case where the user has no credits
		if creditsUint64 == 0 {

			http.Error(w, "user has no credits left", http.StatusBadRequest)
			return

		}

		multiPartFile := r.MultipartForm.File["problem"][0]

		if multiPartFile.Size > (1024 * 1024 * 20) { // 20 mb

			http.Error(w, "file size exceeded", http.StatusBadRequest)
			return

		}

		if ok, err := validateImageFile(r); err != nil || !ok {

			http.Error(w, "invalid image file", http.StatusBadRequest)
			return

		}

		file, err := multiPartFile.Open()
		if err != nil {

			http.Error(w, "error opening submitted image", http.StatusInternalServerError)
			return

		}

		latex, err := rtr.S.ImgLatex.ImageToLatex(file)
		if err != nil {

			http.Error(w, "could not extrapolate information from given file", http.StatusBadRequest)
			return

		}

		responseChannel, err := rtr.S.LLMClient.StreamResponse(context.Background(), fmt.Sprintf(QueryBoilerplate, latex))
		if err != nil {

			http.Error(w, "Could not think at the moment. Try again later", http.StatusInternalServerError)
			return
		}

		flusher, ok := w.(http.Flusher)
		if !ok {

			http.Error(w, "Could not create a streaming response", http.StatusInternalServerError)
			return

		}

		responseStruct := struct {
			Token string `json:"token"`
		}{}

		for llmResponseToken := range responseChannel {

			responseStruct.Token = llmResponseToken

			if err := json.NewEncoder(w).Encode(&responseStruct); err != nil {

				http.Error(w, "error while streaming response", http.StatusInternalServerError)
				return

			}

			// flush changes to the stream
			flusher.Flush()

		}

		rtr.S.CreditCache.Store(userEmail, creditsUint64-1)

	}

}

// POST :: Json(query: string) -> Json(content: string)
func HandleSolveTextInput(rtr *router.Router) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		userSession, err := rtr.S.CookieStore.Get(
			r, "usersession",
		)

		if err != nil {

			http.Error(w, "User session expired", http.StatusBadRequest)
			return

		}

		userEmail := userSession.Values["email"].(string)

		credits, ok := rtr.S.CreditCache.Load(userEmail)
		if !ok {

			http.Error(w, "internal server error occureed", http.StatusInternalServerError)
			return

		}

		creditsUint64 := credits.(uint64)

		// handle case where the user has no credits
		if creditsUint64 == 0 {

			http.Error(w, "user has no credits left", http.StatusBadRequest)
			return

		}
		defer r.Body.Close()

		queryStruct := struct {
			Query string `json:"query"`
		}{}

		if err := json.NewDecoder(r.Body).Decode(&queryStruct); err != nil {

			http.Error(w, "error decoding json query input", http.StatusBadRequest)
			return

		}

		responseChannel, err := rtr.S.LLMClient.StreamResponse(context.Background(), fmt.Sprintf(QueryBoilerplate, queryStruct.Query))
		if err != nil {

			http.Error(w, "Could not think at the moment. Try again later", http.StatusInternalServerError)
			return
		}

		flusher, ok := w.(http.Flusher)
		if !ok {

			http.Error(w, "Could not create a streaming response", http.StatusInternalServerError)
			return

		}

		responseStruct := struct {
			Token string `json:"token"`
		}{}

		for llmResponseToken := range responseChannel {

			responseStruct.Token = llmResponseToken

			if err := json.NewEncoder(w).Encode(&responseStruct); err != nil {

				http.Error(w, "error while streaming response", http.StatusInternalServerError)
				return

			}

			flusher.Flush()

		}

		rtr.S.CreditCache.Store(userEmail, creditsUint64-1)

	}

}

// user authentication middleware, only checks session for now
func AuthMiddleware(next http.HandlerFunc, rtr *router.Router) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if !auth.CheckSessionExists(r, rtr.S.CookieStore) {

			http.Redirect(w, r, "/api", http.StatusTemporaryRedirect)
			return

		}
		next.ServeHTTP(w, r)

	}

}
