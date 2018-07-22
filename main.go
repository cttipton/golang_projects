package main

import (
	"crypto/tls"
	"log"
	"net/http"
	//"io/ioutil"
	//"fmt"
	//"net/url"
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

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "favicon.png")
}


func main() {
	//lets start an unencrypted http server, for the purposes of redirect
	go http.ListenAndServe(":80", http.HandlerFunc(redirectTLS))
	
	
	//mux up a router
	mux := http.NewServeMux()
	mux.HandleFunc("/favicon.png", faviconHandler)	
	//make sure we can renew LetsEncrypt Certs. Below handle will allow LetsEncrypt and our webserver to verify one another.
	mux.Handle("/.well-known/acme-challenge/", http.StripPrefix("/.well-known/acme-challenge/", http.FileServer(http.FileSystem(http.Dir("/tmp/letsencrypt/")))))

	//Point us home
  mux.HandleFunc("/", index) //func(w http.ResponseWriter, req *http.Request) {
    

	//webserver config, this gets A+ on Qualys scan if you're sure to bind to 443 with full cert chain--not just server.crt.
	cfg := &tls.Config{
        	MinVersion:               tls.VersionTLS12,
        	CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
        	PreferServerCipherSuites: true,
        	CipherSuites: []uint16{
        	    tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
        	    tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
        	},
    	}
	srv := &http.Server{
        	Addr:         ":443",
        	Handler:      mux,
        	TLSConfig:    cfg,
        	TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
    	}
	log.Println(srv.ListenAndServeTLS("/etc/ssl/iminshell.com.pem", "/etc/ssl/iminshell.com.key"))

}

