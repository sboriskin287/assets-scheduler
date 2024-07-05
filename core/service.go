package core

import (
	"context"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

const (
	index = "site/index.html"
)

type PeriodService struct {
	periodCol        *mongo.Collection
	periodDetailsCol *mongo.Collection
	totalBalanceCol  *mongo.Collection
}

func NewPeriodService(client *mongo.Client) *PeriodService {
	const dbName = "assets-scheduler"
	var periodsCol = client.Database(dbName).Collection("periods")
	var periodDetailsCol = client.Database(dbName).Collection("period_details")
	var totalBalanceCol = client.Database(dbName).Collection("total_balance")
	return &PeriodService{periodsCol, periodDetailsCol, totalBalanceCol}
}

func (ps *PeriodService) Index(w http.ResponseWriter, _ *http.Request) {
	periods, err := ps.GetPeriods()
	if err != nil {
		return
	}
	tmpl := template.Must(template.ParseFiles(index))
	err = tmpl.Execute(w, periods)
	if err != nil {
		log.Panic(err)
	}
}

func (ps *PeriodService) CreatePeriod(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		return
	}
	startDate := r.Form.Get("start_date")
	endDate := r.Form.Get("end_date")
	period := NewPeriod(startDate, endDate)
	objId, err := ps.createPeriod(period)
	if err != nil {
		return
	}
	insertedPeriod, err := ps.getPeriodById(objId)
	if err != nil {
		return
	}
	tmpl := template.Must(template.ParseFiles(index))
	_ = tmpl.ExecuteTemplate(w, "period", insertedPeriod)
}

func (ps *PeriodService) GetPeriodDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		return
	}
	cur, err := ps.periodDetailsCol.Find(context.TODO(), bson.M{"period_id": id})
	if err != nil {
		return
	}
	periodDetails := make([]PeriodDetailDto, 0)
	for cur.Next(context.TODO()) {
		var periodDetail PeriodDetail
		err := cur.Decode(&periodDetail)
		if err != nil {
			return
		}
		periodDetails = append(periodDetails, periodDetail.toDto())
	}
	tmpl := template.Must(template.ParseFiles(index))
	_ = tmpl.ExecuteTemplate(w, "period_details", map[string][]PeriodDetailDto{"Details": periodDetails})
}

func (ps *PeriodService) CreatePeriodDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	periodId, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		return
	}
	err = r.ParseForm()
	if err != nil {
		return
	}
	name := r.Form.Get("title")
	amount, err := strconv.Atoi(r.Form.Get("amount"))
	if err != nil {
		return
	}
	dType := Type(r.Form.Get("type"))
	periodDetail := NewPeriodDetail(name, amount, dType, periodId)
	objId, err := ps.createPeriodDetail(periodDetail)
	if err != nil {
		return
	}
	insertedPeriodDetail, err := ps.getPeriodDetailById(objId)
	if err != nil {
		return
	}
	tmpl := template.Must(template.ParseFiles(index))
	_ = tmpl.ExecuteTemplate(w, "period_detail", insertedPeriodDetail)
}

func (ps *PeriodService) GetPeriods() ([]PeriodDto, error) {
	var periods []PeriodDto
	cur, err := ps.periodCol.Aggregate(context.TODO(), periodsIncomeExpenseQuery)
	if err != nil {
		return nil, err
	}
	for cur.Next(context.TODO()) {
		var period PeriodIncomeExpense
		err := cur.Decode(&period)
		if err != nil {
			return nil, err
		}
		periods = append(periods, period.toDto())
	}
	return periods, nil
}

func (ps *PeriodService) getPeriodById(id primitive.ObjectID) (*PeriodIncomeExpense, error) {
	var period PeriodIncomeExpense
	res := ps.periodCol.FindOne(context.TODO(), bson.M{"_id": id})
	err := res.Decode(&period)
	if err != nil {
		return nil, err
	}
	return &period, nil
}

func (ps *PeriodService) getPeriodDetailById(id primitive.ObjectID) (*PeriodDetail, error) {
	var periodDetail PeriodDetail
	res := ps.periodDetailsCol.FindOne(context.TODO(), bson.M{"_id": id})
	err := res.Decode(&periodDetail)
	if err != nil {
		return nil, err
	}
	return &periodDetail, nil
}

func (ps *PeriodService) createPeriodDetail(detail *PeriodDetail) (primitive.ObjectID, error) {
	res, err := ps.periodDetailsCol.InsertOne(context.TODO(), detail)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return res.InsertedID.(primitive.ObjectID), nil
}

func (ps *PeriodService) createPeriod(period *Period) (primitive.ObjectID, error) {
	res, err := ps.periodCol.InsertOne(context.TODO(), period)
	if err != nil {
		return primitive.NilObjectID, err
	}
	objId := res.InsertedID.(primitive.ObjectID)
	return objId, nil
}

func (ps *PeriodService) handleError(w http.ResponseWriter, err error) {
	_ = template.
		Must(template.ParseFiles("site/error.html")).
		Execute(w, err)
}
