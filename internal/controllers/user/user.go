package user

import (
	"encoding/json"
	"icealpha/internal/router"

	"net/http"
)

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

		// latex to answer

		response := map[string]string{

			"content": "tbd",
		}

		if err = json.NewEncoder(w).Encode(response); err != nil {

			http.Error(w, "Error writing final response", http.StatusInternalServerError)
			return

		}

	})

}
