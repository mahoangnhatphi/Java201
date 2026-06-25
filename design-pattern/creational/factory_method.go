package creational

import "fmt"

// PaymentMethod defines the interface for different payment methods
type PaymentMethod interface {
	Pay(amount float64) string
}

// CreditCardPayment implements PaymentMethod
type CreditCardPayment struct {
	cardNumber string
}

func NewCreditCardPayment(cardNumber string) *CreditCardPayment {
	return &CreditCardPayment{cardNumber: cardNumber}
}

func (c *CreditCardPayment) Pay(amount float64) string {
	return fmt.Sprintf("Paid %.2f using Credit Card ending in %s", amount, c.cardNumber[len(c.cardNumber)-4:])
}

// PayPalPayment implements PaymentMethod
type PayPalPayment struct {
	email string
}

func NewPayPalPayment(email string) *PayPalPayment {
	return &PayPalPayment{email: email}
}

func (p *PayPalPayment) Pay(amount float64) string {
	return fmt.Sprintf("Paid %.2f using PayPal (%s)", amount, p.email)
}

// CryptoPayment implements PaymentMethod
type CryptoPayment struct {
	walletAddress string
}

func NewCryptoPayment(walletAddress string) *CryptoPayment {
	return &CryptoPayment{walletAddress: walletAddress}
}

func (c *CryptoPayment) Pay(amount float64) string {
	return fmt.Sprintf("Paid %.2f using Crypto wallet %s", amount, c.walletAddress[:8]+"...")
}

// PaymentFactory is the factory that creates payment methods
type PaymentFactory struct{}

func NewPaymentFactory() *PaymentFactory {
	return &PaymentFactory{}
}

func (f *PaymentFactory) CreatePaymentMethod(methodType string, details map[string]string) (PaymentMethod, error) {
	switch methodType {
	case "credit_card":
		cardNumber, ok := details["card_number"]
		if !ok {
			return nil, fmt.Errorf("card_number is required for credit_card")
		}
		return NewCreditCardPayment(cardNumber), nil
	case "paypal":
		email, ok := details["email"]
		if !ok {
			return nil, fmt.Errorf("email is required for paypal")
		}
		return NewPayPalPayment(email), nil
	case "crypto":
		wallet, ok := details["wallet_address"]
		if !ok {
			return nil, fmt.Errorf("wallet_address is required for crypto")
		}
		return NewCryptoPayment(wallet), nil
	default:
		return nil, fmt.Errorf("unsupported payment method: %s", methodType)
	}
}

// FactoryExampleUsage demonstrates the Factory Method pattern
func FactoryExampleUsage() {
	factory := NewPaymentFactory()

	payment, err := factory.CreatePaymentMethod("credit_card", map[string]string{
		"card_number": "1234-5678-9012-3456",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(payment.Pay(100.50))

	payment, err = factory.CreatePaymentMethod("paypal", map[string]string{
		"email": "user@example.com",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(payment.Pay(250.00))

	payment, err = factory.CreatePaymentMethod("crypto", map[string]string{
		"wallet_address": "0x1234567890abcdef",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(payment.Pay(500.75))
}
