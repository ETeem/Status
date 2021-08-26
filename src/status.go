package main

import (
	"os"
	"fmt"
	"strings"
	"io/ioutil"

	"net/http"
	"github.com/gorilla/mux"
)

func handleWhoAreYou(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "SFTP Status")
}

func handlePing(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "pong")
}

func handleDescription(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Blue/Green SFTP Status")
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	stat, err := ioutil.ReadFile("/etc/status/status.dat")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("501 - " + err.Error()))
		return
	}

	strstat := strings.TrimSpace(string(stat))
	if strstat == "500" {
		http.Error(w, http.StatusText(500), 500)
		return 
	}

	fmt.Fprintf(w, "live")
	return 
}

func handleSetStatus(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	if len(status) == 0 {
		fmt.Fprintf(w, "Missing 'status' parameter")
		return
	}

	if status != "standby" && status != "live" {
		fmt.Fprintf(w, "Invalid 'status' parameter.  Can Be: standby or live")
		return
	}

	if status == "standby" {
		d1 := []byte("500\n")
		_ = ioutil.WriteFile("/etc/status/status.dat", d1, 0644)
		fmt.Fprintf(w, "set to standby")
		return
	} 

	d1 := []byte("200\n")
	_ = ioutil.WriteFile("/etc/status/status.dat", d1, 0644)
	fmt.Fprintf(w, "set to live")
}

func handleFlip(w http.ResponseWriter, r *http.Request) {
	stat, err := ioutil.ReadFile("/etc/status/status.dat")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("501 - " + err.Error()))
		return
	}

	strstat := strings.TrimSpace(string(stat))
	if strstat == "500" {
		d1 := []byte("200\n")
		_ = ioutil.WriteFile("/etc/status/status.dat", d1, 0644)
		fmt.Fprintf(w, "set to live")
		return 
	} else {
		d1 := []byte("500\n")
		_ = ioutil.WriteFile("/etc/status/status.dat", d1, 0644)
		fmt.Fprintf(w, "set to standby")
		return 
	}
}

func handleMain(w http.ResponseWriter, r *http.Request) {
	stat, err := ioutil.ReadFile("/etc/status/status.dat")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("501 - " + err.Error()))
		return
	}

	strstat := strings.TrimSpace(string(stat))
	if strstat == "500" {
		fmt.Fprintf(w, "System State Is Currently: Standby")
		return
	}

	fmt.Fprintf(w, "System State Is Currently: Live")
}

func main() {

	if _, xerr := os.Stat("/etc/status/"); xerr != nil {
		os.Mkdir("/etc/status/", 0755)
	}

	if _, yerr := os.Stat("/etc/status/status.dat"); yerr != nil {
		d1 := []byte("500\n")
		_ = ioutil.WriteFile("/etc/status/status.dat", d1, 0644)
	}

        router := mux.NewRouter()
        router.HandleFunc("/whoareyou", handleWhoAreYou)
        router.HandleFunc("/ping", handlePing)
        router.HandleFunc("/description", handleDescription)
        router.HandleFunc("/status", handleStatus)
        router.HandleFunc("/setstatus", handleSetStatus)
        router.HandleFunc("/flip", handleFlip)
        router.HandleFunc("/", handleMain)

        err := http.ListenAndServe(":80", router)
        if err != nil {
                fmt.Println("ListenAndServe: ", err)
        }

}
