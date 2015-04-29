package models

import (
    "github.com/orc/db"
    "strconv"
)

type GroupsModel struct {
    Entity
}

type Groups struct {
    Id    int    `name:"id" type:"int" null:"NOT NULL" extra:"PRIMARY"`
    Name  string `name:"name" type:"text" null:"NOT NULL" extra:"UNIQUE"`
    Owner int    `name:"face_id" type:"int" null:"NOT NULL" extra:"REFERENCES" refTable:"faces" refField:"id" refFieldShow:"id"`
}

func (c *ModelManager) Groups() *GroupsModel {
    model := new(GroupsModel)

    model.TableName = "groups"
    model.Caption = "Группы"

    model.Columns = []string{"id", "name", "face_id"}
    model.ColNames = []string{"ID", "Название", "Владелец"}

    model.Fields = new(Groups)
    model.WherePart = make(map[string]interface{}, 0)
    model.Condition = AND
    model.OrderBy = "id"
    model.Limit = "ALL"
    model.Offset = 0

    model.Sub = true
    model.SubTable = []string{"persons"}
    model.SubField = "group_id"

    return model
}

func (this *GroupsModel) GetModelRefDate() (fields []string, result map[string]interface{}) {
    fields = []string{"name"}

    result = make(map[string]interface{})

    query := `SELECT faces.id as id, array_to_string(array_agg(param_values.value), ' ') as name
        FROM reg_param_vals
        INNER JOIN registrations ON registrations.id = reg_param_vals.reg_id
        INNER JOIN faces ON faces.id = registrations.face_id
        INNER JOIN events ON events.id = registrations.event_id
        INNER JOIN param_values ON param_values.id = reg_param_vals.param_val_id
        INNER JOIN params ON params.id = param_values.param_id
        WHERE params.id in (5, 6, 7) AND events.id = 1 GROUP BY faces.id ORDER BY faces.id;`

    result["face_id"] = db.Query(query, nil)

    return fields, result
}

func (this *GroupsModel) Select(fields []string, filters map[string]interface{}, limit, offset int, sord, sidx string) (result []interface{}) {
    if len(fields) == 0 {
        return nil
    }

    query := `SELECT `

    for _, field := range fields {
        switch field {
        case "id":
            query += "groups.id, "
            break
        case "name":
            query += "groups.name as group_name, "
            break
        case "face_id":
            query += "array_to_string(array_agg(param_values.value), ' ') as face_name, "
            break
        }
    }

    query = query[:len(query)-2]

    query += ` FROM reg_param_vals
        INNER JOIN registrations ON registrations.id = reg_param_vals.reg_id
        INNER JOIN faces ON faces.id = registrations.face_id
        INNER JOIN events ON events.id = registrations.event_id
        INNER JOIN param_values ON param_values.id = reg_param_vals.param_val_id
        INNER JOIN params ON params.id = param_values.param_id
        INNER JOIN groups ON groups.face_id = faces.id`

    where, params := this.Where(filters)
    if where != "" {
        query += where + ` AND params.id in (5, 6, 7) AND events.id = 1 GROUP BY groups.id`
    } else {
        query += ` WHERE params.id in (5, 6, 7) AND events.id = 1 GROUP BY groups.id`
    }

    if sidx != "" {
        query += ` ORDER BY groups.`+sidx
    }

    query += ` `+ sord

    if limit != -1 {
        params = append(params, limit)
        query += ` LIMIT $`+strconv.Itoa(len(params))
    }

    if offset != -1 {
        params = append(params, offset)
        query += ` OFFSET $`+strconv.Itoa(len(params))
    }

    query += `;`

    return db.Query(query, params)
}
