package core

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const dateFormat = "02.01.2006"

type Period struct {
	Id          primitive.ObjectID `bson:"_id,omitempty"`
	StartDate   primitive.DateTime `bson:"start,omitempty"`
	EndDate     primitive.DateTime `bson:"end,omitempty"`
	IsCompleted bool               `bson:"is_completed"`
}

type PeriodIncomeExpense struct {
	Period       `bson:",inline"`
	TotalIncome  int64 `bson:"total_income,omitempty"`
	TotalExpense int64 `bson:"total_expense,omitempty"`
}

type PeriodDto struct {
	Id           string            `json:"id,omitempty"`
	StartDate    Date              `json:"start,omitempty"`
	EndDate      Date              `json:"end,omitempty"`
	TotalIncome  int64             `json:"total_income,omitempty"`
	TotalExpense int64             `json:"total_expense,omitempty"`
	IsCompleted  bool              `json:"is_completed"`
	Details      []PeriodDetailDto `json:"details,omitempty"`
}

type PeriodDetail struct {
	Id       primitive.ObjectID `bson:"_id,omitempty"`
	Title    string             `bson:"title,omitempty"`
	Amount   int                `bson:"amount,omitempty"`
	Status   Status             `bson:"status,omitempty"`
	Type     Type               `bson:"type,omitempty"`
	PeriodId primitive.ObjectID `bson:"period_id,omitempty"`
}

type PeriodDetailDto struct {
	Id       string `json:"id,omitempty"`
	Title    string `json:"title,omitempty"`
	Amount   int    `json:"amount,omitempty"`
	Status   Status `json:"status,omitempty"`
	Type     Type   `json:"type,omitempty"`
	PeriodId string `json:"period_id,omitempty"`
}

func NewPeriod(startDateStr string, endDateStr string) *Period {
	startDate, _ := time.Parse(time.DateOnly, startDateStr)
	endDate, _ := time.Parse(time.DateOnly, endDateStr)
	return &Period{
		StartDate:   primitive.NewDateTimeFromTime(startDate),
		EndDate:     primitive.NewDateTimeFromTime(endDate),
		IsCompleted: false,
	}
}

func (pDto PeriodDto) FormatStartDate() string {
	return time.Time(pDto.StartDate).Format(dateFormat)
}

func (pDto PeriodDto) FormatEndDate() string {
	return time.Time(pDto.EndDate).Format(dateFormat)
}

func NewPeriodDetail(title string, amount int, dType Type, periodId primitive.ObjectID) *PeriodDetail {
	return &PeriodDetail{
		Title:    title,
		Amount:   amount,
		Type:     dType,
		Status:   NEW,
		PeriodId: periodId,
	}
}

func (pDto PeriodDto) EvalColor() string {
	expIncRel := float64(pDto.TotalExpense) / float64(pDto.TotalIncome)
	if expIncRel < 0.7 {
		return "bg-success"
	} else if expIncRel >= 0.7 && expIncRel < 1 {
		return "bg-warning"
	} else {
		return "bg-danger"
	}
}

func (p Period) toDto() PeriodDto {
	start := Date(p.StartDate.Time())
	end := Date(p.EndDate.Time())
	return PeriodDto{
		Id:          p.Id.Hex(),
		StartDate:   start,
		EndDate:     end,
		IsCompleted: p.IsCompleted,
	}
}

func (p PeriodDetail) toDto() PeriodDetailDto {
	return PeriodDetailDto{
		Id:       p.Id.Hex(),
		Title:    p.Title,
		Amount:   p.Amount,
		Status:   p.Status,
		Type:     p.Type,
		PeriodId: p.PeriodId.Hex(),
	}
}

func (p PeriodIncomeExpense) toDto() PeriodDto {
	dto := p.Period.toDto()
	dto.TotalIncome = p.TotalIncome
	dto.TotalExpense = p.TotalExpense
	return dto
}

func (p *PeriodDetailDto) EvalColor() string {
	if p.Type == INCOME {
		return "bg-success"
	} else {
		return "bg-danger"
	}
}

type Status string
type Type string

const (
	NEW     Status = "NEW"
	SAVING  Status = "SAVING"
	REALIZE Status = "REALIZE"
)

const (
	INCOME   Type = "INCOME"
	SPENDING Type = "EXPENSE"
)

type Date time.Time
