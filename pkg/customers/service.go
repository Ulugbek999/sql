package customers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
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

// Customer npenctasnaet codoi GaHHep.
type Customer struct {
	ID 		int64		`json:"id"`
	Name 	string		`json:"name"`
	Phone 	string		`json:"phone"`	
//	Password string 	`json:"password"`
	Active 	bool		`json:"active"`
	Created time.Time 	`json:"created"`
}

//Registration for
type Registration struct {
	Name		string	`json:"name"`
	Phone		string	`json:"phone"`
	Password	string	`json:"password"`
}


// ByID Bo3BpawaeT OaHHep no upeHTHOuKaTopy.
func (s *Service) ByID(ctx context.Context, id int64) (*Customer, error) {
	item := &Customer{}

	err := s.pool.QueryRow(ctx, `
		SELECT id, name, phone, active, created FROM customers WHERE id = $1
	`, id).Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)
	
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}

	return item, nil
}

// All for
func (s *Service) All(ctx context.Context) ([]*Customer, error) {
	
	items := make([]*Customer,0)

	rows, err := s.pool.Query(ctx, `
		SELECT id, name, phone, active, created FROM customers
	`)
	
	if errors.Is(err, pgx.ErrNoRows) {
		log.Print("No rows")
		return nil, ErrNoRows
	}

	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}
	
	defer rows.Close()	
	
	for rows.Next() {
		item := &Customer{}
		err = rows.Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)
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
	//s.mu.RLock()
	//defer s.mu.RUnlock()
	// for _, Customer := range s.items {
	// 	if Customer.ID == id {
	// 		return Customer, nil
	// 	}
	// }
	// Customers := s.items
	// if len(s.items) == 0 {
	// 	return nil, errors.New("no items found")
	// }
	return items, nil
	//panic("not implemented")
}

//AllActive for
func (s *Service) AllActive(ctx context.Context) ([]*Customer, error) {
	
	items := make([]*Customer,0)

	rows, err := s.pool.Query(ctx, `
		SELECT id, name, phone, active, created FROM customers WHERE active
	`)
	
	if errors.Is(err, pgx.ErrNoRows) {
		log.Print("No rows")
		return nil, ErrNoRows
	}

	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}
	
	defer rows.Close()
	
	for rows.Next() {
		item := &Customer{}
		err = rows.Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)
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
	//s.mu.RLock()
	//defer s.mu.RUnlock()
	// for _, Customer := range s.items {
	// 	if Customer.ID == id {
	// 		return Customer, nil
	// 	}
	// }
	// Customers := s.items
	// if len(s.items) == 0 {
	// 	return nil, errors.New("no items found")
	// }
	return items, nil
	//panic("not implemented")
}

// BlockByID for
func (s *Service) BlockByID(ctx context.Context, id int64) (*Customer, error) {
	
	item := &Customer{}
	active := false

//	items := make([]*Customer,0)

	err := s.pool.QueryRow(ctx, `
		UPDATE customers SET active = $2 WHERE id = $1 RETURNING id, name, phone, active, created
	`, id, active).Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)
		
	if errors.Is(err, pgx.ErrNoRows) {
		log.Print("No rows")
		return nil, ErrNoRows
	}

	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}
	
	return item, nil
	//log.Print(result)	
	
	//s.mu.RLock()
	//defer s.mu.RUnlock()
	// for _, Customer := range s.items {
	// 	if Customer.ID == id {
	// 		return Customer, nil
	// 	}
	// }
	// Customers := s.items
	// if len(s.items) == 0 {
	// 	return nil, errors.New("no items found")
	// }
	//return items, nil
	//panic("not implemented")
}

// UnBlockByID for
func (s *Service) UnBlockByID(ctx context.Context, id int64) (*Customer, error) {
	
	item := &Customer{}
	active := true

//	items := make([]*Customer,0)

	err := s.pool.QueryRow(ctx, `
		UPDATE customers SET active = $2 WHERE id = $1 RETURNING id, name, phone, active, created
	`, id, active).Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)
		
	if errors.Is(err, pgx.ErrNoRows) {
		log.Print("No rows")
		return nil, ErrNoRows
	}

	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}
	
	return item, nil
	
}

// RemoveByID for
func (s *Service) RemoveByID(ctx context.Context, id int64) (*Customer, error) {
	
	item := &Customer{}
	

//	items := make([]*Customer,0)

	err := s.pool.QueryRow(ctx, `
		DELETE FROM customers WHERE id = $1 RETURNING id, name, phone, active, created
	`, id).Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)
		
	if errors.Is(err, pgx.ErrNoRows) {
		log.Print("No rows")
		return nil, ErrNoRows
	}

	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}
	
	return item, nil
	
}


// Save for
func (s *Service) Save(ctx context.Context, itemCustomer Customer) (*Customer, error) {
	
	item := &Customer{}

	if itemCustomer.ID == 0 {
	
		err := s.pool.QueryRow(ctx, `
		INSERT INTO customers(name, phone) VALUES($1, $2) RETURNING id, name, phone, active, created
		`, itemCustomer.Name, itemCustomer.Phone).Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)
		
		if errors.Is(err, pgx.ErrNoRows) {
			log.Print("No rows")
			return nil, ErrNoRows
		}

		if err != nil {
			log.Print(err)
			return nil, ErrInternal
		}
	
	}
	
	if itemCustomer.ID != 0 {
		
		err := s.pool.QueryRow(ctx, `
		UPDATE customers SET name = $2, phone = $3 WHERE id = $1 RETURNING id, name, phone, active, created
		`, itemCustomer.ID, itemCustomer.Name, itemCustomer.Phone).Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)
		
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

// Register for
func (s *Service) Register(ctx context.Context, registration *Registration) (*Customer, error) {
	
	item := &Customer{}

	hash, err := bcrypt.GenerateFromPassword([]byte(registration.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}

	err = s.pool.QueryRow(ctx, `
		INSERT INTO customers(name, phone, password)
		VALUES($1, $2, $3)
		ON CONFLICT (phone) DO NOTHING RETURNING id, name, phone, active, created
		`, registration.Name, registration.Phone, hash).Scan(
		&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)
		
	if errors.Is(err, pgx.ErrNoRows) {
		log.Print("No rows")
		return nil, ErrNoRows
	}

	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}

	return item, nil
		
}

//Auth for
type Auth struct {
	Login string `json:"login"`
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
	//item := &Token{}
	err = s.pool.QueryRow(ctx, `SELECT id, password FROM customers WHERE phone = $1`, phone).Scan(&id, &hash)

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
		INSERT INTO customers_tokens(token, customer_id) VALUES($1, $2)`,token, id)

	if err != nil {
		log.Print(err)
		return "", ErrInternal
	}


	return token, nil
}

//Product for
type Product struct {
	ID		int64	`json:"id"`
	Name	string	`json:"name"`
	Price	int		`json:"price"`
	Qty		int		`json:"qty"`
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

//IDByToken for
func (s *Service) IDByToken(ctx context.Context, token string) (int64, error) {
	var id int64
	err := s.pool.QueryRow(ctx, `
	 SELECT customer_id FROM customers_tokens WHERE token = $1
	 `, token).Scan(&id)

	 if err == pgx.ErrNoRows {
		return 0, nil
	}

	if err != nil {
		return 0, ErrInternal
	}
	log.Print("ID of customer ", id)
	return id, nil

}

//Purchase for
type Purchase struct {
	ID		int64	`json:"id"`
	ManagerID	int	`json:"manager_id"`
	CustomerID	int		`json:"customer_id"`
	
}

//Purchases for
func (s *Service) Purchases(ctx context.Context, id int64) ([]*Purchase, error) {
	items := make([]*Purchase, 0)
	rows, err := s.pool.Query(ctx, `
		SELECT id, manager_id, customer_id FROM sales WHERE customer_id = $1 ORDER BY id LIMIT 500 
	`,id)
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
		err = rows.Scan(&item.ID, &item.ManagerID, &item.CustomerID)
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


// //MakePurchases for
// func (s *Service) MakePurchases(ctx context.Context, id int64) ([]*Purchase, error) {
	
// 	err := s.pool.QueryRow(ctx, `
// 		INSERT INTO customers(name, phone) VALUES($1, $2) RETURNING id, name, phone, active, created
// 		`, itemCustomer.Name, itemCustomer.Phone).Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)
		
// 		if errors.Is(err, pgx.ErrNoRows) {
// 			log.Print("No rows")
// 			return nil, ErrNoRows
// 		}

// 		if err != nil {
// 			log.Print(err)
// 			return nil, ErrInternal
// 		}
	
	
// 	items := make([]*Purchase, 0)
// 	rows, err := s.pool.Query(ctx, `
// 		SELECT id, manager_id, customer_id FROM sales WHERE customer_id = $1 ORDER BY id LIMIT 500 
// 	`,id)
// 	if errors.Is(err, pgx.ErrNoRows) {
// 		return items, nil
// 	}
// 	if err != nil {
// 		log.Print(err)
// 		return nil, ErrInternal
// 	}
// 	defer rows.Close()
// 	for rows.Next() {
// 		item := &Purchase{}
// 		err = rows.Scan(&item.ID, &item.ManagerID, &item.CustomerID)
// 		if err != nil {
// 			log.Print(err)
// 			return nil, err
// 		}
// 		items = append(items, item)
// 	}
// 	err = rows.Err()
// 	if err != nil {
// 		log.Print(err)
// 		return nil, err
// 	}

// 	return items, nil
// }