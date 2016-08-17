package main

import (
    "flag"
    "log"
    "net/http"
    "database/sql"
    "github.com/klenin/orc/db"
    "github.com/klenin/orc/initial"
    "github.com/klenin/orc/router"
    "github.com/klenin/orc/mvc/controllers"
    "github.com/klenin/orc/config"
    "strings"
)

var err error

type stopOnNotFoundResponseWriter struct {
    http.ResponseWriter
    IsNotFound bool
}

func (ww *stopOnNotFoundResponseWriter) WriteHeader(status int) {
    if status == http.StatusNotFound {
        ww.IsNotFound = true
        for k := range ww.ResponseWriter.Header() {
            delete(ww.ResponseWriter.Header(), k)
        }
    } else {
        ww.ResponseWriter.WriteHeader(status)
    }
}

func (ww *stopOnNotFoundResponseWriter) Write(p []byte) (int, error) {
    if ww.IsNotFound {
        return len(p), nil
    }
    return ww.ResponseWriter.Write(p)
}

func main() {
    db.DB, err = sql.Open("postgres", config.GetValue("DATABASE_URL"))
    defer db.DB.Close()

    if err != nil {
        log.Fatalln("Error DB open:", err.Error())
    }

    if err = db.DB.Ping(); err != nil {
        log.Fatalln("Error DB ping:", err.Error())
    }

    log.Println("Connected to DB")

    testData := flag.Bool("test-data", false, "to load test data")
    resetDB := flag.Bool("reset-db", false, "reset the database")
    flag.Parse()

    initial.Init(*resetDB, *testData)

    // base := new(controllers.BaseController)
    // base.Index().LoadContestsFromCats()

    http.Handle("/", new(router.FastCGIServer))
    http.HandleFunc("/wellcometoprofile/", controllers.WellcomeToProfile)

    fileServer := http.FileServer(http.Dir("./static"))
    http.Handle("/js/", fileServer)
    http.Handle("/css/", fileServer)
    http.Handle("/img/", fileServer)
    http.Handle("/vendor/", fileServer)
    http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
        r.URL.Path = strings.Replace(r.URL.Path, "static", "frontend", 1)
        ww := &stopOnNotFoundResponseWriter{ResponseWriter: w}
        fileServer.ServeHTTP(ww, r)
        if ww.IsNotFound {
            http.ServeFile(w, r, "./static/frontend/index.html")
        }
    })

    addr := config.GetValue("HOSTNAME") + ":" + config.GetValue("PORT")
    log.Println("Server listening on", addr)
    log.Fatalln("Error listening:", http.ListenAndServe(addr, nil))
}
