package user

import (
	"icealpha/internal/router"

	"net/http"
)

func HandleSolveInputImage(pattern string, rtr *router.Router) {

	rtr.R.Post(pattern, func(w http.ResponseWriter, r *http.Request) {

		file := r.MultipartForm.File["problem"][0]

		if file.Size > (1024 * 1024 * 20) { // 20 mb

			http.Error(w, "file size exceeded", http.StatusBadRequest)
			return

		}

		// add checks to filter out non-image files

	})

}
