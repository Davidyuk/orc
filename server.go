package main

import (
    "flag"
    "github.com/orc/db"
    "github.com/orc/resources"
    "github.com/orc/router"
    "github.com/orc/mvc/controllers"
    "log"
    "net/http"
    "os"
)

func main() {
    log.Println("Server started.")

    testData := flag.Bool("test-data", false, "to load test data")
    flag.Parse()

    db.Init()
    controllers.CreateRegistrationEvent()

    if *testData == true {
        resources.Load()
    }

    base := new(controllers.BaseController)
    base.Index().LoadContestsFromCats()

    http.Handle("/", new(router.FastCGIServer))
    http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./static/js"))))
    http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./static/css"))))
    http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("./static/img"))))

    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Println("Error listening: ", err.Error())
        os.Exit(1)
    }
}
