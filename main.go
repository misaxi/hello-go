package main

import (
  "encoding/json"
  "net/http"
  "strings"
)

func main() {
  http.HandleFunc("/", hello)
  http.HandleFunc("/followers/", queryWeatherByCity)
  http.ListenAndServe(":8080", nil)
}

func hello(w http.ResponseWriter, r *http.Request) {
  w.Write([]byte("hello!"))
}

type followerData struct {
  Name string `json:"login"`
  Url string `json:"url"`
}

func queryWeatherByCity(w http.ResponseWriter, r *http.Request) {
  user := strings.SplitN(r.URL.Path, "/", 3)[2]

  data, err := query(user)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json; charset=utf-8")
  json.NewEncoder(w).Encode(data)
}

func query(user string) ([]followerData, error) {
  resp, err := http.Get("https://api.github.com/users/" + user + "/followers")
  if err != nil {
    return []followerData{}, err
  }

  defer resp.Body.Close()

  var d []followerData

  if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
    return []followerData{}, err
  }

  return d, nil
}
