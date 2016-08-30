package controllers

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"time"
	"github.com/klenin/orc/utils"
	"github.com/klenin/orc/config"
	"strconv"
	"net/http"
)

type Event struct {
	Id          int        `json:"id" gorm:"primary_key"`
	Name        string     `json:"name"`
	DateStart   time.Time  `json:"dateStart"`
	DateFinish  time.Time  `json:"dateFinish"`
	Time        time.Time  `json:"time"`
	Team        bool       `json:"isTeam"`
	Url         string     `json:"url"`
	Forms       []Form     `json:"forms,omitempty" gorm:"many2many:events_forms"`
}

type Form struct {
	Id        int      `json:"id" gorm:"primary_key"`
	Name      string   `json:"name"`
	Personal  bool     `json:"isPersonal"`
	Params    []Param  `json:"params"`
}

type Param struct {
	Id             int        `json:"id" gorm:"primary_key"`
	Name           string     `json:"name"`
	FormId         int        `json:"-"`
	ParamType      ParamType  `json:"-"`
	ParamTypeId    int        `json:"-"`
	ParamTypeName  string     `json:"type" gorm:"-"`
	Identifier     int        `json:"identifier"`
	Required       bool       `json:"isRequired"`
	Editable       bool       `json:"isEditable"`
}

type ParamType struct {
	Id    int     `gorm:"primary_key"`
	Name  string
}

var gormDb *gorm.DB

func init() {
	var err error
	gormDb, err = gorm.Open("postgres", config.GetValue("DATABASE_URL"))
	if err != nil {
		log.Fatal("failed to connect database")
	}
}

func (c *BaseController) Api() *Api {
	return new(Api)
}

type Api struct {
	Controller
}

func (a *Api) Event(q string) {
	a.Response.Header().Add("Access-Control-Allow-Origin", "http://localhost:3000")

	if q == "" {
		var events []Event
		query := gormDb.Order("date_start desc")
		if name := a.Request.URL.Query().Get("name"); name != "" {
			query = query.Where("lower(name) LIKE lower(?)", "%" + name + "%")
		}
		query.Find(&events)
		utils.SendJSReply(events, a.Response)
		return
	}

	id, err := strconv.Atoi(q)
	if err != nil {
		http.Error(a.Response, "Invalid id", http.StatusBadRequest)
		return
	}

	var event Event
	gormDb.First(&event, id)
	if event.Id == 0 {
		http.NotFound(a.Response, a.Request)
		return
	}

	gormDb.Model(&event).Related(&event.Forms, "Forms")
	for i, form := range event.Forms {
		gormDb.Model(&form).Related(&event.Forms[i].Params)
		for j, param := range event.Forms[i].Params {
			gormDb.Model(&param).Related(&event.Forms[i].Params[j].ParamType)
			event.Forms[i].Params[j].ParamTypeName = event.Forms[i].Params[j].ParamType.Name
		}
	}

	utils.SendJSReply(event, a.Response)
}
