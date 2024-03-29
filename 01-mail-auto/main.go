package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strings"
)
func main(){

	http.HandleFunc("/mail",Mail)
	fmt.Println(http.ListenAndServe(":9191",nil))
}
func Mail(w http.ResponseWriter, r *http.Request)  {
	var index int
	t1, err := template.ParseFiles("mail.html")
	if err != nil {
		log.Fatal(err)
	}

	subject:="Subject:"+r.FormValue("subject")+"\n"
	data:=r.FormValue("comment")
	file,f1,_:=r.FormFile("uploadfile")
	err1 := t1.ExecuteTemplate(w, "mail.html", nil)
	if err1 != nil {
		log.Fatal(err1)
	}
	if file!=nil {
		csvfile, _ := f1.Open()
		reader := csv.NewReader(bufio.NewReader(csvfile))
		var fname []string
		line, _ := reader.ReadAll()
		for i,em:=range line[0]{
			if strings.Contains(strings.ToLower(em),"mail"){
				index=i
			}
		}
		for _, name := range line[1:] {
			fname = append(fname, name[index])
		}
		//fmt.Println(fname)
		go mailsend(data,subject,fname)
	}
}
func mailsend(data string,subject string,fname []string){
	auth := smtp.PlainAuth(
		"",
		os.Getenv("USER"),//sender email address
		os.Getenv("PASS"),// password
		"smtp.gmail.com",
	)
	msg := []byte(subject+data)

	err := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		os.Getenv("USER"),//sender email address
		fname,
		msg,
	)

	if err != nil {
		log.Fatal(err)
	}
}