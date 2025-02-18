package user

import (
	"context"
	"encoding/json"
	"fmt"
	"icealpha/internal/router"

	"net/http"
)

const QUERYBOILERPLATE = "Give the result of the following problem: %s\n Give the result in the first line and the explanation in the following lines"

// POST(problem: multipart[image]) -> Json(content: string)
func HandleSolveInputImage(pattern string, rtr *router.Router) {

	rtr.R.Post(pattern, func(w http.ResponseWriter, r *http.Request) {

		multiPartFile := r.MultipartForm.File["problem"][0]

		if multiPartFile.Size > (1024 * 1024 * 20) { // 20 mb

			http.Error(w, "file size exceeded", http.StatusBadRequest)
			return

		}

		// add checks to filter out non-image files

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

		responseChannel, err := rtr.S.LLMClient.StreamResponse(context.Background(), fmt.Sprintf(QUERYBOILERPLATE, latex))
		if err != nil {

			http.Error(w, "Could not think at the moment. Try again later")
			return
		}

		//		for msg := range responseChannel {
		//
		//		}

	})

}

// POST :: Json(query: string) -> Json(content: string)
func HandleSolveTextInput(pattern string, rtr *router.Router) {

	rtr.R.Post(pattern, func(w http.ResponseWriter, r *http.Request) {

		defer r.Body.Close()

		queryStruct := struct {
			Query string `json:"query"`
		}{}

		if err := json.NewDecoder(r.Body).Decode(&queryStruct); err != nil {

			http.Error(w, "error decoding json query input", http.StatusBadRequest)
			return

		}

		// convert queryStruct.Query to the result

	})

}
