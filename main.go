package main

import (
	"crypto/tls"
	"log"
	"net/http"
	//"io/ioutil"
	//"fmt"
	//"net/url"
)

//set header for all handlers instead of configuring each one.  Credit Adam Ng https://www.socketloop.com/tutorials/golang-set-or-add-headers-for-many-or-different-handlers
func SetHeaders(w http.ResponseWriter) {
	w.Header().Set("X-Frame-Options", "SAMEORIGIN")
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
	w.Header().Set("Allow", "OPTIONS, GET, HEAD, POST")
	w.Header().Set("Content-Security-Policy", "default-src 'none'; font-src 'https://fonts.googleapis.com'; img-src 'self' https://i.imgur.com; object-src 'none'; script-src 'self'; style-src 'self'; frame-ancestors 'self'")
	w.Header().Set("Set-Cookie", "__Host-BMOSESSIONID=YnVnemlsbGE=; Max-Age=2592000; Path=/; Secure; HttpOnly; SameSite=Strict")
	w.Header().Set("Referrer-Policy", "no-referrer, strict-origin-when-cross-origin")
	w.Header().Set("X-Content-Type-Options","nosniff")
	w.Header().Set("X-XSS-Protection", "1; mode=block")

}




func redirectTLS(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://"+r.Host+":443"+r.RequestURI, http.StatusMovedPermanently)
}


func index(w http.ResponseWriter, req *http.Request) {
	SetHeaders(w)	
//	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
//        w.Write([]byte("This is an example server.\n"))
	if req.URL.Path != "/" {
		log.Printf("404: %s", req.URL.String())	
		http.NotFound(w, req)
		return
	}
	http.ServeFile(w, req, "index.html")
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	http.ServeFile(w, r, "favicon.png")
}
func art(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	http.ServeFile(w, r, "art.png")
}


func main() {
	//lets start an unencrypted http server, for the purposes of redirect
	go http.ListenAndServe(":80", http.HandlerFunc(redirectTLS))
	
	
	//mux up a router
	mux := http.NewServeMux()
	mux.HandleFunc("/favicon.png", faviconHandler)	
	mux.HandleFunc("/art.png", art)
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

