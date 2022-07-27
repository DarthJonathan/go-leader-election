package main

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	. "github.com/onsi/ginkgo/v2"
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
			It("should initialize and run cron", func() {
				mock.ExpectExec(`CREATE TABLE`).WillReturnResult(sqlmock.NewResult(1, 1))
				InitializeClient(sqlx.NewDb(db, "postgres"))
			})
		})

		//Context("with fewer than 300 pages", func() {
		//	It("should be a short story", func() {
		//		Expect(foxInSocks.Category()).To(Equal(books.CategoryShortStory))
		//	})
		//})
	})

	AfterEach(func() {
		db.Close()
	})
})
