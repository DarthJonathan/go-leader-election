package main

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	. "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"log"
)

var _ = Describe("Election", func() {

	var db *sql.DB
	var mock sqlmock.Sqlmock

	BeforeEach(func() {
		var err error
		db, mock, err = sqlmock.New()
		if err != nil {
			log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
	})

	Describe("Leader Election", func() {
		Context("with normal database", func() {
			It("should initialize and run migration", func() {
				//Expect
				mock.ExpectExec(`CREATE TABLE`).WillReturnResult(sqlmock.NewResult(1, 1))

				//Run
				InitializeClient(sqlx.NewDb(db, "postgres"))
			})

			It("should handle migration error", func() {
				//Expect
				mock.ExpectExec(`CREATE TABLE`).WillReturnError(errors.New("Error"))

				//Run
				gomega.Expect(func() {
					InitializeClient(sqlx.NewDb(db, "postgres"))
				}).To(gomega.Panic())
			})
		})

		Context("with normal database", func() {
			BeforeEach(func() {
				DB = sqlx.NewDb(db, "postgres")
			})

			It("update heartbeat successfully", func() {
				//Expect
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE golead_hosts SET last_heartbeat`).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()

				//Run
				updateHeartbeat()

				// we make sure that all expectations were met
				if err := mock.ExpectationsWereMet(); err != nil {
					log.Fatalf("there were unfulfilled expectations: %s", err)
				}
			})
		})
	})

	AfterEach(func() {
		db.Close()
	})
})
