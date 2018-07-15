package main

import (
    "crypto/tls"
    "log"
    "net/http"
    "io/ioutil"
    "fmt"
    "net/url"
)

func redirectTLS(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://"+r.Host+":443"+r.RequestURI, http.StatusMovedPermanently)
}


func index(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
//        w.Write([]byte("This is an example server.\n"))

	if req.URL.Path != "/" {
		log.Printf("404: %s", req.URL.String())
		http.NotFound(w, req)
		return
	}
	http.ServeFile(w, req, "index.html")
}


func main() {
	go http.ListenAndServe(":80", http.HandlerFunc(redirectTLS))
    mux := http.NewServeMux()
    mux.HandleFunc("/", index) //func(w http.ResponseWriter, req *http.Request) {
 //       w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
  //      w.Write([]byte("This is an example server.\n"))
    
    cfg := &tls.Config{
        MinVersion:               tls.VersionTLS12,
        CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
        PreferServerCipherSuites: true,
        CipherSuites: []uint16{
            tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
            tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
            tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
            tls.TLS_RSA_WITH_AES_256_CBC_SHA,
        },
    }
    srv := &http.Server{
        Addr:         ":443",
        Handler:      mux,
        TLSConfig:    cfg,
        TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
    }
    log.Fatal(srv.ListenAndServeTLS("/etc/ssl/certs/server.crt", "/etc/ssl/certs/server.key"))

}


const apikey = "juu7m6ckf9z3rf2g5a8hwc9577amh26e" 

//func getkey() {
//	b, err := ioutil.ReadFile("/srv/apikey")
//	if err != nil {
//		fmt.Print(err)
//	}
//	var apikey string = string(b)
//}

func getzipurl() {
	u, err := url.Parse("api.ziprecruiter.com/jobs/")
	if err != nil {
	log.Fatal(err)
	}
	u.Scheme = "https"
	u.Host = "api.ziprecruiter.com"
	u.Path = "jobs/v1"
	q := u.Query()
	q.Set("search", job)
	q.Add("&location", locale)
	q.Add("&radius_miles", miles)
	q.Add("&days_ago", days_ago)
	q.Add("&jobs_per_page", perp)
	q.Add("&api_key", apikey)
	u.RawQuery = q.Encode()
	fmt.Println(u)
}




const job = "golang"
const locale = "Michigan"
const miles = "25"
const days_ago = "30"
const perp = "20"

func zipsearch() {
	fmt.Println("Starting...")
	search := http.NewRequest("GET", string getzipurl)
	client := &http.Client{}
	response, err := client.Do(search)
	if err != nil {
		fmt.Printf("HTTP req failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data))
	}
	
	//"https://api.ziprecruiter.com/jobs/v1?search=" + job + "&location=" + locale + "&radius_miles=" + miles + "&days_ago=" + days_ago + "&jobs_per_page=10&page=1&api_key=" + apikey)

}
