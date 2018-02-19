package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/rs/cors"
)

type question struct {
	title   string
	options []string
}

var db *sql.DB
var envvars map[string]string
var passwords map[int]string
var host string
var questions []question
var originals []article
var modified []article
var user1 []article
var user2 []article

func seedQuestions() {
	questions = make([]question, 12)
	questions[0] = question{
		"Please indicate which of the following best fits the description of the scientific technology",
		[]string{"No description OR mention of technology", "The description was wrong/could not make sense of it", "Only provided the name of the technology", "They provided a correct description of the technology", "Mentioned the name of the technology and provided a correct description of it"},
	}
	questions[1] = question{
		"How many individuals were mentioned in the replication?",
		[]string{"0", "1", "2"},
	}
	questions[2] = question{
		"Individual 1: Was the individual mentioned an expert or a non-expert?",
		[]string{"Expert", "Non-expert", "Expertise not specified", "Did not mention an individual"},
	}
	questions[3] = question{
		"Individual 1: Did the individual provide a general indication of support or opposition for the scientified topic?",
		[]string{"Yes", "No", "Did not mention an individual"},
	}
	questions[4] = question{
		"Individual 1: Did they provide the CORRECT support or opposition for the first individual?",
		[]string{"Yes", "No", "Did not mention an individual"},
	}
	questions[5] = question{
		"Individual 1: Did they provide an argument for their position?",
		[]string{"Yes", "No", "Did not mention an individual"},
	}
	questions[6] = question{
		"Individual 1: Did they provide the correct argument for the first individual?",
		[]string{"Yes", "No", "Did not mention an individual"},
	}
	questions[7] = question{
		"Individual 2: Was the individual mentioned an expert or a non-expert?",
		[]string{"Expert", "Non-expert", "Expertise not specified", "Did not mention an individual"},
	}
	questions[8] = question{
		"Individual 2: Did the individual provide a general indication of support or opposition for the scientific topic?",
		[]string{"Yes", "No", "Did not mention a second individual"},
	}
	questions[9] = question{
		"Individual 2: Did they provide the CORRECT support or opposition for the first individual?",
		[]string{"Yes", "No", "Did not mention a second individual"},
	}
	questions[10] = question{
		"Individual 2: Did they provide an argument for their position?",
		[]string{"Yes", "No", "Did not mention a second individual"},
	}
	questions[11] = question{
		"Individual 2: Did they provide the correct argument for the first individual?",
		[]string{"Yes", "No", "Did not mention a second individual"},
	}
}

func main() {
	seedQuestions()

	passwords = make(map[int]string)
	passwords[0] = "admin0"
	passwords[1] = "first-user1"
	passwords[2] = "second-user2"

	envvars = getEnv()
	dbURL := envvars["CP_DB_URL"]
	//dbURL := "user=samlerner dbname=cp sslmode=disable"
	port := envvars["PORT"]
	addr := ":" + port
	host = envvars["HOST"]

	db = connect(dbURL)
	if db == nil {
		log.Fatal("db is nil")
	}
	defer close(db)

	//originals = getArticles(db, "Originals")
	//modified = getArticles(db, "Modified")

	mux := http.NewServeMux()
	mux.HandleFunc("/", login)
	mux.HandleFunc("/sendLogin", sendLogin)
	mux.HandleFunc("/response", response)
	mux.HandleFunc("/sendResponse", sendResponse)
	mux.HandleFunc("/finished", finished)
	mux.HandleFunc("/seeAll", seeAll)

	handler := cors.Default().Handler(mux)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
