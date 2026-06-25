package behavioral

import "fmt"

// DatabaseConnection defines the interface for database operations
type DatabaseConnection interface {
	Query(query string) ([]map[string]interface{}, error)
	Execute(query string, args ...interface{}) error
}

// PostgresConnection implements DatabaseConnection
type PostgresConnection struct {
	connectionString string
}

func NewPostgresConnection(connectionString string) *PostgresConnection {
	return &PostgresConnection{connectionString: connectionString}
}

func (p *PostgresConnection) Query(query string) ([]map[string]interface{}, error) {
	fmt.Printf("[Postgres] Executing query: %s on %s\n", query, p.connectionString)
	return []map[string]interface{}{
		{"id": 1, "name": "John Doe"},
		{"id": 2, "name": "Jane Smith"},
	}, nil
}

func (p *PostgresConnection) Execute(query string, args ...interface{}) error {
	fmt.Printf("[Postgres] Executing: %s with args %v on %s\n", query, args, p.connectionString)
	return nil
}

// MySQLConnection implements DatabaseConnection
type MySQLConnection struct {
	connectionString string
}

func NewMySQLConnection(connectionString string) *MySQLConnection {
	return &MySQLConnection{connectionString: connectionString}
}

func (m *MySQLConnection) Query(query string) ([]map[string]interface{}, error) {
	fmt.Printf("[MySQL] Executing query: %s on %s\n", query, m.connectionString)
	return []map[string]interface{}{
		{"id": 1, "name": "Alice Johnson"},
		{"id": 2, "name": "Bob Williams"},
	}, nil
}

func (m *MySQLConnection) Execute(query string, args ...interface{}) error {
	fmt.Printf("[MySQL] Executing: %s with args %v on %s\n", query, args, m.connectionString)
	return nil
}

// UserRepository handles user data operations
type UserRepository struct {
	db DatabaseConnection
}

func NewUserRepository(db DatabaseConnection) *UserRepository {
	return &UserRepository{db: db}
}

func (u *UserRepository) FindAll() ([]map[string]interface{}, error) {
	return u.db.Query("SELECT * FROM users")
}

func (u *UserRepository) Create(name, email string) error {
	return u.db.Execute("INSERT INTO users (name, email) VALUES (?, ?)", name, email)
}

// EmailService defines the interface for sending emails
type EmailService interface {
	Send(to, subject, body string) error
}

// SendGridEmailService implements EmailService
type SendGridEmailService struct {
	apiKey string
}

func NewSendGridEmailService(apiKey string) *SendGridEmailService {
	return &SendGridEmailService{apiKey: apiKey}
}

func (s *SendGridEmailService) Send(to, subject, body string) error {
	fmt.Printf("[SendGrid] Sending email to %s with subject '%s'\n", to, subject)
	return nil
}

// SMTPEmailService implements EmailService
type SMTPEmailService struct {
	smtpHost string
	smtpPort int
}

func NewSMTPEmailService(smtpHost string, smtpPort int) *SMTPEmailService {
	return &SMTPEmailService{smtpHost: smtpHost, smtpPort: smtpPort}
}

func (s *SMTPEmailService) Send(to, subject, body string) error {
	fmt.Printf("[SMTP] Sending email to %s via %s:%d with subject '%s'\n", to, s.smtpHost, s.smtpPort, subject)
	return nil
}

// UserService handles user business logic
type UserService struct {
	userRepo *UserRepository
	emailSvc EmailService
	logger   Logger
}

func NewUserService(userRepo *UserRepository, emailSvc EmailService, logger Logger) *UserService {
	return &UserService{
		userRepo: userRepo,
		emailSvc: emailSvc,
		logger:   logger,
	}
}

func (u *UserService) RegisterUser(name, email string) error {
	u.logger.Log("Registering new user: " + name)

	if err := u.userRepo.Create(name, email); err != nil {
		u.logger.Log("Failed to create user: " + err.Error())
		return err
	}

	if err := u.emailSvc.Send(email, "Welcome!", "Welcome to our platform!"); err != nil {
		u.logger.Log("Failed to send welcome email: " + err.Error())
		return err
	}

	u.logger.Log("User registered successfully: " + name)
	return nil
}

func (u *UserService) GetAllUsers() ([]map[string]interface{}, error) {
	u.logger.Log("Fetching all users")
	return u.userRepo.FindAll()
}

// Logger defines the interface for logging
type Logger interface {
	Log(message string)
}

// ConsoleLogger implements Logger
type ConsoleLogger struct{}

func NewConsoleLogger() *ConsoleLogger {
	return &ConsoleLogger{}
}

func (c *ConsoleLogger) Log(message string) {
	fmt.Printf("[ConsoleLogger] %s\n", message)
}

// FileLogger implements Logger
type FileLogger struct {
	filePath string
}

func NewFileLogger(filePath string) *FileLogger {
	return &FileLogger{filePath: filePath}
}

func (f *FileLogger) Log(message string) {
	fmt.Printf("[FileLogger] Writing to %s: %s\n", f.filePath, message)
}

// DependencyInjectionExampleUsage demonstrates Dependency Injection
func DependencyInjectionExampleUsage() {
	postgres := NewPostgresConnection("postgres://localhost:5432/myapp")
	userRepo := NewUserRepository(postgres)

	emailSvc := NewSendGridEmailService("SG.xxxxx")
	consoleLogger := NewConsoleLogger()

	userService := NewUserService(userRepo, emailSvc, consoleLogger)

	userService.RegisterUser("John Doe", "john@example.com")
	users, _ := userService.GetAllUsers()
	fmt.Printf("Users: %v\n", users)
}
