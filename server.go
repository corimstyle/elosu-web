// websockets.go
package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"

	"github.com/kabukky/httpscerts"
)

// Variables used by the Elo Calculator
var (
	player1elo, player1pc        int
	player2elo, player2pc        int
	winner, finallelo, finalwelo int
	name1, name2                 string
)

func main() {
	// Check if the cert files are available.
	httpscerts.Check("certs/server.pem", "certs/key.pem")
	http.Handle("/elosu/", http.StripPrefix("/elosu/", http.FileServer(http.Dir("elosu"))))
	http.Handle("/home/", http.StripPrefix("/home/", http.FileServer(http.Dir("home"))))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// If root directory is called in /elosu/
		if r.URL.Path[1:] == "elosu" {
			http.ServeFile(w, r, r.URL.Path[1:])
		} else if r.URL.Path[1:] == "" {
			http.ServeFile(w, r, "home/index.html")
		} else {
			fmt.Fprintf(w, "Hello, you've requested \"%s\" and that does not exist on this server.", r.URL.Path)
		}
	})
	// Serve /calc with a text response.
	http.HandleFunc("/calc", func(w http.ResponseWriter, r *http.Request) {
		// Parses Form
		err := r.ParseForm()
		if err != nil {
			http.Error(w, fmt.Sprintf("error parsing url %v", err), 500)
		}
		// Extracts information passed from AJAX statement on examplecalc.html
		p1elo := r.FormValue("P1")
		p1eloint, _ := strconv.Atoi(p1elo)
		p2elo := r.FormValue("P2")
		p2eloint, _ := strconv.Atoi(p2elo)
		output := calcK(1, p1eloint, p2eloint, 0, 0, "Player 1", "Player 2")
		// Display all calc through the console
		println(p1eloint, p2eloint)
		println(output[24:28], output[53:57])
		fmt.Fprintln(w, output)
	})
	// Clears the output
	http.HandleFunc("/clear", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "")
	})
	// Serves the webpage
	errhttps := http.ListenAndServeTLS(":443", "certs/server.pem", "certs/key.pem", nil)
	if errhttps != nil {
		log.Fatal("Web server (HTTPS): ", errhttps)
	}

}

// ********** START ELO CALCULATOR **********
// Does calculatons to find the correct elo of the 2 players and returns an HTML string with the new values
func calcElo(welo, lelo, wk, lk int, wname, lname string) string {
	// Fuck what does this do again
	var wdiv = float64(float64(welo) / 400.0)
	var rw = float64(math.Pow(10, wdiv))
	var ldiv = float64(float64(lelo) / 400.0)
	var rl = float64(math.Pow(10, ldiv))
	finalwelo = welo + int(float64(wk)*(1-(rw/(rw+rl))))
	finallelo = lelo + int(float64(lk)*(0-(rl/(rw+rl))))
	//return formatted elo values
	return fmt.Sprintf("<h3>%s final elo: %d <br>%s final elo: %d \n</h3>", wname, finalwelo, lname, finallelo)
}

// Calculates the K value and calls calcElo using the correct winner/loser ordering
func calcK(winner, player1elo, player2elo, player1pc, player2pc int, name1, name2 string) string {
	var wk, lk int
	var newK = 75
	var oldK = 25
	// Logic for faster new player growth
	if player1pc <= 10 {
		if winner == 1 {
			wk = newK
		} else {
			lk = newK
		}
	} else if player1pc > 10 {
		if winner == 1 {
			wk = oldK
		} else {
			lk = oldK
		}
	}
	if player2pc <= 10 {
		if winner == 2 {
			wk = newK
		} else {
			lk = newK
		}
	} else if player2pc > 10 {
		if winner == 2 {
			wk = oldK
		} else {
			lk = oldK
		}
	}
	// Choses correct order for player that won.
	if winner == 1 {
		return (calcElo(player1elo, player2elo, wk, lk, name1, name2))
	}
	return (calcElo(player2elo, player1elo, wk, lk, name2, name1))

}

// ********** STOP ELO CALCULATOR **********
