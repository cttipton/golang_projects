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
	go getkey()
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

func getkey() {
	b, err := ioutil.ReadFile("/srv/apikey")
	if err != nil {
		fmt.Print(err)
	}
	apikey := string(b)
}


func getzipurl() {
		u, err := url.Parse("http://bing.com/search?q=dotnet")
	if err != nil {
		log.Fatal(err)
	}
	u.Scheme = "https"
	u.Host = "api.ziprecruiter.com"
	u.Path = "jobs/"
	q := u.Query()
	q.Set("v1?search", getjob())
	u.RawQuery = q.Encode()
	fmt.Println(u)
}

func getjob() {
	Println("golang")
}


func zipsearch() {
	fmt.Println("Starting...")
//	jsdonData := map[

	go getzipurl()
	http.NewRequest("GET", "golang")
	
	//"https://api.ziprecruiter.com/jobs/v1?search=" + job + "&location=" + locale + "&radius_miles=" + miles + "&days_ago=" + days_ago + "&jobs_per_page=10&page=1&api_key=" + apikey)

}
