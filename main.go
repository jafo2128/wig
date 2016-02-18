package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	_version string = "v0.1"
)

var _settings = LoadSettings()

type Server struct {
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Upgrade") == "websocket" {
		fmt.Printf("New client: %s\n", r.RemoteAddr)
		NewClient(w, r)
	} else {
		body := "Hello World\n"
		w.Header().Set("Server", "wig")
		w.Header().Set("Connection", "close")
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Length", fmt.Sprint(len(body)))
		fmt.Fprint(w, body)
	}
}

func main() {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)
	signal.Notify(sigchan, syscall.SIGTERM)

	addr := fmt.Sprintf("%v:%v", _settings.WsHost, _settings.WsPort)

	//check if certs exist, otherwise generate new cert
	if _, err := os.Stat(_settings.SslCert); os.IsNotExist(err) {
		if _settings.AutoGenCert {
			//generate cert
			serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
			serial, err := rand.Int(rand.Reader, serialNumberLimit)
			bits := 2048
			priv, err := rsa.GenerateKey(rand.Reader, bits)
			if err == nil {
				cert := x509.Certificate{
					SerialNumber: serial,
					Subject: pkix.Name{
						CommonName: _settings.WsHost,
					},
					DNSNames:              []string{_settings.WsHost},
					NotBefore:             time.Now(),
					NotAfter:              time.Now().Add(time.Hour * 24 * 365),
					KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDataEncipherment,
					ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
					BasicConstraintsValid: true,
				}

				derBytes, ccer := x509.CreateCertificate(rand.Reader, &cert, &cert, &priv.PublicKey, priv)
				if ccer == nil {
					cout, coer := os.Create(_settings.SslCert)
					if coer == nil {
						pem.Encode(cout, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
						cout.Close()

						keyOut, koer := os.OpenFile(_settings.SslKey, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
						if koer == nil {
							pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
							keyOut.Close()

							fmt.Printf("***** ATTENTION *****\n\tA certificate was generated automatically, please visit https://%s before opening site.\n***** ATTENTION *****\n\n", addr)
						} else {
							fmt.Println("Failed to write key:", koer)
							os.Exit(99)
						}
					} else {
						fmt.Println("Failed to write cert:", coer)
						os.Exit(99)
					}
				} else {
					fmt.Println("Failed to create cert:", ccer)
					os.Exit(99)
				}
			} else {
				fmt.Println("Failed to create key:", err)
				os.Exit(99)
			}
		} else {
			fmt.Println("SSL certificate was not found:", _settings.SslCert)
			os.Exit(99)
		}
	}

	server := Server{}
	go func() {
		http.HandleFunc("/ws", server.ServeHTTP)
		hter := http.ListenAndServeTLS(addr, _settings.SslCert, _settings.SslKey, nil)
		if hter != nil {
			fmt.Errorf(hter.Error())
			os.Exit(99)
		}
	}()

	fmt.Println("WS port started on:", addr)
	<-sigchan
}
