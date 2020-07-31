package post

import (
	"bytes"
	"encoding/json"
	"html/template"
	"net/http"
	"path"
	"time"

	"github.com/kielkow/Post-Service/shared/providers/apperror"
)

// ReportFilter struct
type ReportFilter struct {
	Author      string `json:"author"`
	Description string `json:"description"`
}

func handlePostReport(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var postFilter ReportFilter

		err := json.NewDecoder(r.Body).Decode(&postFilter)

		if err != nil {
			error := apperror.GenerateError(400, err.Error())

			w.WriteHeader(http.StatusBadRequest)
			w.Write(error)
			return
		}

		posts, err := searchPostData(postFilter)

		if err != nil {
			error := apperror.GenerateError(500, err.Error())

			w.WriteHeader(http.StatusInternalServerError)
			w.Write(error)
			return
		}

		t := template.New("report.gotmpl").Funcs(template.FuncMap{"mod": func(i, x int) bool { return i%x == 0 }})
		t, err = t.ParseFiles(path.Join("templates", "report.gotmpl"))

		if err != nil {
			error := apperror.GenerateError(500, err.Error())

			w.WriteHeader(http.StatusInternalServerError)
			w.Write(error)
			return
		}

		var tmpl bytes.Buffer

		if len(posts) > 0 {
			err = t.Execute(&tmpl, posts)
		} else {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		rdr := bytes.NewReader(tmpl.Bytes())

		w.Header().Set("Content-Disposition", "Attachement")

		http.ServeContent(w, r, "report.html", time.Now(), rdr)

	case http.MethodOptions:
		return

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
