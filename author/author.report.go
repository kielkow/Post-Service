package author

import (
	"bytes"
	"encoding/json"
	"html/template"
	"net/http"
	"path"
	"time"

	"github.com/kielkow/Post-Service/apperror"
)

// ReportFilter struct
type ReportFilter struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

func handleAuthorReport(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var authorFilter ReportFilter

		err := json.NewDecoder(r.Body).Decode(&authorFilter)

		if err != nil {
			error := apperror.GenerateError(400, err.Error())

			w.WriteHeader(http.StatusBadRequest)
			w.Write(error)
			return
		}

		authors, err := searchAuthorData(authorFilter)

		if err != nil {
			error := apperror.GenerateError(500, err.Error())

			w.WriteHeader(http.StatusInternalServerError)
			w.Write(error)
			return
		}

		t := template.New("report_author.gotmpl").Funcs(template.FuncMap{"mod": func(i, x int) bool { return i%x == 0 }})
		t, err = t.ParseFiles(path.Join("templates", "report_author.gotmpl"))

		if err != nil {
			error := apperror.GenerateError(500, err.Error())

			w.WriteHeader(http.StatusInternalServerError)
			w.Write(error)
			return
		}

		var tmpl bytes.Buffer

		if len(authors) > 0 {
			err = t.Execute(&tmpl, authors)
		} else {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		rdr := bytes.NewReader(tmpl.Bytes())

		w.Header().Set("Content-Disposition", "Attachement")

		http.ServeContent(w, r, "report_author.html", time.Now(), rdr)

	case http.MethodOptions:
		return

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
