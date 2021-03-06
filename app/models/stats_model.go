package models

import "labix.org/v2/mgo/bson"

//Stats object
type Stat struct {
	Id         bson.ObjectId `bson:"_id"`
	StatName   string        `bson:"statname"`
	StatMetric string        `bson:"statmetric"`
}

type UserStat struct {
	Id       bson.ObjectId `bson:"_id"`
	StatName string        `bson:"statistic"`
	Value    float64       `bson:"value"`
	Username string        `bson:"user"`
}

//Achievements object
type Ach struct {
	Id       bson.ObjectId `bson:"_id"`
	AchName  string        `bson:"achname"`
	StatName string        `bson:"stat"`
	MinVal   float64       `bson:"minVal"`
}

type UserAch struct {
	Id       bson.ObjectId `bson:"_id"`
	AchName  string        `bson:"achievement"`
	Complete bool          `bson:"complete"`
	Username string        `bson:"user"`
}
