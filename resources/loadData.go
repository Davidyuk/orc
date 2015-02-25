package resources

import (
    "github.com/orc/db"
    "github.com/orc/mvc/controllers"
    "github.com/orc/mvc/models"
    "io/ioutil"
    "math/rand"
    "strconv"
    "strings"
    "time"
)

const USER_COUNT = 20

var base = new(models.ModelManager)

func random(min, max int) int {
    rand.Seed(int64(time.Now().Second()))
    return rand.Intn(max-min) + min
}

func addDate(d, m, y int) string {
    return strconv.Itoa(d) + "-" + strconv.Itoa(m) + "-" + strconv.Itoa(y)
}

func addTime(h, m, s int) string {
    return strconv.Itoa(h) + ":" + strconv.Itoa(m) + ":" + strconv.Itoa(s)
}

func Load() {
    loadUsers()
    loadEvents()
    loadEventTypes()
    loadForms()
    loadParamTypes()
}

func loadUsers() {
    base := new(controllers.BaseController)
    for i := 0; i < USER_COUNT; i++ {
        rand.Seed(int64(i))
        result, reg_id := base.Handler().HandleRegister_("user"+strconv.Itoa(i), "secret"+strconv.Itoa(i), "user")
        if result == "ok" {
            eventsRegs := controllers.GetModel("events_regs")
            eventsRegs.LoadModelData(map[string]interface{}{"reg_id": reg_id, "event_id": 1})
            db.QueryInsert_(eventsRegs, "")
        }
    }
    result, reg_id := base.Handler().HandleRegister_("admin", "password", "admin")
    if result == "ok" {
        eventsRegs := controllers.GetModel("events_regs")
        eventsRegs.LoadModelData(map[string]interface{}{"reg_id": reg_id, "event_id": 1})
        db.QueryInsert_(eventsRegs, "")
    }
}

func loadEvents() {
    eventNames, _ := ioutil.ReadFile("./resources/event-name")
    subjectNames, _ := ioutil.ReadFile("./resources/subject-name")
    eventNameSource := strings.Split(string(eventNames), "\n")
    subjectNameSource := strings.Split(string(subjectNames), "\n")
    for i := 0; i < len(eventNameSource); i++ {
        rand.Seed(int64(i))
        eventName := strings.TrimSpace(eventNameSource[rand.Intn(len(eventNameSource))])
        eventName += " по дисциплине "
        eventName += "\"" + strings.TrimSpace(subjectNameSource[rand.Intn(len(subjectNameSource))]) + "\""
        dateStart := addDate(random(1894, 2014), random(1, 12), random(1, 28))
        dateFinish := addDate(random(1894, 2014), random(1, 12), random(1, 28))
        time := addTime(random(0, 11), random(1, 60), random(1, 60))
        params := []interface{}{eventName, dateStart, dateFinish, time, ""}
        entity := base.Events()
        db.QueryInsert("events", entity.GetColumnSlice(1), params, "")
    }
}

func loadEventTypes() {
    eventTypeNames, _ := ioutil.ReadFile("./resources/event-type-name")
    eventTypeNamesSourse := strings.Split(string(eventTypeNames), "\n")
    topicality := []bool{true, false}
    for i := 0; i < len(eventTypeNamesSourse); i++ {
        //rand.Seed(int64(i))
        eventTypeName := strings.TrimSpace(eventTypeNamesSourse[i])
        params := []interface{}{eventTypeName, "", topicality[rand.Intn(2)]}
        entity := base.EventTypes()
        db.QueryInsert("event_types", entity.GetColumnSlice(1), params, "")
    }
}

func loadForms() {
    formNames, _ := ioutil.ReadFile("./resources/form-name")
    formNamesSourse := strings.Split(string(formNames), "\n")
    for i := 0; i < len(formNamesSourse); i++ {
        formName := strings.TrimSpace(formNamesSourse[i])
        entity := base.Forms()
        db.QueryInsert("forms", entity.GetColumnSlice(1), []interface{}{formName}, "")
    }
}

func loadParamTypes() {
    paramTypesNames, _ := ioutil.ReadFile("./resources/param-type-name")
    paramTypesSourse := strings.Split(string(paramTypesNames), "\n")
    for i := 0; i < len(paramTypesSourse); i++ {
        paramType := strings.TrimSpace(paramTypesSourse[i])
        entity := base.ParamTypes()
        db.QueryInsert("param_types", entity.GetColumnSlice(1), []interface{}{paramType}, "")
    }
}
