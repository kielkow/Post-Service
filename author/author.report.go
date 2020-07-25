package author

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"path"
	"time"
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
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		authors, err := searchAuthorData(authorFilter)

		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		t := template.New("report_author.gotmpl").Funcs(template.FuncMap{"mod": func(i, x int) bool { return i%x == 0 }})
		t, err = t.ParseFiles(path.Join("templates", "report_author.gotmpl"))

		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
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
