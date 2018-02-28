package writers

import (
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/blevesearch/bleve"
	bolt "github.com/boltdb/bolt"
	"github.com/jtarchie/syslog/pkg/log"
)

type Server struct {
	port       int
	index      bleve.Index
	httpServer *http.Server
	db         *bolt.DB
}

type doc struct {
	Version   int
	Priority  int
	Timestamp time.Time
	Hostname  string
	AppName   string
	ProcID    string
	MsgID     string
	Message   string
}

func NewServer(port int) *Server {
	index, err := bleve.NewMemOnly(bleve.NewIndexMapping())
	if err != nil {
		log.Fatalf("could not start indexer: %s", err)
	}

	tmpFile, err := ioutil.TempFile("", "messages")
	if err != nil {
		log.Fatalf("could not create db: %s", err)
	}

	db, err := bolt.Open(tmpFile.Name(), 0600, &bolt.Options{Timeout: 10 * time.Second})
	if err != nil {
		log.Fatalf("could not start db: %s", err)
	}

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("messages"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	return &Server{
		index: index,
		port:  port,
		db:    db,
	}
}

func (s *Server) Write(l *syslog.Log) error {
	id := strconv.FormatInt(time.Now().UnixNano(), 10)

	err := s.index.Index(id, doc{
		Version:   l.Version(),
		Priority:  l.Priority(),
		Timestamp: l.Timestamp(),
		Hostname:  l.Hostname(),
		AppName:   l.Appname(),
		ProcID:    l.ProcID(),
		MsgID:     l.MsgID(),
		Message:   l.Message(),
	})

	if err != nil {
		return err
	}

	return s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("messages"))
		return bucket.Put([]byte(id), []byte(l.String()))
	})
}

func (s *Server) Start() error {
	log.Printf("web: starting search index")

	mux := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		html := fmt.Sprintf(`
		<html>
		<head>
			<link href="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-Gn5384xqQ1aoWXA+058RXPxPg6fy4IWvTNh0E263XmFcJlSAwiGgFAW/dAiS6JXm" crossorigin="anonymous">
		</head>
		<body>
		<div class="container">
		<form type="GET" action="/">
			<div class="form-group row">
				<label for="q" class="col-sm-2 col-form-label">Query</label>
				<input type="text" id="q" name="q" value="%s" class="col-sm-9">
			</div>
		</form>
		`, html.EscapeString(r.URL.Query().Get("q")))

		query := bleve.NewQueryStringQuery(r.URL.Query().Get("q"))
		search := bleve.NewSearchRequest(query)
		search.Size = 1000
		search.SortBy([]string{"Timestamp"})

		result, err := s.index.Search(search)
		if err != nil {
			html += fmt.Sprintf("Something went wrong: %s", err)
		} else {
			s.db.View(func(tx *bolt.Tx) error {
				bucket := tx.Bucket([]byte("messages"))

				for _, hit := range result.Hits {
					value := bucket.Get([]byte(hit.ID))
					html += fmt.Sprintf("<p>%s</p>", value)
				}
				return nil
			})
		}
		html += `</div></body></html>`
		w.Write([]byte(html))
	})

	log.Printf("web: starting webserver on port %d", s.port)
	s.httpServer = &http.Server{Addr: fmt.Sprintf("localhost:%d", s.port), Handler: mux}
	return s.httpServer.ListenAndServe()
}

func (s *Server) Addr() string {
	return s.httpServer.Addr
}
