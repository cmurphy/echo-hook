package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	admission "k8s.io/api/admission/v1beta1"
	"net/http"
)

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(string(body))
	admissionRequest := admission.AdmissionReview{}
	if err := json.Unmarshal(body, &admissionRequest); err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	admissionResponse := admission.AdmissionReview{
		Response: &admission.AdmissionResponse{
			Allowed: true,
			UID:     admissionRequest.Request.UID,
		},
	}
	admissionResponse.TypeMeta.Kind = "AdmissionReview"
	admissionResponse.TypeMeta.APIVersion = "admission.k8s.io/v1"
	httpResponse, err := json.Marshal(admissionResponse)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if _, err := w.Write(httpResponse); err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	certFile := "cert.pem"
	keyFile := "key.pem"
	pair, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		fmt.Printf("Failed to load key pair: %v\n", err)
		return
	}

	server := &http.Server{
		Addr:      ":8080",
		TLSConfig: &tls.Config{Certificates: []tls.Certificate{pair}},
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", webhookHandler)
	server.Handler = mux
	server.ListenAndServeTLS("", "")
}
