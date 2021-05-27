package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"


	
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/docs/v1"
)

func (s *Server) handleRegisterManager(writer http.ResponseWriter, request *http.Request) {
	// _, err := middleware.Authentication(request.Context())
	// if err != nil {
	// 	log.Print(err)
	// 	http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	// 	return
	// }
	log.Print(request.RequestURI)
	log.Print(request.Method)
	log.Print(request.Header)
	log.Print(request.Header.Get("Content-Type"))

	//log.Print(request.FormValue("id"))
	log.Print(request.FormValue("name"))
	log.Print(request.FormValue("phone"))
	log.Print(request.FormValue("password"))
	
	var item *managers.Registration	
	err := json.NewDecoder(request.Body).Decode(&item)
		
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	
	saved, err := s.managersSvc.Register(request.Context(), item)

	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	
	data, err := json.Marshal(&managers.Token{Token: saved})

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



func (s *Server) handleManagersGetToken(writer http.ResponseWriter, request *http.Request) {
	
	log.Print(request.RequestURI)
	log.Print(request.Method)
	log.Print(request.Header)
	log.Print(request.Header.Get("Content-Type"))

	//log.Print(request.FormValue("id"))
	log.Print(request.FormValue("login"))
//	log.Print(request.FormValue("phone"))
	log.Print(request.FormValue("password"))
	
	
	var item *managers.Auth
	

// 	// 
	 err := json.NewDecoder(request.Body).Decode(&item)

	//var customer1 *managers.Customer	

	
		
	if err != nil {
		log.Print("We are hereeeeeeeeeeeeeeeeeeeeeeeeeeeeee5")
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	token, err := s.managersSvc.Token(request.Context(), item.Login, item.Password)
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

func (s *Server) handleManagersGetProducts(writer http.ResponseWriter, request *http.Request) {
	items, err := s.managersSvc.Products(request.Context())
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

func (s *Server) handleManagersChangeProduct(writer http.ResponseWriter, request *http.Request) {
	
	var item *managers.Product	
	err := json.NewDecoder(request.Body).Decode(&item)
		
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	
		
	items, err := s.managersSvc.ChangeProducts(request.Context(), item)
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

func (s *Server) handleManagerGetPurchases(writer http.ResponseWriter, request *http.Request) {
	id, err := middleware.Authentication(request.Context())
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	items, err := s.managersSvc.Purchases(request.Context(), id)
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


func (s *Server) handleManagersMakeSale(writer http.ResponseWriter, request *http.Request) {
	id, err := middleware.Authentication(request.Context())
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	var item *managers.Sale	
	err = json.NewDecoder(request.Body).Decode(&item)
	item.ManagerID = id	
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	
		
	items, err := s.managersSvc.MakeSale(request.Context(), item)
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


func (s *Server) handleManagersGetSale(writer http.ResponseWriter, request *http.Request) {
	id, err := middleware.Authentication(request.Context())
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	var item managers.SaleTotal	
			
	items, err := s.managersSvc.GetSale(request.Context())
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	for _, v := range items {
		if v.ManagerID == id {
			item.ManagerID = id
			item.Total = v.Total
		}
	}
	log.Printf("%v", item)
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


func (s *Server) handleManagerMakePurchases(writer http.ResponseWriter, request *http.Request) {
	id, err := middleware.Authentication(request.Context())
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	items, err := s.managersSvc.Purchases(request.Context(), id)
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


func (s *Server) handleManagersGetGoogleAuth(writer http.ResponseWriter, request *http.Request) {
	//	id, err := middleware.Authentication(request.Context())
	// if err != nil {
	// 	log.Print(err)
	// 	http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	// 	return
	// }

	wd, err := os.Getwd()
	if err != nil {
		log.Print(err)
		return
	}
	wd1 := wd + "/credentials.json" 
	b, err := ioutil.ReadFile(wd1)
	if err != nil {
			log.Fatalf("Unable to read client secret file: %v", err)
	}
	
	
	log.Print(request.RequestURI)
	log.Print(request.Method)
	log.Print(request.Header)
	log.Print(request.Header.Get("Content-Type"))
	log.Print(request.Header.Get("Authorization"))
	codeParam := request.URL.Query().Get("code")
//	log.Print(codeParam)
	//config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/drive.metadata.readonly")
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/documents.readonly")
	if err != nil {
			log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, codeParam)
	if err != nil {
			log.Fatalf("Unable to retrieve token from web: %v", err)
	}

	client := config.Client(context.Background(), tok)

	srv, err := docs.New(client)
	if err != nil {
			log.Fatalf("Unable to retrieve Docs client: %v", err)
	}

	// Prints the title of the requested doc:
	// 

	// var docId string
	// if _, err := fmt.Scan(&docId); err != nil {
	// 		log.Fatalf("Unable to read authorization code: %v", err)
	// }


	docId := "1mXe0p2-W0EEyx6k6hC4RoS95ZDY-PXuDEPF1WRy7MOg"
	doc, err := srv.Documents.Get(docId).Do()
	if err != nil {
			log.Fatalf("Unable to retrieve data from document: %v", err)
	}
	fmt.Printf("The title of the doc is: %s\n", doc.Title)


	
}

func (s *Server) handleManagersGetGoogle(writer http.ResponseWriter, request *http.Request) {
	//	id, err := middleware.Authentication(request.Context())
	// if err != nil {
	// 	log.Print(err)
	// 	http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	// 	return
	// }
	
	wd, err := os.Getwd()
	if err != nil {
		log.Print(err)
		return
	}
	wd1 := wd + "/credentials.json" 
	b, err := ioutil.ReadFile(wd1)
	if err != nil {
			log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	//config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/drive.metadata.readonly")
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/documents.readonly")
	if err != nil {
			log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := docs.New(client)
	if err != nil {
			log.Fatalf("Unable to retrieve Docs client: %v", err)
	}

	// Prints the title of the requested doc:
	// https://docs.google.com/document/d/195j9eDD3ccgjQRttHhJPymLJUCOUjs-jmwTrekvdjFE/edit
	docId := "1mXe0p2-W0EEyx6k6hC4RoS95ZDY-PXuDEPF1WRy7MOg"
	doc, err := srv.Documents.Get(docId).Do()
	if err != nil {
			log.Fatalf("Unable to retrieve data from document: %v", err)
	}
	fmt.Printf("The title of the doc is: %s\n", doc.Title)

	
}


func (s *Server) handleUnites(writer http.ResponseWriter, request *http.Request) {

	var unit *managers.Unites
	err := json.NewDecoder(request.Body).Decode(&unit)
		
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	
	item, err := s.managersSvc.SaveUnit(request.Context(), *unit)

//	item, err := s.customersSvc.ByID(request.Context(), id)

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


func (s *Server) handleUnitesConversion(writer http.ResponseWriter, request *http.Request) {

	var unit *managers.UnitesConversion
	err := json.NewDecoder(request.Body).Decode(&unit)
		
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	
	item, err := s.managersSvc.GetUnitConversion(request.Context(), *unit)

//	item, err := s.customersSvc.ByID(request.Context(), id)

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

// Retrieves a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
			tok = getTokenFromWeb(config)
			saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Requests a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
			"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
			log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, authCode)
	if err != nil {
			log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
			return nil, err
	}
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	defer f.Close()
	if err != nil {
			log.Fatalf("Unable to cache OAuth token: %v", err)
	}
	json.NewEncoder(f).Encode(token)
}


func (s *Server) handleOmdb(writer http.ResponseWriter, request *http.Request) {
	// _, err := middleware.Authentication(request.Context())
	// if err != nil {
	// 	log.Print(err)
	// 	http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	// 	return
	// }
	log.Print(request.RequestURI)
	log.Print(request.Method)
	log.Print(request.Header)
	log.Print(request.Header.Get("Content-Type"))

	//log.Print(request.FormValue("id"))
	log.Print(request.FormValue("title"))
	
	
	var item *managers.Search	
	err := json.NewDecoder(request.Body).Decode(&item)
		
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	
	saved, err := s.managersSvc.SearchOmdb(request.Context(), item)

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


func (s *Server) handleImport(writer http.ResponseWriter, request *http.Request) {

		
	err := s.managersSvc.GetVotes(request.Context())

//	item, err := s.customersSvc.ByID(request.Context(), id)

	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	

}