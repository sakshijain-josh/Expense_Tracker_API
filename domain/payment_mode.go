package domain

// PaymentMode represents the mode of payment for an expense
type PaymentMode string

const (
	PaymentModeUPI  PaymentMode = "UPI"
	PaymentModeCash PaymentMode = "Cash"
)

// IsValid checks if the payment mode is valid
func (pm PaymentMode) IsValid() bool {
	return pm == PaymentModeUPI || pm == PaymentModeCash
}
