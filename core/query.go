package core

import "go.mongodb.org/mongo-driver/bson"

var periodsIncomeExpenseQuery = bson.A{
	bson.M{
		"$lookup": bson.M{
			"from":         "period_details",
			"localField":   "_id",
			"foreignField": "period_id",
			"as":           "details",
		},
	},
	bson.M{
		"$addFields": bson.M{
			"total_income": bson.M{
				"$sum": bson.M{
					"$map": bson.M{
						"input": bson.M{
							"$filter": bson.M{
								"input": "$details",
								"as":    "detail",
								"cond": bson.M{
									"$eq": bson.A{"$$detail.type", "INCOME"},
								},
							},
						},
						"as": "income_detail",
						"in": "$$income_detail.amount",
					},
				},
			},
			"total_expense": bson.M{
				"$sum": bson.M{
					"$map": bson.M{
						"input": bson.M{
							"$filter": bson.M{
								"input": "$details",
								"as":    "detail",
								"cond": bson.M{
									"$eq": bson.A{"$$detail.type", "EXPENSE"},
								},
							},
						},
						"as": "expense_detail",
						"in": "$$expense_detail.amount",
					},
				},
			},
		},
	},
	bson.M{
		"$project": bson.M{
			"details": 0,
		},
	},
}
