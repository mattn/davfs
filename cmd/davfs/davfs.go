package main

import (
	"flag"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/mattn/go-colorable"

	"github.com/mattn/davfs"
	_ "github.com/mattn/davfs/plugin/file"
	_ "github.com/mattn/davfs/plugin/memory"
	_ "github.com/mattn/davfs/plugin/mysql"
	_ "github.com/mattn/davfs/plugin/postgres"
	_ "github.com/mattn/davfs/plugin/sqlite3"
	"golang.org/x/net/webdav"
)

var (
	addr   = flag.String("addr", ":9999", "server address")
	driver = flag.String("driver", "file", "database driver")
	source = flag.String("source", ".", "database connection string")
	cred   = flag.String("cred", "", "credential for basic auth")
	create = flag.Bool("create", false, "create filesystem")
)

func errorString(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}

func main() {
	flag.Parse()

	log.SetOutput(colorable.NewColorableStderr())

	if *create {
		err := davfs.CreateFS(*driver, *source)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}
	fs, err := davfs.NewFS(*driver, *source)
	if err != nil {
		log.Fatal(err)
	}

	dav := &webdav.Handler{
		FileSystem: fs,
		LockSystem: webdav.NewMemLS(),
		Logger: func(r *http.Request, err error) {
			litmus := r.Header.Get("X-Litmus")
			if len(litmus) > 19 {
				litmus = litmus[:16] + "..."
			}

			switch r.Method {
			case "COPY", "MOVE":
				dst := ""
				if u, err := url.Parse(r.Header.Get("Destination")); err == nil {
					dst = u.Path
				}
				log.Printf("%-18s %s %s %s",
					color.GreenString(r.Method),
					r.URL.Path,
					dst,
					color.RedString(errorString(err)))
			default:
				log.Printf("%-18s %s %s",
					color.GreenString(r.Method),
					r.URL.Path,
					color.RedString(errorString(err)))
			}
		},
	}

	var handler http.Handler
	if *cred != "" {
		token := strings.SplitN(*cred, ":", 2)
		if len(token) != 2 {
			flag.Usage()
			return
		}
		user, pass := token[0], token[1]
		handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			username, password, ok := r.BasicAuth()
			if !ok || username != user || password != pass {
				w.Header().Set("WWW-Authenticate", `Basic realm="davfs"`)
				http.Error(w, "authorization failed", http.StatusUnauthorized)
				return
			}
			dav.ServeHTTP(w, r)
		})
	} else {
		handler = dav
	}

	log.Print(color.CyanString("Server started %v", *addr))
	http.Handle("/", handler)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
