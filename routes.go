package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func sendTemplate(f string, d map[string]interface{}, w http.ResponseWriter) {
	t := templateHandler{filename: f}
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("public", t.filename)))
	})
	d["host"] = host
	t.templ.Execute(w, d)
}

func login(w http.ResponseWriter, r *http.Request) {
	sendTemplate("index.html", make(map[string]interface{}), w)
}

func sendLogin(w http.ResponseWriter, r *http.Request) {
	password := r.URL.Query().Get("password")
	if password == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("No password entered"))
		return
	}

	for user, pass := range passwords {
		if pass == password {
			modTmp := make([]article, 0)
			for _, a := range modified {
				if a.user == user {
					modTmp = append(modTmp, a)
				}
			}

			if user == 1 {
				user1 = modTmp
			} else if user == 2 {
				user2 = modTmp
			} else if user != 0 {
				log.Fatal("bad user id", user)
				return
			}

			w.WriteHeader(http.StatusOK)
			return
		}
	}
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("This password was not correct"))
}

func response(w http.ResponseWriter, r *http.Request) {
	numberStr := r.URL.Query().Get("n")
	number, err := strconv.ParseInt(numberStr, 10, 64)
	if err != nil {
		log.Fatal("unable to convert n param to int")
	}

	userStr := r.URL.Query().Get("user")
	user, err := strconv.ParseInt(userStr, 10, 64)
	if err != nil {
		log.Fatal("unable to convert user param to int")
	}

	var currMod []article
	if int(user) == 1 {
		currMod = user1
	} else if int(user) == 2 {
		currMod = user2
	} else {
		log.Fatal("invalid user id", user)
		return
	}

	mod := currMod[int(number)-1]
	articleID := mod.articleId

	data := make(map[string]interface{})
	for _, o := range originals {
		if o.articleId == articleID {
			data["original"] = o.text
			break
		}
	}
	data["modified"] = mod.text
	data["current"] = int(number)
	data["total"] = len(currMod)
	data["user"] = int(user)
	data["questions"] = convToJSON(questions)
	sendTemplate("response.html", data, w)
}

func sendResponse(w http.ResponseWriter, r *http.Request) {
	numberStr := r.URL.Query().Get("n")
	number, err := strconv.ParseInt(numberStr, 10, 64)
	if err != nil {
		log.Fatal("unable to convert n param to int")
	}

	userStr := r.URL.Query().Get("user")
	user, err := strconv.ParseInt(userStr, 10, 64)
	if err != nil {
		log.Fatal("unable to convert user param to int")
	}

	answers := r.URL.Query().Get("answers")
	answerStrArr := strings.Split(answers, "-")
	answersArr := make([]int, 12)

	for i, answerStr := range answerStrArr {
		answer, err := strconv.ParseInt(answerStr, 10, 64)
		if err != nil {
			log.Println("could not parse answer int")
			continue
		}
		answersArr[i] = int(answer)
	}

	if user == 1 {
		saveAnswers(db, user1[int(number)-1], answersArr)
	} else if user == 2 {
		saveAnswers(db, user2[int(number)-1], answersArr)
	}

	w.WriteHeader(http.StatusOK)
}

func finished(w http.ResponseWriter, r *http.Request) {
	sendTemplate("finished.html", make(map[string]interface{}), w)
}

func seeAll(w http.ResponseWriter, r *http.Request) {
	var buffer bytes.Buffer

	answers := getAnswers(db)
	buffer.WriteString("article_id, text, id, l/m, wave, subj_id, user_id, q1, q2, q3, q4, q5, q6, q7, q8, q9, q10, q11, q12\n")
	for _, answer := range answers {
		buffer.WriteString(strings.Join(answer, ", "))
		buffer.WriteString("\n")
	}

	w.Write(buffer.Bytes())
}

func convToJSON(qs []question) string {
	qj := "{\"questions\": ["
	for i, q := range qs {
		qj = fmt.Sprintf("%s{\"title\": \"%s\", \"choices\": [", qj, q.title)
		for j, o := range q.options {
			qj = fmt.Sprintf("%s\"%s\"", qj, o)
			if j < len(q.options)-1 {
				qj = fmt.Sprintf("%s, ", qj)
			}
		}
		qj = fmt.Sprintf("%s]}", qj)
		if i < len(qs)-1 {
			qj = fmt.Sprintf("%s,", qj)
		}
	}
	return fmt.Sprintf("%s]}", qj)
}
