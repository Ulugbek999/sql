package app

import (
	"encoding/json"
	"log"
	"net/http"

	
)

func (s *Server) handleRegisterCustomer(writer http.ResponseWriter, request *http.Request) {
	
	log.Print(request.RequestURI)
	log.Print(request.Method)
	log.Print(request.Header)
	log.Print(request.Header.Get("Content-Type"))

	//log.Print(request.FormValue("id"))
	log.Print(request.FormValue("name"))
	log.Print(request.FormValue("phone"))
	log.Print(request.FormValue("password"))
	
	var item *customers.Registration	
	err := json.NewDecoder(request.Body).Decode(&item)
		
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	
	saved, err := s.customersSvc.Register(request.Context(), item)

	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	
	data, err := json.Marshal(saved)

	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)

	if err != nil {
		log.Print(err)
	}

	log.Print(item)

}



func (s *Server) handleCustomersGetToken(writer http.ResponseWriter, request *http.Request) {
	
	log.Print(request.RequestURI)
	log.Print(request.Method)
	log.Print(request.Header)
	log.Print(request.Header.Get("Content-Type"))

	//log.Print(request.FormValue("id"))
	log.Print(request.FormValue("login"))
//	log.Print(request.FormValue("phone"))
	log.Print(request.FormValue("password"))
	
	
	var item *customers.Auth
	

// 	// 
	 err := json.NewDecoder(request.Body).Decode(&item)

	//var customer1 *customers.Customer	

	
		
	if err != nil {
		log.Print("We are hereeeeeeeeeeeeeeeeeeeeeeeeeeeeee5")
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	token, err := s.customersSvc.Token(request.Context(), item.Login, item.Password)
	//item, err := s.securitySvc.TokenForCustomer(request.Context(), customer.Phone, customer.Password)
	
	if err != nil {
		log.Print("We are hereeeeeeeeeeeeeeeeeeeeeeeeeeeeee2")
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	//token := jwt.New(jwt.SigningMethodRS256)

	//dataToken := json.Token(item)
	
	 data, err := json.Marshal(&security.Token{Token: token})

	if err != nil {
		log.Print("We are hereeeeeeeeeeeeeeeeeeeeeeeeeeeeee3")
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	 writer.Header().Set("Content-Type", "application/json")
	 _, err = writer.Write(data)

	if err != nil {
		log.Print("We are hereeeeeeeeeeeeeeeeeeeeeeeeeeeeee4")
		log.Print(err)
	}

	log.Print(item)

}

func (s *Server) handleCustomersGetProducts(writer http.ResponseWriter, request *http.Request) {
	items, err := s.customersSvc.Products(request.Context())
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(items)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Print(err)
	}
}

func (s *Server) handleCustomerGetPurchases(writer http.ResponseWriter, request *http.Request) {
	id, err := middleware.Authentication(request.Context())
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	items, err := s.customersSvc.Purchases(request.Context(), id)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(items)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Print(err)
	}

}


func (s *Server) handleCustomerMakePurchases(writer http.ResponseWriter, request *http.Request) {
	id, err := middleware.Authentication(request.Context())
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	items, err := s.customersSvc.Purchases(request.Context(), id)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(items)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Print(err)
	}

}