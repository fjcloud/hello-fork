package main
import (
	"fmt"
	"net/http"
	"html/template"
	"os"
	"runtime"
	"embed"
)
//go:embed static/*
var content embed.FS
func helloHandler(w http.ResponseWriter, r *http.Request) {
	response := os.Getenv("RESPONSE")
	if len(response) == 0 {
		response = "I love Application Platforms!"
	}
	// Get CPU architecture
	arch := runtime.GOARCH
	// Get cloud region
	region := os.Getenv("CLOUD_REGION")
	if len(region) == 0 {
		region = "Unknown"
	}
	// Get OpenShift type
	openshiftType := os.Getenv("OPENSHIFT_TYPE")
	if len(openshiftType) == 0 {
		openshiftType = "Unknown"
	}
	// Get cloud type
	cloudType := os.Getenv("CLOUD_TYPE")
	if len(cloudType) == 0 {
		cloudType = "Unknown"
	}
	tmpl := `
	<!DOCTYPE html>
	<html>
	<head>
		<style>
			.container {
				display: flex;
				flex-direction: column;
				align-items: center;
				justify-content: center;
				height: 100vh;
				text-align: center;
			}
			h1 {
				font-size: 2em;
			}
			.arch {
				font-size: 1.5em; /* Larger font size for CPU architecture */
				color: #2c3e50; /* A contrasting color */
				margin-top: 20px; /* Space above the architecture text */
				font-weight: bold; /* Make the text bold */
			}
		</style>
	</head>
	<body>
		<div class="container">
			<img src="https://raw.githubusercontent.com/andyrepton/hello/main/static/openshift.jpg" alt="OpenShift" style="max-width: 100%; max-height: 50%;">
			<h1>{{.Response}}</h1>
			<p class="arch">CPU Architecture: {{.Arch}}</p>
			<p class="arch">{{.OpenshiftType}} running on {{.CloudType}} in {{.Region}}</p>
		</div>
	</body>
	</html>
	`
	t, err := template.New("webpage").Parse(tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := struct {
		Response      string
		Arch         string
		Region       string
		OpenshiftType string
		CloudType    string
	}{
		Response:      response,
		Arch:         arch,
		Region:       region,
		OpenshiftType: openshiftType,
		CloudType:    cloudType,
	}
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	fmt.Println("Servicing an impatient beginner's request.")
}
func listenAndServe(port string) {
	fmt.Printf("serving on %s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
func main() {
	http.HandleFunc("/", helloHandler)
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	go listenAndServe(port)
	select {}
}
