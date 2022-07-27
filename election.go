package main

import (
	"github.com/go-co-op/gocron"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"log"
	"time"
)

var (
	UUID     string
	DB       *sqlx.DB
	isLeader = false
)

var schema = `
CREATE TABLE IF NOT EXISTS golead_hosts (
    uuid text,
    last_heartbeat timestamp,
    startup timestamp,
    inactive_flag int,
    is_leader boolean
)`

func InitializeClient(db *sqlx.DB) {
	//Generate client UUID
	UUID = uuid.New().String()

	DB = db

	//Migrate
	DB.MustExec(schema)

	//Update Heartbeat calls
	s := gocron.NewScheduler(time.UTC)

	s.Every(5).Seconds().Do(func() {
		updateHeartbeat()
	})

	s.Every(5).Seconds().Do(func() {
		electLeader()
	})
}

func updateHeartbeat() {
	query := `UPDATE golead_hosts SET last_heartbeat = now() WHERE UUID = $1`
	tx, err := DB.Begin()
	if err != nil {
		log.Fatalln(err)
	}
	_, err = tx.Exec(query, UUID)
	if err != nil {
		log.Fatalln(err)
	}
}

func electLeader() {
	//Check if last leader was indisposed
	tx := DB.MustBegin()

	//Update inactive flag
	query := `UPDATE golead_hosts SET inactive_flag = inactive_flag + 1 WHERE last_heartbeat < current_timestamp - interval '1 minute'`
	_, err := tx.Exec(query)
	if err != nil {
		log.Fatalln(err)
	}

	//Delete inactive flag if more than 1
	query = `DELETE FROM golead_hosts WHERE inactive_flag > 3`
	_, err = tx.Exec(query)
	if err != nil {
		log.Fatalln(err)
	}

	//Check If any leader is still active, if more than one, clean
	query = `SELECT count(1) WHERE is_leader = true FOR UPDATE`
	leader := 0
	err = tx.Get(&leader, query)
	if err != nil {
		log.Fatalln(err)
	}

	if leader > 1 {
		//Purge leaders
		query = `UPDATE golead_hosts SET is_leader = false`
		_, err = tx.Exec(query)
		if err != nil {
			log.Fatalln(err)
		}
	} else if leader < 1 {
		//Elect leader
		query = `UPDATE golead_hosts SET is_leader = true WHERE UUID = $1`
		_, err = tx.Exec(query, UUID)
		if err != nil {
			log.Fatalln(err)
		}
		isLeader = true
	}

	//Commit
	err = tx.Commit()
	if err != nil {
		log.Fatalln(err)
	}
}

func IsLeader() bool {
	return isLeader
}
