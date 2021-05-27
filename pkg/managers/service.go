
package managers

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	
	"time"

)

//ErrNoSuchUser for
var ErrNoSuchUser = errors.New("no such user")

//ErrNotFound for
var ErrNotFound = errors.New("item not found")

//ErrInternal for
var ErrInternal = errors.New("internal error")

//ErrNoRows for
var ErrNoRows = errors.New("No rows")

//ErrPhoneUsed for
var ErrPhoneUsed = errors.New("phone already registered")

//ErrInvalidPassword for
var ErrInvalidPassword = errors.New("invalid password")

//ErrTokenNotFound for
var ErrTokenNotFound = errors.New("token not found")

//ErrTokenExpired for
var ErrTokenExpired = errors.New("token expired")

// Service npenctasnset co6oi cepsuc no ynpasnenwo OaHHepamn.
type Service struct {
	pool *pgxpool.Pool
}

// NewService co3qaét cepsuc.
func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}

}

// Manager npenctasnaet codoi GaHHep.
type Manager struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
	//	Password string 	`json:"password"`
	Active  bool      `json:"active"`
	Created time.Time `json:"created"`
}

//Registration for
type Registration struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
	//	Password  	string 	`json:"password"`
	Roles []string `json:"roles"`
}

// Unites npenctasnaet codoi GaHHep.
type Unites struct {
	ID      int64     `json:"id"`
	Name    string    `json:"name"`
	Meter   float64   `json:"meter"`
	Created time.Time `json:"created"`
}

// UnitesConversion npenctasnaet codoi GaHHep.
type UnitesConversion struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Value  float64 `json:"value"`
	Result float64 `json:"result"`
}

// Votes npenctasnaet codoi GaHHep.
type Votes struct {
	VoterID     int64 `json:"voterid"`
	CandidateID int64 `json:"candidateid"`
}

//Search for
type Search struct {
	Title string `json:"title"`
}

// Register for
func (s *Service) Register(ctx context.Context, registration *Registration) (token string, err error) {

	item := &Manager{}

	// hash, err := bcrypt.GenerateFromPassword([]byte(registration.Password), bcrypt.DefaultCost)
	// if err != nil {
	// 	log.Print(err)
	// 	return nil, ErrInternal
	// }

	err = s.pool.QueryRow(ctx, `
		INSERT INTO Managers(name, login, roles)
		VALUES($1, $2, $3)
		ON CONFLICT (login) DO NOTHING RETURNING id, name, login, active, created
		`, registration.Name, registration.Phone, registration.Roles).Scan(
		&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)

	if errors.Is(err, pgx.ErrNoRows) {
		log.Print("No rows")
		return "", ErrNoRows
	}

	if err != nil {
		log.Print(err)
		return "", ErrInternal
	}

	log.Print(item)

	buffer := make([]byte, 256)
	n, err := rand.Read(buffer)
	if n != len(buffer) || err != nil {
		return "", ErrInternal
	}

	token = hex.EncodeToString(buffer)

	_, err = s.pool.Exec(ctx, `
		INSERT INTO managers_tokens(token, manager_id) VALUES($1, $2)`, token, item.ID)

	if err != nil {
		log.Print(err)
		return "", ErrInternal
	}

	return token, nil

}

//Auth for
type Auth struct {
	Login    string `json:"phone"`
	Password string `json:"password"`
}

//Token for
type Token struct {
	Token string `json:"token"`
}

//Token for
func (s *Service) Token(ctx context.Context, phone string, password string) (token string, err error) {
	var hash string
	var id int64
	// hashPas, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	// log.Print(hex.EncodeToString(hashPas))
	err = s.pool.QueryRow(ctx, `SELECT id, password FROM managers WHERE login = $1`, phone).Scan(&id, &hash)

	if err == pgx.ErrNoRows {
		return "", ErrNoSuchUser
	}

	if err != nil {
		return "", ErrInternal
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {

		return "", ErrInvalidPassword
	}

	buffer := make([]byte, 256)
	n, err := rand.Read(buffer)
	if n != len(buffer) || err != nil {
		return "", ErrInternal
	}

	token = hex.EncodeToString(buffer)

	_, err = s.pool.Exec(ctx, `
		INSERT INTO managers_tokens(token, manager_id) VALUES($1, $2)`, token, id)

	if err != nil {
		log.Print(err)
		return "", ErrInternal
	}

	return token, nil
}

//TokenNew for
func (s *Service) TokenNew(ctx context.Context, phone string, password string) (token string, err error) {
	var hash string
	var id int64
	// hashPas, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	// log.Print(hex.EncodeToString(hashPas))
	err = s.pool.QueryRow(ctx, `SELECT id, password FROM managers WHERE login = $1`, phone).Scan(&id, &hash)

	if err == pgx.ErrNoRows {
		return "", ErrNoSuchUser
	}

	// if err != nil {
	// 	return "", ErrInternal
	// }

	// err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	// if err != nil {

	// 	return "", ErrInvalidPassword
	// }

	buffer := make([]byte, 256)
	n, err := rand.Read(buffer)
	if n != len(buffer) || err != nil {
		return "", ErrInternal
	}

	token = hex.EncodeToString(buffer)

	_, err = s.pool.Exec(ctx, `
		INSERT INTO managers_tokens(token, manager_id) VALUES($1, $2)`, token, id)

	if err != nil {
		log.Print(err)
		return "", ErrInternal
	}

	return token, nil
}

//Sale for
type Sale struct {
	ID         int64          `json:"id"`
	CustomerID string         `json:"customer_id"`
	Positions  []SalePosition `json:"positions"`
	ManagerID  int64          `json:"manager_id"`
}

//SalePosition for
type SalePosition struct {
	ID        int64 `json:"id"`
	ProductID int64 `json:"product_id"`
	Price     int   `json:"price"`
	Qty       int   `json:"qty"`
}

//NewNullString for
func NewNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

//MakeSale for
func (s *Service) MakeSale(ctx context.Context, itemSale *Sale) (*Sale, error) {
	products, _ := s.Products(ctx)
	name := ""
	item := &Sale{}
	saleItem := &SalePosition{}

	for _, v := range itemSale.Positions {
		for _, product := range products {
			if product.ID == v.ProductID {

				quantity := product.Qty
				if quantity < v.Qty {
					log.Print("Not enough quantity")
					return nil, ErrInternal
				}
			}
		}
	}

	if itemSale.ID == 0 {

		err := s.pool.QueryRow(ctx, `
		INSERT INTO sales(manager_id, customer_id) VALUES($1, $2) RETURNING id, manager_id
		`, itemSale.ManagerID, NewNullString(itemSale.CustomerID)).Scan(&item.ID, &item.ManagerID)

		if errors.Is(err, pgx.ErrNoRows) {
			log.Print("No rows")
			return nil, ErrNoRows
		}

		if err != nil {
			log.Print(err)
			return nil, ErrInternal
		}

		for _, v := range itemSale.Positions {
			for _, product := range products {
				if product.ID == v.ProductID {
					name = product.Name
					quantity := product.Qty
					if quantity < v.Qty {
						log.Print("Not enough quantity")
						return nil, ErrInternal
					}
				}
			}
			err := s.pool.QueryRow(ctx, `
			INSERT INTO sale_positions(sale_id, product_id, price, qty, name) VALUES($1, $2, $3, $4, $5) RETURNING id, product_id
			`, item.ID, v.ProductID, v.Price, v.Qty, name).Scan(&saleItem.ID, &saleItem.ProductID)

			if errors.Is(err, pgx.ErrNoRows) {
				log.Print("No rows")
				return nil, ErrNoRows
			}

			if err != nil {
				log.Print(err)
				return nil, ErrInternal
			}

		}

	}

	// if itemSale.ID != 0 {

	// 	err := s.pool.QueryRow(ctx, `
	// 	UPDATE products SET name = $2, qty = $3, price = $4 WHERE id = $1 RETURNING id, name, price, qty
	// 	`, itemSale.ID, itemSale.Name, itemSale.Qty, itemProduct.Price).Scan(&item.ID, &item.Name, &item.Price, &item.Qty)

	// 	if errors.Is(err, pgx.ErrNoRows) {
	// 		log.Print("No rows")
	// 		return nil, ErrNoRows
	// 	}

	// 	if err != nil {
	// 		log.Print(err)
	// 		return nil, ErrInternal
	// 	}
	// }

	log.Print(saleItem)

	return item, nil

}

//SaleTotal for
type SaleTotal struct {
	ManagerID int64 `json:"manager_id"`
	Total     int   `json:"total"`
}

//GetSale for
func (s *Service) GetSale(ctx context.Context) ([]*SaleTotal, error) {
	items := make([]*SaleTotal, 0)

	rows, err := s.pool.Query(ctx, `
		SELECT  
		m.id,				
		CASE 
			WHEN sum(ss.total) IS NULL THEN 0 
			ELSE sum(ss.total) 
		END total
		FROM managers AS m  
		LEFT JOIN sales AS s on m.id = s.manager_id  
		LEFT JOIN (
		SELECT  sp.sale_id,
				sum(sp.price * sp.qty) AS total
		FROM sale_positions AS sp
		GROUP BY sp.sale_id
		) ss ON s.id = ss.sale_id
		GROUP BY 
		m.id 
	`)
	if errors.Is(err, pgx.ErrNoRows) {
		return items, nil
	}
	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}
	defer rows.Close()
	for rows.Next() {
		item := &SaleTotal{}
		err = rows.Scan(&item.ManagerID, &item.Total)
		if err != nil {
			log.Print(err)
			return nil, err
		}
		items = append(items, item)
	}
	err = rows.Err()
	if err != nil {
		log.Print(err)
		return nil, err
	}

	log.Printf("%v", items)
	return items, nil

}

//Product for
type Product struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Qty   int    `json:"qty"`
}

//Products for
func (s *Service) Products(ctx context.Context) ([]*Product, error) {
	items := make([]*Product, 0)
	rows, err := s.pool.Query(ctx, `
		SELECT id, name, price, qty FROM products WHERE active ORDER BY id LIMIT 500 
	`)
	if errors.Is(err, pgx.ErrNoRows) {
		return items, nil
	}
	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}
	defer rows.Close()
	for rows.Next() {
		item := &Product{}
		err = rows.Scan(&item.ID, &item.Name, &item.Price, &item.Qty)
		if err != nil {
			log.Print(err)
			return nil, err
		}
		items = append(items, item)
	}
	err = rows.Err()
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return items, nil
}

//ChangeProducts for
func (s *Service) ChangeProducts(ctx context.Context, itemProduct *Product) (*Product, error) {

	item := &Product{}

	if itemProduct.ID == 0 {

		err := s.pool.QueryRow(ctx, `
		INSERT INTO products(name, qty, price) VALUES($1, $2, $3) RETURNING id, name, price, qty
		`, itemProduct.Name, itemProduct.Qty, itemProduct.Price).Scan(&item.ID, &item.Name, &item.Price, &item.Qty)

		if errors.Is(err, pgx.ErrNoRows) {
			log.Print("No rows")
			return nil, ErrNoRows
		}

		if err != nil {
			log.Print(err)
			return nil, ErrInternal
		}

	}

	if itemProduct.ID != 0 {

		err := s.pool.QueryRow(ctx, `
		UPDATE products SET name = $2, qty = $3, price = $4 WHERE id = $1 RETURNING id, name, price, qty
		`, itemProduct.ID, itemProduct.Name, itemProduct.Qty, itemProduct.Price).Scan(&item.ID, &item.Name, &item.Price, &item.Qty)

		if errors.Is(err, pgx.ErrNoRows) {
			log.Print("No rows")
			return nil, ErrNoRows
		}

		if err != nil {
			log.Print(err)
			return nil, ErrInternal
		}
	}

	return item, nil

}

//IDByToken for
func (s *Service) IDByToken(ctx context.Context, token string) (int64, error) {
	var id int64
	err := s.pool.QueryRow(ctx, `
	 SELECT Manager_id FROM Managers_tokens WHERE token = $1
	 `, token).Scan(&id)

	if err == pgx.ErrNoRows {
		return 0, nil
	}

	if err != nil {
		return 0, ErrInternal
	}
	log.Print("ID of Manager ", id)
	return id, nil

}

//ManagerRole for
func (s *Service) ManagerRole(ctx context.Context, roles ...string) bool {
	id, err := middleware.Authentication(ctx)

	if err != nil {
		log.Print(err)
		return false
	}
	err = s.pool.QueryRow(ctx, `
	 SELECT roles FROM managers WHERE id = $1
	 `, id).Scan(&roles)

	if err == pgx.ErrNoRows {
		return false
	}

	if err != nil {
		return false
	}
	for _, v := range roles {
		if v == "ADMIN" {
			return true
		}
	}
	return false

}

//Purchase for
type Purchase struct {
	ID         int64 `json:"id"`
	ManagerID  int   `json:"manager_id"`
	CustomerID int   `json:"customer_id"`
}

//Purchases for
func (s *Service) Purchases(ctx context.Context, id int64) ([]*Purchase, error) {
	items := make([]*Purchase, 0)
	rows, err := s.pool.Query(ctx, `
		SELECT id, manager_id, Manager_id FROM sales WHERE Manager_id = $1 ORDER BY id LIMIT 500 
	`, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return items, nil
	}
	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}
	defer rows.Close()
	for rows.Next() {
		item := &Purchase{}
		err = rows.Scan(&item.ID, &item.ManagerID, &item.ManagerID)
		if err != nil {
			log.Print(err)
			return nil, err
		}
		items = append(items, item)
	}
	err = rows.Err()
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return items, nil
}

// SearchOmdb for
func (s *Service) SearchOmdb(ctx context.Context, search *Search) (result interface{}, err error) {

	//item := &Search{}

	url := "http://www.omdbapi.com/?s=" + search.Title + "&apikey=4de9f5a6"

	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(response.StatusCode)

	var itemResponse interface{}
	err = json.NewDecoder(response.Body).Decode(&itemResponse)

	if err != nil {
		log.Print(err)
		return nil, err
	}

	defer response.Body.Close()
	log.Print(response.Body)
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(data))

	return itemResponse, nil

}

// SaveUnit for
func (s *Service) SaveUnit(ctx context.Context, itemUnit Unites) (*Unites, error) {

	item := &Unites{}

	if itemUnit.ID == 0 {

		err := s.pool.QueryRow(ctx, `
		INSERT INTO unites(name, meter) VALUES($1, $2) RETURNING id, name, meter, created
		`, itemUnit.Name, itemUnit.Meter).Scan(&item.ID, &item.Name, &item.Meter, &item.Created)

		if errors.Is(err, pgx.ErrNoRows) {
			log.Print("No rows")
			return nil, ErrNoRows
		}

		if err != nil {
			log.Print(err)
			return nil, ErrInternal
		}

		// lenCustomers := len(s.items) - 1
		// for i, Customer := range s.items {
		// 	if i == lenCustomers {
		// 		lastID = Customer.ID
		// 	}
		// }

	}

	if itemUnit.ID != 0 {

		err := s.pool.QueryRow(ctx, `
		UPDATE unites SET name = $2, meter = $3 WHERE id = $1 RETURNING id, name, meter, created
		`, itemUnit.ID, itemUnit.Name, itemUnit.Meter).Scan(&item.ID, &item.Name, &item.Meter, &item.Created)

		if errors.Is(err, pgx.ErrNoRows) {
			log.Print("No rows")
			return nil, ErrNoRows
		}

		if err != nil {
			log.Print(err)
			return nil, ErrInternal
		}
	}
	return item, nil

}

// GetUnitConversion for
func (s *Service) GetUnitConversion(ctx context.Context, itemUnit UnitesConversion) (result float64, err error) {

	//	item := &UnitesConversion{}

	var from float64
	var to float64

	err = s.pool.QueryRow(ctx, `
	 SELECT meter FROM unites WHERE name = $1
	 `, itemUnit.From).Scan(&from)

	if err == pgx.ErrNoRows {
		return 0, nil
	}

	if err != nil {
		return 0, ErrInternal
	}

	err = s.pool.QueryRow(ctx, `
	 SELECT meter FROM unites WHERE name = $1
	 `, itemUnit.To).Scan(&to)

	if err == pgx.ErrNoRows {
		return 0, nil
	}

	if err != nil {
		return 0, ErrInternal
	}

	result = from / to

	return result, nil

}

// GetVotes for
func (s *Service) GetVotes(ctx context.Context) (err error) {

	votes, err := Import()
	if err != nil {
		log.Print(err)
		return ErrInternal
	}

	for _, vote := range votes {

		_, err = s.pool.Exec(ctx, `
		INSERT INTO votes(voter_id, candidate_id) VALUES($1, $2)`, vote.VoterID, vote.CandidateID)

		if err != nil {
			log.Print(err)
			return ErrInternal
		}

	}

	return nil

}

// Import for
func Import() ([]*Votes, error) {

	items := make([]*Votes, 0)

	dir, err := os.Getwd()
	if err != nil {
		log.Print(err)
		return nil, err
	}

	dirImport := dir + "/votes.dump"
	fileImport, err := os.Open(dirImport)
	if err != nil {
		log.Print(err)
		err = ErrNotFound
	}
	if err != ErrNotFound {
		defer func() {
			err := fileImport.Close()
			if err != nil {
				log.Print(err)
			}
		}()

		//	log.Printf("%#v", fileImport)

		contentFile := make([]byte, 0)
		bufFavorite := make([]byte, 4)
		for {
			read, err := fileImport.Read(bufFavorite)
			if err == io.EOF {
				break
			}
			contentFile = append(contentFile, bufFavorite[:read]...)
		}

		dataFavorite := string(contentFile)
		newDataFavorite := strings.Split(dataFavorite, "\r\n")
		//log.Print(data)
		//log.Print(newData)

		for _, stroka := range newDataFavorite {
			item := &Votes{}

			newData2 := strings.Split(stroka, ",")
			for ind, stroka2 := range newData2 {

				if ind == 0 {
					voterID, _ := strconv.ParseInt(stroka2, 10, 64)
					item.VoterID = int64(voterID)

				}
				if ind == 1 {
					
					candidateID, err := strconv.ParseInt(stroka2, 10, 64)
					if err != nil {
						log.Print(err)
					}
					item.CandidateID = int64(candidateID)
				}

			}

			items = append(items, item)

		}

	}
	return items, nil
}