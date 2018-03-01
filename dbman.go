package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type article struct {
	id        string
	articleId string
	text      string
	user      int
}

func connect(url string) *sql.DB {
	db, err := sql.Open("postgres", url)
	if err != nil {
		log.Fatal("Couldn't connect to postgres", err)
		return nil
	}
	return db
}

func close(db *sql.DB) {
	if err := db.Close(); err != nil {
		log.Println("couldn't close db", err)
	}
}

func getArticles(db *sql.DB, table string) []article {
	arts := make([]article, 0)

	queryStr := fmt.Sprintf("select id, article_id, text, user_id from %s", table)
	if table == "Modified" {
		queryStr = fmt.Sprintf("%s where coded=false", queryStr)
	}
	rows, err := db.Query(queryStr)
	if err != nil {
		log.Fatal(fmt.Sprintf("error querying for %s", table, err))
		return arts
	}
	defer rows.Close()

	for rows.Next() {
		a := article{}
		if err := rows.Scan(&a.id, &a.articleId, &a.text, &a.user); err != nil {
			log.Fatal("error scanning row")
			return arts
		}
		arts = append(arts, a)
	}
	return arts
}

func saveAnswers(db *sql.DB, a article, answers []int) {
	queryStr := "update Modified set coded=true, "
	for i, answer := range answers {
		queryStr = fmt.Sprintf("%sq%d=%d", queryStr, i+1, answer)
		if i < len(answers)-1 {
			queryStr = fmt.Sprintf("%s, ", queryStr)
		}
	}
	queryStr = fmt.Sprintf("%s where id=%s", queryStr, a.id)
	_, err := db.Exec(queryStr)
	if err != nil {
		log.Fatal("error saving answers: ", err)
		return
	}
}

func getAnswers(db *sql.DB) [][]string {
	answers := [][]string{}

	rows, err := db.Query("select * from Modified")
	if err != nil {
		log.Fatal("error getting all answers", err)
		return answers
	}
	defer rows.Close()

	for rows.Next() {
		row := make([]interface{}, 20)
		if err := rows.Scan(&row[0], &row[1], &row[2], &row[3], &row[4], &row[5], &row[6], &row[7], &row[8], &row[9], &row[10], &row[11], &row[12], &row[13], &row[14], &row[15], &row[16], &row[17], &row[18], &row[19]); err != nil {
			log.Fatal("error scanning row", err)
			return answers
		}
		answer := make([]string, 19)
		answer[0] = row[0].(string)
		answer[1] = row[1].(string)
		answer[2] = fmt.Sprintf("%d", int(row[2].(int64)))
		answer[3] = row[3].(string)
		answer[4] = fmt.Sprintf("%d", int(row[4].(int64)))
		answer[5] = row[5].(string)
		for i := 0; i < 12; i++ {
			if row[i+6] == nil {
				answer[i+6] = "N/A"
				continue
			}
			answer[i+6] = fmt.Sprintf("%d", int(row[i+6].(int64)))
		}
		answer[1] = fmt.Sprintf("\\\"%s\\\"", answer[1])
		answers = append(answers, answer)
	}
	return answers
}
