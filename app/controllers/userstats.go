package controllers

import (
	"fmt"
	"github.com/revel/revel"
	"labix.org/v2/mgo/bson"
	"leaderboard/app/models"
)

func (c App) GetStat(name string) *models.Stat {

	// connect to DB server(s)
	d, s := db(statcol)

	// Query
	results := models.Stat{}
	query := bson.M{"statname": name}
	err := d.Find(query).One(&results)

	if err != nil {
		panic(err)
	}
	if len(results.StatName) == 0 {
		return nil
	}

	s.Close()

	return &results

}

func (c App) GetAchStat(name string) *models.Ach {

	// connect to DB server(s)
	d, s := db(achcol)

	// Query
	results := models.Ach{}
	query := bson.M{"stat": name}
	err := d.Find(query).One(&results)

	if err != nil {
		panic(err)
	}
	if len(results.AchName) == 0 {
		return nil
	}

	s.Close()

	return &results

}

func (c App) SaveUserStat(statName string, statValue float64) revel.Result {

	if c.Session["user"] == "" || c.Session["role"] != playerrole {
		return c.RenderJson("User is not logged in, or user is not a player")
	} else {
		username := c.Session["user"]
		stat := c.GetStat(statName)
		user := c.GetUser(username)
		ach := c.GetAchStat(stat.StatName)
		// connect to DB server
		d, s := db(userstatcol)

		// Query
		var doc models.UserStat
		results := models.UserStat{}
		query := bson.M{"statistic": stat.StatName, "user": user.Username}
		err := d.Find(query).One(&results)

		if len(results.StatName) == 0 {
			//do DB operations
			doc = models.UserStat{Id: bson.NewObjectId(), StatName: statName, Value: statValue, Username: user.Username}
			err = d.Insert(doc)
			if err != nil {
				panic(err)
			} else {
				if results.Value > ach.MinVal {
					// works OK but what if we insert the new stat increment and
					// inserting of achievement fails? then user would have
					// missed an achievement, so we need a way of controlling
					// when achievements are missed. Perhaps we need to check
					// if achievement has been reached before incrementing stat
					// in database... Same for line 104
					c.Achieve(ach.AchName, true)
				}
				s.Close()
				return c.RenderJson(doc)
			}
		} else {

			newValue := statValue + results.Value
			fmt.Println(newValue) // debug
			colQuerier := bson.M{"statistic": stat.StatName, "user": user.Username}
			change := bson.M{"$set": bson.M{"value": newValue}}
			err = d.Update(colQuerier, change)
			if err != nil {
				panic(err)
			} else {
				if newValue > ach.MinVal {
					c.Achieve(ach.AchName, true)
				}
				s.Close()
				return c.RenderJson(err)
			}
		}

		if err != nil {
			panic(err)
		}

		return c.RenderJson(results)
	}

}

func (c App) GetUserStats(username string) revel.Result {

	// connect to DB server(s)
	d, s := db(userstatcol)

	var user string

	if len(username) == 0 {
		user = c.Session["user"]
	} else {
		user = username
	}

	// Query
	var results []models.UserStat
	query := bson.M{"user": user}
	err := d.Find(query).Sort("-timestamp").All(&results)

	if err != nil {
		panic(err)
	}
	if len(results) == 0 {
		return nil
	}

	s.Close()

	return c.RenderJson(results)

}
