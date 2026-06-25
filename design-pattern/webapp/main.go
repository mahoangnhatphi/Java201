package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"design-pattern/behavioral"
	"design-pattern/creational"
)

func main() {
	server := NewWebServer()
	fmt.Println("Starting server on :8080 ...")
	fmt.Println("Try these curl commands:")
	fmt.Println("  curl http://localhost:8080/")
	fmt.Println("  curl -X POST http://localhost:8080/api/users -d 'name=John&email=john@example.com'")
	fmt.Println("  curl -X POST http://localhost:8080/api/payments -d 'method=credit_card&amount=99.99&card_number=1234'")
	fmt.Println("  curl http://localhost:8080/api/users")
	fmt.Println()
	if err := server.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Server failed: %v\n", err)
		os.Exit(1)
	}
}

type WebServer struct {
	port           string
	eventBus       *behavioral.EventBus
	paymentFactory *creational.PaymentFactory
	userService    *behavioral.UserService
	config         *creational.Config
}

func NewWebServer() *WebServer {
	paymentFactory := creational.NewPaymentFactory()
	config := creational.GetConfigManager().GetConfig()

	eventBus := behavioral.NewEventBus()
	eventBus.Register(behavioral.NewEmailNotifier("SendGrid"))
	eventBus.Register(behavioral.NewLoggerObserver())
	eventBus.Register(behavioral.NewAnalyticsObserver("GA-123456"))

	postgres := behavioral.NewPostgresConnection("postgres://localhost:5432/webstore")
	userRepo := behavioral.NewUserRepository(postgres)
	emailSvc := behavioral.NewSendGridEmailService("SG.api-key")
	consoleLogger := behavioral.NewConsoleLogger()
	userService := behavioral.NewUserService(userRepo, emailSvc, consoleLogger)

	return &WebServer{
		port:           ":8080",
		eventBus:       eventBus,
		paymentFactory: paymentFactory,
		userService:    userService,
		config:         config,
	}
}

func (s *WebServer) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleHome)
	mux.HandleFunc("/api/users", s.handleGetUsers)
	mux.HandleFunc("/api/payments", s.handleCreatePayment)
	mux.HandleFunc("/api/users/create", s.handleCreateUser)

	fmt.Printf("Server running at http://localhost%s\n", s.port)
	return http.ListenAndServe(s.port, mux)
}

func (s *WebServer) handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"app":     s.config.AppName,
		"version": s.config.Version,
		"status":  "running",
		"endpoints": []string{
			"GET  /",
			"POST /api/users/create",
			"POST /api/payments",
			"GET  /api/users",
		},
	})
}

func (s *WebServer) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	email := r.FormValue("email")

	if err := s.userService.RegisterUser(name, email); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.eventBus.Notify(behavioral.Event{
		Type:    behavioral.UserCreated,
		Payload: map[string]string{"name": name, "email": email},
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User created successfully",
		"name":    name,
		"email":   email,
	})
}

func (s *WebServer) handleCreatePayment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	methodType := r.FormValue("method")
	amount := r.FormValue("amount")

	details := map[string]string{
		"method": methodType,
	}

	switch methodType {
	case "credit_card":
		details["card_number"] = r.FormValue("card_number")
	case "paypal":
		details["email"] = r.FormValue("email")
	case "crypto":
		details["wallet_address"] = r.FormValue("wallet_address")
	}

	payment, err := s.paymentFactory.CreatePaymentMethod(methodType, details)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var amountFloat float64
	fmt.Sscanf(amount, "%f", &amountFloat)

	result := payment.Pay(amountFloat)

	s.eventBus.Notify(behavioral.Event{
		Type:    behavioral.OrderPlaced,
		Payload: map[string]interface{}{"method": methodType, "amount": amountFloat},
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"result": result,
	})
}

func (s *WebServer) handleGetUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	users, err := s.userService.GetAllUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
