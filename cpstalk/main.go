package main

import "fmt"
// import "./scripts"
import "net/http"

func main() {

	// fmt.Println(scripts.GetContests())

	// Subs := scripts.GetSubmissions("tryingtocode")
	// for i:=0; i<len(Subs.Data); i++ {
	// 	fmt.Println(Subs.Data[i])
	// }

	// Graph := scripts.GetGraphData("tryingtocode")
	// for i:=0; i<len(Graph); i++ {
	// 	fmt.Println(Graph[i])
	// }

	// fmt.Println(scripts.GetProfileInfo("tryingtocode"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Work in progress...")
	})

	http.ListenAndServe(":80", nil)
}