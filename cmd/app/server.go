package app

import (
	"encoding/json"
	"io/ioutil"

	//	"io"

	"log"
	"net/http"

	//	"os"
	"strconv"
	//	"strings"

	

)

const (
	//GET for
	GET = "GET"
	//POST for
	POST = "POST"
	//DELETE for
	DELETE = "DELETE"
)

// Server npegctasnseT coOow normyeckwi CepBep Hawero npunomeHna.
type Server struct {
	mux *mux.Router

	customersSvc *customers.Service

	managersSvc *managers.Service
	
	securitySvc *security.Service
}

// NewServer - OyHKUMA-KOHCTpykTOp pina co3maHna cepsepa.
func NewServer(mux *mux.Router, customersSvc *customers.Service, managersSvc *managers.Service, securitySvc *security.Service) *Server {

	return &Server{mux: mux, customersSvc: customersSvc, managersSvc: managersSvc, securitySvc: securitySvc}

}

//ServeHTTP for
func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	//request.SetBasicAuth("Login1", "password1")
	s.mux.ServeHTTP(writer, request)

}

// Init wHmunannsupyet cepsep (permctpupyet sce Handler's)
func (s *Server) Init() {

	customersAuthenticateMd := middleware.Authenticate(s.customersSvc.IDByToken)
	//customersRoleMd := middleware.CheckRole()
	customersSubrouter := s.mux.PathPrefix("/api/customers").Subrouter()
	customersSubrouter.Use(customersAuthenticateMd)
	
	s.mux.HandleFunc("/customers/{id}/block", s.handleCustomerblockByID).Methods(POST)
	//s.mux.HandleFunc("/customers.unblockById", s.handleCustomerunblockByID)
	s.mux.HandleFunc("/customers/{id}/block", s.handleCustomerunblockByID).Methods(DELETE)
	//s.mux.HandleFunc("/customers.removeById", s.handleCustomerremoveByID)
	s.mux.HandleFunc("/customers/{id}", s.handleCustomerRemoveByID).Methods(DELETE)
	
	customersSubrouter.HandleFunc("", s.handleRegisterCustomer).Methods(POST)
	
	customersSubrouter.HandleFunc("/token", s.handleCustomersGetToken).Methods(POST)
	//s.mux.HandleFunc("/api/customers/token/validate", s.handleCheckToken).Methods(POST)
	customersSubrouter.HandleFunc("/purchases", s.handleCustomerGetPurchases).Methods(GET)
	customersSubrouter.HandleFunc("/products", s.handleCustomersGetProducts).Methods(GET)
	//customersSubrouter.HandleFunc("/purchases", s.handleCustomerMakePurchases).Methods(POST)
	//s.mux.HandleFunc("/api/customers/products", s.handleCustomersGetProducts).Methods(GET)

	managersAuthenticateMd := middleware.Authenticate(s.managersSvc.IDByToken)
	chMd := middleware.CheckRole(s.managersSvc.ManagerRole)
	//customersRoleMd := middleware.CheckRole()
	managersSubrouter := s.mux.PathPrefix("/api/managers").Subrouter()
	managersSubrouter.Use(managersAuthenticateMd)
	managersSubrouter.Handle("", chMd(http.HandlerFunc(s.handleRegisterManager))).Methods(POST)
	managersSubrouter.HandleFunc("/token", s.handleManagersGetToken).Methods(POST)
	managersSubrouter.HandleFunc("/products", s.handleManagersChangeProduct).Methods(POST)
	managersSubrouter.HandleFunc("/sales", s.handleManagersMakeSale).Methods(POST)
	managersSubrouter.HandleFunc("/sales", s.handleManagersGetSale).Methods(GET)

	managersSubrouter.HandleFunc("/googleAuth", s.handleManagersGetGoogleAuth).Methods(GET)


	omdbSubrouter := s.mux.PathPrefix("/api/omdb").Subrouter()
	
	omdbSubrouter.HandleFunc("", s.handleOmdb).Methods(GET)

	s.mux.HandleFunc("/google", s.handleManagersGetGoogle).Methods(GET)
	
	s.mux.HandleFunc("/unites", s.handleUnites).Methods(POST)
	s.mux.HandleFunc("/unites", s.handleUnitesConversion).Methods(GET)
	s.mux.HandleFunc("/import", s.handleImport).Methods(GET)

}



func (s *Server) handleAuth(writer http.ResponseWriter, request *http.Request) {

}



func (s *Server) handleGetCustomerByID(writer http.ResponseWriter, request *http.Request) {

	//idParam := request.URL.Query().Get("id")
	idParam, ok := mux.Vars(request)["id"]
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idParam, 10, 64)

	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	item, err := s.customersSvc.ByID(request.Context(), id)

	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(item)

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

func (s *Server) handleCustomerblockByID(writer http.ResponseWriter, request *http.Request) {

	//idParam := request.URL.Query().Get("id")
	idParam, ok := mux.Vars(request)["id"]
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idParam, 10, 64)

	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	item, err := s.customersSvc.BlockByID(request.Context(), id)

	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(item)

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


func (s *Server) handleCustomerunblockByID(writer http.ResponseWriter, request *http.Request) {

	//idParam := request.URL.Query().Get("id")

	idParam, ok := mux.Vars(request)["id"]
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idParam, 10, 64)

	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	item, err := s.customersSvc.UnBlockByID(request.Context(), id)

	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(item)

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


func (s *Server) handleCustomerRemoveByID(writer http.ResponseWriter, request *http.Request) {

	//idParam := request.URL.Query().Get("id")
	idParam, ok := mux.Vars(request)["id"]
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idParam, 10, 64)

	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	item, err := s.customersSvc.RemoveByID(request.Context(), id)

	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(item)

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


func (s *Server) handleSaveCustomer(writer http.ResponseWriter, request *http.Request) {
	
	log.Print(request.RequestURI)
	log.Print(request.Method)
	log.Print(request.Header)
	log.Print(request.Header.Get("Content-Type"))

	log.Print(request.FormValue("id"))
	log.Print(request.FormValue("name"))
	log.Print(request.FormValue("phone"))
	
	
	// body, err := ioutil.ReadAll(request.Body)
	// if err != nil {
	// 	log.Print(err)
	// }
	// log.Printf("%s", body)

	// err = request.ParseMultipartForm(10 * 1024 * 1024)
	// if err != nil {
	// 	log.Print(err)
	// }

	// log.Print(request.Form)
	// log.Print(request.PostForm)
	// idParam := request.FormValue("id")
	
	// id, err := strconv.ParseInt(idParam, 10, 64)
	// if err != nil {
	// 	log.Print(err)

	// }

	// nameParam := request.FormValue("name")
	// phoneParam := request.FormValue("phone")
	
	// customer := customers.Customer{
	// 	ID: id,

	// 	Name: nameParam,

	// 	Phone: phoneParam,
		
	// }

	var customer *customers.Customer
	err := json.NewDecoder(request.Body).Decode(&customer)

	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	item, err := s.customersSvc.Save(request.Context(), *customer)

	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	
	data, err := json.Marshal(item)

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


func (s *Server) handleGetAllcustomers(writer http.ResponseWriter, request *http.Request) {
	log.Print(request)
	log.Print(request.Header)
	log.Print(request.Body)
	item, err := s.customersSvc.All(request.Context())

	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(item)

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

	log.Printf("%#v", item)
}

func (s *Server) handleGetAllActivecustomers(writer http.ResponseWriter, request *http.Request) {
	log.Print(request)
	log.Print(request.Header)
	log.Print(request.Body)
	item, err := s.customersSvc.AllActive(request.Context())

	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(item)

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

	log.Printf("%#v", item)
}



func (s *Server) handleGetToken1(writer http.ResponseWriter, request *http.Request) {
	
	log.Print(request.RequestURI)
	log.Print(request.Method)
	log.Print(request.Header)
	log.Print(request.Header.Get("Content-Type"))

	//log.Print(request.FormValue("id"))
	log.Print(request.FormValue("name"))
	log.Print(request.FormValue("phone"))
	log.Print(request.FormValue("password"))

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Print(err)
	}
	log.Printf("%s", body)

	var customer *customers.Auth	
// 	// 
	 err = json.NewDecoder(request.Body).Decode(&customer)
		
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	// log.Print(customer.Name)
	item, err := s.customersSvc.Token(request.Context(), customer.Login, customer.Password)

	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	
	

	log.Print(item)

}



func (s *Server) handleCheckToken(writer http.ResponseWriter, request *http.Request) {
	
	log.Print(request.RequestURI)
	log.Print(request.Method)
	log.Print(request.Header)
	log.Print(request.Header.Get("Content-Type"))

	//log.Print(request.FormValue("id"))
	log.Print(request.FormValue("token"))
//	log.Print(request.FormValue("phone"))
	//log.Print(request.FormValue("password"))
	
	// body, err := ioutil.ReadAll(request.Body)
	// if err != nil {
	// 	log.Print(err)
	// }
	// log.Printf("%s", body)

	// err = request.ParseForm()
	// if err != nil {
	// 	log.Print(err)
	// }

	// body, err := ioutil.ReadAll(request.Body)
	// if err != nil {
	// 	log.Print(err)
	// }
	// log.Printf("%s", body)

	//var customer *customers.Customer
	type tokeNew struct {
		Token 		string		`json:"token"`	
	}

	var customer1 *tokeNew
	

// 	// 
	 err := json.NewDecoder(request.Body).Decode(&customer1)

	//var customer1 *customers.Customer	

	
		
	if err != nil {
		log.Print("We are hereeeeeeeeeeeeeeeeeeeeeeeeeeeeee5")
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	item, status := s.securitySvc.CheckTokenCustomer(request.Context(), customer1.Token)
	//item, err := s.securitySvc.TokenForCustomer(request.Context(), customer.Phone, customer.Password)
	
	// if item == nil {
	// 	log.Print("We are hereeeeeeeeeeeeeeeeeeeeeeeeeeeeee2")
	// 	//log.Print(err)
	// 	http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	// 	return
	// }

	//token := jwt.New(jwt.SigningMethodRS256)

	//dataToken := json.Token(item)
	
	
	 data, err := json.Marshal(item)

	if err != nil {
		log.Print("We are hereeeeeeeeeeeeeeeeeeeeeeeeeeeeee3")
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if status == 404 {
		writer.WriteHeader(http.StatusNotFound)
	} else if status == 400 {
		writer.WriteHeader(http.StatusBadRequest)
	} else {
		writer.WriteHeader(http.StatusOK)
	}
	 writer.Header().Set("Content-Type", "application/json")
	 _, err = writer.Write(data)

	if err != nil {
		log.Print("We are hereeeeeeeeeeeeeeeeeeeeeeeeeeeeee4")
		log.Print(err)
	}

	log.Print(item)

}