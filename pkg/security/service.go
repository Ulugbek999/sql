package security

import (
	"context"
	"encoding/hex"
	"errors"
	"log"
	"math/rand"
	"time"


)

// Service for
type Service struct {
	pool *pgxpool.Pool
}

// NewService co3qaét cepsuc.
func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
	
}

// Auth for
func (s *Service) Auth(login, password string) (ok bool) {
	ctx := context.Background()

	//item := &Customer{}

	rows, err := s.pool.Query(ctx, `
		SELECT id, name FROM managers WHERE login = $1 AND password = $2
	`, login, password)
	
	if errors.Is(err, pgx.ErrNoRows) {
		log.Print("No rows")
		return false
	}

	if err != nil {
		log.Print(err)
		return false
	}

	if !rows.Next() {
		return false
	}

	defer rows.Close()
	
	return true
}

// Token for
type Token struct {
	Token 	string		`json:"token"`
	ID 		int64	 	`json:"id"`
	Expire	time.Time 	`json:"expire"`
	Created time.Time 	`json:"created"`
}

//TokenStatus for
type TokenStatus struct {
	Status string `json:"status"`
	CustomerID int64 `json:"customerId"`
}
//TokenReason for
type TokenReason struct {
	Status string `json:"status"`
	Reason string `json:"reason"`
}

//ErrNoSuchUser for
var ErrNoSuchUser = errors.New("no such user")
//ErrInvalidPassword for
var ErrInvalidPassword = errors.New("invalid password")
//ErrInternal for
var ErrInternal = errors.New("internal error")

//TokenForCustomer for
func (s *Service) TokenForCustomer(ctx context.Context, phone string, password string) (token string, err error) {
	
	var hash string
	var id int64
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
	_, err = s.pool.Exec(ctx, `INSERT INTO customers_tokens(token, customer_id) VALUES($1, $2)`, token, id)
	if err != nil {
		return "", ErrInternal
	}

	return token, nil
}

//AuthenticateCustomer for
func (s *Service) AuthenticateCustomer(ctx context.Context, token string) (id int64, err error) {
	err = s.pool.QueryRow(ctx, `SELECT customer_id FROM customers_tokens WHERE token = $1`, token).Scan(&id)
	if err == pgx.ErrNoRows {
		return 0, ErrNoSuchUser
	}

	if err != nil {
		return 0, ErrInternal
	}	

	return id, nil
}



//CheckTokenCustomer for
func (s *Service) CheckTokenCustomer(ctx context.Context, token string) (inter interface{}, status int64) {
	//var hash string
	var id int64
	var item1 TokenReason
	var item TokenStatus
	var expire time.Time

	log.Print("token", token)
	
	//log.Print(timeNow)
	err := s.pool.QueryRow(ctx, `SELECT customer_id, expire FROM customers_tokens WHERE token = $1`, token).Scan(&id, &expire)

	if err == pgx.ErrNoRows {
		item1.Status = "fail"
		item1.Reason = "not found"
		return item1, 404
	}

	if err != nil {
		log.Print("We are hereeeeeeeeeeeeeeeeeeeeeeeeeeeeee")
		item1.Status = "fail"
		item1.Reason = "not found"
		return item1, 404
	}

	if time.Now().After(expire) {
		log.Print("We are hereeeeeeeeeeeeeeeeeeeeeeeeeeeeee8")
		item1.Status = "fail"
		item1.Reason = "expired"
		return item1, 400
	}

	item.Status = "ok"
	item.CustomerID = id
		
	return item, 200
}