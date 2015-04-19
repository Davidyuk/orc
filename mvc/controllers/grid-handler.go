package controllers

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "github.com/orc/db"
    "github.com/orc/mailer"
    "github.com/orc/mvc/models"
    "github.com/orc/sessions"
    "github.com/orc/utils"
    "math"
    "net/http"
    "strconv"
    "strings"
)

func (c *BaseController) GridHandler() *GridHandler {
    return new(GridHandler)
}

type GridHandler struct {
    Controller
}

func (this *GridHandler) GetSubTable() {
    if !sessions.CheackSession(this.Response, this.Request) {
        http.Redirect(this.Response, this.Request, "/", http.StatusUnauthorized)
        return
    }

    if !this.isAdmin() {
        return
    }

    request, err := utils.ParseJS(this.Request, this.Response)
    if err != nil {
        utils.SendJSReply(err.Error(), this.Response)
        return
    }

    model := GetModel(request["table"].(string))
    index, _ := strconv.Atoi(request["index"].(string))
    subModel := GetModel(model.GetSubTable(index))
    subModel.LoadWherePart(map[string]interface{}{model.GetSubField(): request["id"]})
    result := db.Select(subModel, subModel.GetColumns() )
    refFields, refData := GetModelRefDate(subModel)

    response, err := json.Marshal(map[string]interface{}{
        "data":      result,
        "name":      subModel.GetTableName(),
        "caption":   subModel.GetCaption(),
        "colnames":  subModel.GetColNames(),
        "columns":   subModel.GetColumns(),
        "reffields": refFields,
        "refdata":   refData})
    if utils.HandleErr("[GridHandler::GetSubTable] Marshal: ", err, this.Response) {
        return
    }

    fmt.Fprintf(this.Response, "%s", string(response))
}

func (this *GridHandler) Load(tableName string) {
    if !sessions.CheackSession(this.Response, this.Request) {
        http.Redirect(this.Response, this.Request, "/", http.StatusUnauthorized)
        return
    }

    if !this.isAdmin() {
        return
    }

    limit, err := strconv.Atoi(this.Request.PostFormValue("rows"))
    if utils.HandleErr("[GridHandler::Load]  limit Atoi: ", err, this.Response) {
        return
    }

    page, err := strconv.Atoi(this.Request.PostFormValue("page"))
    if utils.HandleErr("[GridHandler::Load] page Atoi: ", err, this.Response) {
        return
    }

    sidx := this.Request.FormValue("sidx")
    start := limit*page - limit

    model := GetModel(tableName)
    model.SetOrder(sidx)
    model.SetLimit(limit)
    model.SetOffset(start)

    rows := db.Select(model, model.GetColumns())
    count := db.SelectCount(tableName)

    var totalPages int
    if count > 0 {
        totalPages = int(math.Ceil(float64(count) / float64(limit)))
    } else {
        totalPages = 0
    }

    result := make(map[string]interface{}, 4)
    result["rows"] = rows
    result["page"] = page
    result["total"] = totalPages
    result["records"] = count

    utils.SendJSReply(result, this.Response)
}

func (this *GridHandler) Select(tableName string) {
    if !sessions.CheackSession(this.Response, this.Request) {
        http.Redirect(this.Response, this.Request, "/", http.StatusUnauthorized)
        return
    }

    if !this.isAdmin() {
        return
    }

    model := GetModel(tableName)
    refFields, refData := GetModelRefDate(model)

    this.Render([]string{"mvc/views/table.html"}, "table", Model{
        RefData:   refData,
        RefFields: refFields,
        TableName: model.GetTableName(),
        ColNames:  model.GetColNames(),
        Columns:   model.GetColumns(),
        Caption:   model.GetCaption(),
        Sub:       model.GetSub()})
}

func (this *GridHandler) Edit(tableName string) {
    if !sessions.CheackSession(this.Response, this.Request) {
        http.Error(this.Response, "", http.StatusUnauthorized)
        return
    }

    if !this.isAdmin() {
        return
    }

    model := GetModel(tableName)
    if model == nil {
        utils.HandleErr("[Grid-Handler::Edit] GetModel: invalid model", nil, this.Response)
        return
    }

    params := make(map[string]interface{}, len(model.GetColumns()))
    for i := 0; i < len(model.GetColumns()); i++ {
        params[model.GetColumnByIdx(i)] = this.Request.PostFormValue(model.GetColumnByIdx(i))
    }

    oper := this.Request.PostFormValue("oper")
    switch oper {
    case "edit":
        id, err := strconv.Atoi(this.Request.PostFormValue("id"))
        if utils.HandleErr("[Grid-Handler::Edit] strconv.Atoi: ", err, this.Response) {
            return
        }
        model.LoadModelData(params)
        model.LoadWherePart(map[string]interface{}{"id": id})
        err = db.QueryUpdate_(model).Scan()
        utils.HandleErr("", err, this.Response)
        break
    case "add":
        model.LoadModelData(params)
        db.QueryInsert_(model, "").Scan()
        break
    case "del":
        db.QueryDeleteByIds(tableName, this.Request.PostFormValue("id"))
        break
    }
}

func (this *GridHandler) ResetPassword() {
    if !sessions.CheackSession(this.Response, this.Request) {
        http.Redirect(this.Response, this.Request, "/", http.StatusUnauthorized)
        return
    }

    if !this.isAdmin() {
        return
    }

    this.Response.Header().Set("Access-Control-Allow-Origin", "*")
    this.Response.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
    this.Response.Header().Set("Content-type", "application/json")

    request, err := utils.ParseJS(this.Request, this.Response)
    if err != nil {
        utils.SendJSReply(err.Error(), this.Response)
        return
    }

    pass1 := request["pass1"].(string)
    pass2 := request["pass2"].(string)

    if !utils.MatchRegexp("^.{6,36}$", pass1) || !utils.MatchRegexp("^.{6,36}$", pass2) {
        utils.SendJSReply(map[string]interface{}{"result": "badPassword"}, this.Response)
        return
    } else if pass1 != pass2 {
        utils.SendJSReply(map[string]interface{}{"result": "differentPasswords"}, this.Response)
        return
    }

    id, err :=  strconv.Atoi(request["id"].(string))
    if utils.HandleErr("[Grid-Handler::ResetPassword] strconv.Atoi: ", err, this.Response) {
        return
    }

    user := GetModel("users")
    user.LoadWherePart(map[string]interface{}{"id": id})

    var salt string
    var enabled bool
    db.SelectRow(user, []string{"salt", "enabled"}).Scan(&salt, &enabled)

    user.GetFields().(*models.User).Enabled = enabled

    user.LoadModelData(map[string]interface{}{"pass": utils.GetMD5Hash(pass1 + salt)})
    db.QueryUpdate_(user).Scan()

    utils.SendJSReply(map[string]interface{}{"result": "ok"}, this.Response)
}

func (this *GridHandler) isAdmin() bool {
    var role string

    user_id := sessions.GetValue("id", this.Request)
    if user_id == nil {
        http.Redirect(this.Response, this.Request, "/", http.StatusUnauthorized)
        return false
    }

    user := GetModel("users")
    user.LoadWherePart(map[string]interface{}{"id": user_id})
    err := db.SelectRow(user, []string{"role"}).Scan(&role)
    if err != nil || role == "user" {
        http.Redirect(this.Response, this.Request, "/", http.StatusForbidden)
        return false
    }

    return role == "admin"
}

func (this *GridHandler) GetEventTypesByEventId() {
    if !sessions.CheackSession(this.Response, this.Request) {
        http.Redirect(this.Response, this.Request, "/", http.StatusUnauthorized)
        return
    }

    if !this.isAdmin() {
        return
    }

    request, err := utils.ParseJS(this.Request, this.Response)
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
    } else {
        event_id, err := strconv.Atoi(request["event_id"].(string))
        if utils.HandleErr("[GridHandler::GetEventTypesByEventId] event_id Atoi: ", err, this.Response) {
            return
        }

        query := `SELECT event_types.id, event_types.name FROM events_types
            INNER JOIN events ON events.id = events_types.event_id
            INNER JOIN event_types ON event_types.id = events_types.type_id
            WHERE events.id = $1;`
        result := db.Query(query, []interface{}{event_id})

        utils.SendJSReply(map[string]interface{}{"result": "ok", "data": result}, this.Response)
    }
}

func (this *GridHandler) ImportForms() {
    if !sessions.CheackSession(this.Response, this.Request) {
        http.Redirect(this.Response, this.Request, "/", http.StatusUnauthorized)
        return
    }

    if !this.isAdmin() {
        return
    }

    request, err := utils.ParseJS(this.Request, this.Response)
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
        return
    }

    event_id, err := strconv.Atoi(request["event_id"].(string))
    if utils.HandleErr("[GridHandler::ImportForms] event_id Atoi: ", err, this.Response) {
        return
    }

    for _, v := range request["event_types_ids"].([]interface{}) {
        println("event_types_ids: ", v)
        type_id, err := strconv.Atoi(v.(string))
        if err != nil {
            utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
            return
        }
        query := `SELECT events.id FROM events
            INNER JOIN events_types ON events_types.event_id = events.id
            INNER JOIN event_types ON event_types.id = events_types.type_id
            WHERE event_types.id=$1 AND events.id <> $2
            ORDER BY id DESC LIMIT 1;`

        eventResult := db.Query(query, []interface{}{type_id, event_id})

        query = `SELECT forms.id FROM forms
            INNER JOIN events_forms ON events_forms.form_id = forms.id
            INNER JOIN events ON events.id = events_forms.event_id
            WHERE events.id=$1;`

        formsResult := db.Query(query, []interface{}{int(eventResult[0].(map[string]interface{})["id"].(int64))})

        for i := 0; i < len(formsResult); i++ {
            form_id := int(formsResult[i].(map[string]interface{})["id"].(int64))
            eventsForms := GetModel("events_forms")
            eventsForms.LoadWherePart(map[string]interface{}{"event_id":  event_id, "form_id": form_id})
            var p int
            err := db.SelectRow(eventsForms, []string{"id"}).Scan(&p)
            if err != sql.ErrNoRows {
                continue
            }
            eventsForms.LoadModelData(map[string]interface{}{"event_id":  event_id, "form_id": form_id})
            db.QueryInsert_(eventsForms, "").Scan()
        }
    }

    utils.SendJSReply(map[string]interface{}{"result": "ok"}, this.Response)
}

func (this *GridHandler) GetPersonsByEventId() {
    if !sessions.CheackSession(this.Response, this.Request) {
        http.Redirect(this.Response, this.Request, "/", http.StatusUnauthorized)
        return
    }

    if !this.isAdmin() {
        return
    }

    request, err := utils.ParseJS(this.Request, this.Response)
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
    } else {
        event_id, err := strconv.Atoi(request["event_id"].(string))
        if utils.HandleErr("[GridHandler::GetPersonsByEventId] event_id Atoi: ", err, this.Response) {
            return
        }

        params := request["params_ids"].([]interface{})

        if len(params) == 0 {
            utils.SendJSReply(map[string]interface{}{"result": "Выберите параметры."}, this.Response)
            return
        }

        q := "SELECT params.name FROM params WHERE params.id in ("+strings.Join(db.MakeParams(len(params)), ", ")+") ORDER BY id"

        var caption []string
        for _, v := range db.Query(q, params) {
            caption = append(caption, v.(map[string]interface{})["name"].(string))
        }

        result := []interface{}{0: map[string]interface{}{"id": -1, "name": strings.Join(caption, " ")}}

        query := `SELECT reg_param_vals.reg_id as id, array_to_string(array_agg(param_values.value), ' ') as name
            FROM reg_param_vals
            INNER JOIN registrations ON registrations.id = reg_param_vals.reg_id
            INNER JOIN events ON events.id = registrations.event_id
            INNER JOIN param_values ON param_values.id = reg_param_vals.param_val_id
            INNER JOIN params ON params.id = param_values.param_id
            WHERE params.id in (` + strings.Join(db.MakeParams(len(params)), ", ")
        query += ") AND events.id = $" + strconv.Itoa(len(params)+1) + " GROUP BY reg_param_vals.reg_id ORDER BY reg_param_vals.reg_id;"

        data := db.Query(query, append(params, event_id))
        utils.SendJSReply(map[string]interface{}{"result": "ok", "data": append(result, data...)}, this.Response)
    }
}

func (this *GridHandler) GetParamsByEventId() {
    if !sessions.CheackSession(this.Response, this.Request) {
        http.Redirect(this.Response, this.Request, "/", http.StatusUnauthorized)
        return
    }

    if !this.isAdmin() {
        return
    }

    request, err := utils.ParseJS(this.Request, this.Response)
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
    } else {
        event_id, err := strconv.Atoi(request["event_id"].(string))
        if utils.HandleErr("[GridHandler::GetParamsByEventId] event_id Atoi: ", err, this.Response) {
            return
        }

        query := `SELECT DISTINCT params.id, params.name
            FROM reg_param_vals
            INNER JOIN param_values ON param_values.id = reg_param_vals.param_val_id
            INNER JOIN params ON params.id = param_values.param_id
            INNER JOIN registrations ON registrations.id = reg_param_vals.reg_id
            INNER JOIN events ON events.id = registrations.event_id
            WHERE events.id = $1 ORDER BY params.id;`

        result := db.Query(query, []interface{}{event_id})

        utils.SendJSReply(map[string]interface{}{"result": "ok", "data": result}, this.Response)
    }
}

func (this *GridHandler) GetPersonRequest() {
    if !sessions.CheackSession(this.Response, this.Request) {
        http.Redirect(this.Response, this.Request, "/", http.StatusUnauthorized)
        return
    }

    // if !this.isAdmin() {
    //     return
    // }

    request, err := utils.ParseJS(this.Request, this.Response)
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
    } else {
        face_id, err := strconv.Atoi(request["face_id"].(string))
        if utils.HandleErr("[GridHandler::GetPersonRequest] reg_id Atoi: ", err, this.Response) {
            return
        }

        query := `SELECT param_values.id, params.name, param_values.value, param_types.name as type FROM param_values
            INNER JOIN params ON params.id = param_values.param_id
            INNER JOIN param_types ON param_types.id = params.param_type_id
            INNER JOIN reg_param_vals ON reg_param_vals.param_val_id = param_values.id
            INNER JOIN registrations ON registrations.id = reg_param_vals.reg_id
            INNER JOIN faces ON registrations.face_id = faces.id
            WHERE faces.id=$1`

        data := db.Query(query, []interface{}{face_id})
        utils.SendJSReply(map[string]interface{}{"result": "ok", "data": data}, this.Response)
    }
}

func (this *GridHandler) ConfirmOrRejectPersonRequest() {
    if !sessions.CheackSession(this.Response, this.Request) {
        http.Redirect(this.Response, this.Request, "/", http.StatusUnauthorized)
        return
    }

    if !this.isAdmin() {
        return
    }

    request, err := utils.ParseJS(this.Request, this.Response)
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)

    } else {
        event_id, err := strconv.Atoi(request["event_id"].(string))
        if utils.HandleErr("[GridHandler::GetPersonRequest] event_id Atoi: ", err, this.Response) {
            return
        }
        reg_id, err := strconv.Atoi(request["reg_id"].(string))
        if utils.HandleErr("[GridHandler::GetPersonRequest] reg_id Atoi: ", err, this.Response) {
            return
        }

        query := `SELECT param_values.value, users.id as user_id
            FROM reg_param_vals
            INNER JOIN registrations ON registrations.id = reg_param_vals.reg_id
            INNER JOIN param_values ON param_values.id = reg_param_vals.param_val_id
            INNER JOIN params ON params.id = param_values.param_id
            INNER JOIN events ON events.id = registrations.event_id
            INNER JOIN faces ON faces.id = registrations.face_id
            INNER JOIN users ON users.id = faces.user_id
            WHERE params.id in (1, 4) AND users.id in (
                SELECT users.id FROM registrations INNER JOIN events ON events.id = registrations.event_id
                INNER JOIN faces ON faces.id = registrations.face_id
                INNER JOIN users ON users.id = faces.user_id
                WHERE registrations.id = $1
            ) ORDER BY params.id;`

        data := db.Query(query, []interface{}{reg_id})

        if len(data) < 2 {
            utils.SendJSReply(map[string]interface{}{"result": "Нет данных о логине или e-mail пользователя."}, this.Response)
            return
        }

        to := data[0].(map[string]interface{})["value"].(string)
        email := data[1].(map[string]interface{})["value"].(string)
        event := db.Query("SELECT name FROM events WHERE id=$1;", []interface{}{event_id})[0].(map[string]interface{})["name"].(string)

        if request["confirm"].(bool) {
            if event_id == 1 {
                utils.SendJSReply(map[string]interface{}{"result": "Эту заявку нельзя подтвердить письмом."}, this.Response)

            } else {
                if event_id == 2 {
                    user_id := int(data[0].(map[string]interface{})["user_id"].(int64))
                    user := GetModel("users")
                    user.LoadModelData(map[string]interface{}{"role": "head"})
                    user.GetFields().(*models.User).Enabled = true
                    user.LoadWherePart(map[string]interface{}{"id": user_id})
                    db.QueryUpdate_(user).Scan()
                }

                mailer.SendEmailToConfirmRejectPersonRequest(to, email, event, true)
                utils.SendJSReply(map[string]interface{}{"result": "Письмо с подтверждением заявки отправлено."}, this.Response)
            }

        } else {
            if event_id == 1 {
                utils.SendJSReply(map[string]interface{}{"result": "Эту заявку нельзя отклонить письмом."}, this.Response)

            } else {
                query := `DELETE
                    FROM param_values USING reg_param_vals
                    WHERE param_values.id in (SELECT reg_param_vals.param_val_id WHERE reg_param_vals.reg_id = $1);`
                db.Query(query, []interface{}{reg_id})

                query = `DELETE FROM registrations WHERE id = $1;`
                db.Query(query, []interface{}{reg_id})

                mailer.SendEmailToConfirmRejectPersonRequest(to, email, event, false)
                utils.SendJSReply(map[string]interface{}{"result": "Письмо с отклонением заявки отправлено."}, this.Response)
            }
        }
    }
}

func (this *GridHandler) EditParams() {
    if !sessions.CheackSession(this.Response, this.Request) {
        http.Redirect(this.Response, this.Request, "/", http.StatusUnauthorized)
        return
    }

    // if !this.isAdmin() {
    //     return
    // }

    request, err := utils.ParseJS(this.Request, this.Response)
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
        return
    }

    for _, v := range request["data"].([]interface{}) {

        param_val_id, err := strconv.Atoi(v.(map[string]interface{})["id"].(string))
        if err != nil {
            utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
            return
        }

        value := v.(map[string]interface{})["value"].(string)

        // !!!
        if value == "" {
            value = " "
        }

        param_value := GetModel("param_values")
        param_value.LoadModelData(map[string]interface{}{"value": value})
        param_value.LoadWherePart(map[string]interface{}{"id": param_val_id})
        db.QueryUpdate_(param_value).Scan()
    }

    utils.SendJSReply(map[string]interface{}{"result": "Изменения сохранены."}, this.Response)
}
