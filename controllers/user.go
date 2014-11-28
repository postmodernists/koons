package controllers

import (
  "net/http"
  "encoding/json"
  "fmt"
  "../models"
  "log"
  // "time"
  "io/ioutil"
  "github.com/coopernurse/gorp"
)

func usersHandler(cb1 func(r *http.Request) *models.User, cb2 func(r *http.Request) *[]models.User) func(w http.ResponseWriter, r*http.Request) {

  return func(w http.ResponseWriter, r *http.Request) {
    var js []byte
    var dbError error

    if(cb1 != nil){
      js, dbError = json.Marshal( cb1(r) )
    } else{
      js, dbError = json.Marshal( cb2(r) )
    }

    if dbError != nil {
      http.Error(w, dbError.Error(), http.StatusInternalServerError)
      log.Fatal(dbError)
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(js)
  }
}

func GenericHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Whats up breh?")
}

func GetAllUsers(dbmap *gorp.DbMap) func(w http.ResponseWriter, r *http.Request) {

  return usersHandler(nil, func(r *http.Request) *[]models.User {
    var users []models.User
    _, dbError := dbmap.Select(&users, "select * from \"user\"")
    if dbError != nil {
      log.Fatal(dbError)
    }

    return &users
  })

}

func GetUser(dbmap *gorp.DbMap) func(w http.ResponseWriter, r *http.Request) {

  return usersHandler(func(r *http.Request) *models.User  {
    user, dbError := dbmap.Get(models.User{}, 1)
    if dbError != nil {
      log.Fatal(dbError)
    }

    return user.(*models.User)
  }, nil)
}

func CreateUser(dbmap *gorp.DbMap) func(w http.ResponseWriter, r *http.Request) {

  return usersHandler(func(r *http.Request) *models.User {
    body, readError := ioutil.ReadAll(r.Body)
    if readError != nil {
        log.Fatal(readError)
    }

    var user *models.User
    jsError := json.Unmarshal(body, &user)
    if jsError != nil {
        log.Fatal(jsError)
    }

    dbError := dbmap.Insert(user)
    if dbError != nil {
      log.Fatal(dbError)
    }

    return user
  }, nil)
}
