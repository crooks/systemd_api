package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/coreos/go-systemd/v22/dbus"
	"github.com/gorilla/mux"
)

var (
	units []string
)

func init() {
	units = initUnits()
}

// initUnits returns a slice of systemd units to interact with
func initUnits() []string {
	return []string{
		"systemd-networkd.service",
		"systemd-timesyncd.service",
	}
}

func returnAllUnits(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()
	conn, err := dbus.NewWithContext(ctx)
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()

	// https://pkg.go.dev/github.com/coreos/go-systemd/v22/dbus#Conn.ListUnitsByNamesContext
	status, err := conn.ListUnitsByNamesContext(ctx, units)
	if err != nil {
		log.Println(err)
	}
	//fmt.Println(reflect.TypeOf(status))
	json.NewEncoder(w).Encode(status)
}

type startStatus struct {
	Name string
	ID   int
	Err  string
}

func startBackup(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()
	conn, err := dbus.NewWithContext(ctx)
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()
	jid, err := conn.StartUnitContext(ctx, "backup.service", "fail", nil)
	if err != nil {
		log.Println(err)
	}
	resp := startStatus{Name: "backup.service", ID: jid, Err: string(err.Error())}
	json.NewEncoder(w).Encode(resp)
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/status", returnAllUnits)
	myRouter.HandleFunc("/backup", startBackup)
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func main() {
	handleRequests()
}
