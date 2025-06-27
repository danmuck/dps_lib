package finance

import (
	"sync"

	"github.com/google/uuid"
)

// get - Transaction and expenses by year, month, day
// post - add income or expense for a specific date
// put - update income or expense for a specific date
// delete - remove income or expense for a specific date
// props ->
// - year
// - month
// - day
// - type (income/expense)
// - amount
// - description
// - category (e.g., food, transport, etc.)
const (
	IncomeTxn  = "income"  // Type for income transactions
	ExpenseTxn = "expense" // Type for expense transactions
)

type Transaction struct {
	Type     string  `bson:"type" json:"type"`         // Type of transaction (income or expense)
	Amount   float64 `bson:"amount" json:"amount"`     // Amount of the transaction
	Category string  `bson:"category" json:"category"` // Category of the transaction (e.g., food, transport)
	Source   string  `bson:"source" json:"source"`     // Source of the transaction (e.g., bank account, cash)
	Note     string  `bson:"note" json:"note"`         // Additional note for the transaction
	Date     string  `bson:"date" json:"date"`         // Date of the transaction in YYYY-MM-DD format
}

func NewTransaction(category, source, note, date string, amount float64) Transaction {
	var tType string
	if amount >= 0 {
		tType = IncomeTxn // Positive amount indicates income
	} else {
		tType = ExpenseTxn // Negative amount indicates expense
		amount = -amount   // Store amount as positive for expenses
	}
	return Transaction{
		Type:     tType,
		Amount:   amount,
		Category: category,
		Source:   source,
		Note:     note,
		Date:     date,
	}
}

type Asset struct {
	ID          string  `json:"id"`          // Unique identifier for the asset
	Name        string  `json:"name"`        // Name of the asset
	Description string  `json:"description"` // Description of the asset
	Value       float64 `json:"value"`       // Current value of the asset
	Category    string  `json:"category"`    // Category of the asset (e.g., real estate, stocks)
	Date        string  `json:"date"`        // Date of the asset valuation in YYYY-MM-DD format
}

func NewAsset(name, description, category, date string, value float64) Asset {
	tmp := uuid.New().String()
	return Asset{
		ID:          tmp,
		Name:        name,
		Description: description,
		Value:       value,
		Category:    category,
		Date:        date,
	}
}

type Portfolio struct {
	ID          string
	Name        string
	Description string

	Assets   []string               // List of assets to contribute to total wealth
	Income   map[string]Transaction // List of income associated with the portfolio
	Expenses map[string]Transaction // List of expenses associated with the portfolio

	AssetTotal      float64 // total value of assets
	CashTotal       float64 // cash in the bank
	CreditLimit     float64 // credit limit
	CreditUsed      float64 // credit used
	LoansTotal      float64 // total loaned
	MonthlyIncome   float64 // total monthly income
	MonthlyExpenses float64 // total monthly expenses
	NetIncome       float64 // net income calculated as (MonthlyIncome - MonthlyExpenses)
	NetWorth        float64 // net worth calculated as (AssetTotal + CashTotal + CreditLimit - CreditUsed - LoansTotal)

	mu sync.Mutex // Mutex to protect concurrent access to the portfolio
}

func NewPortfolio(name, description string) Portfolio {
	tmp := uuid.New().String()
	return Portfolio{
		ID:          tmp,
		Name:        name,
		Description: description,
		Assets:      make([]string, 0),
		Income:      make(map[string]Transaction),
		Expenses:    make(map[string]Transaction),
		AssetTotal:  0.0,
		CashTotal:   0.0,
		CreditLimit: 0.0,
		CreditUsed:  0.0,
		LoansTotal:  0.0,

		MonthlyIncome:   0.0,
		MonthlyExpenses: 0.0,
	}
}

// AddAsset adds an asset to the portfolio and updates the total asset value
func (p *Portfolio) AddAsset(asset Asset) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Assets = append(p.Assets, asset.ID)
	p.AssetTotal += asset.Value
}

// AddIncome adds an income transaction to the portfolio and updates the monthly income
func (p *Portfolio) AddIncome(txn Transaction) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if txn.Type != IncomeTxn {
		return // Only add income transactions
	}
	p.Income[txn.Date] = txn
	p.MonthlyIncome += txn.Amount
	p.NetIncome = p.MonthlyIncome - p.MonthlyExpenses
}

// AddExpense adds an expense transaction to the portfolio and updates the monthly expenses
func (p *Portfolio) AddExpense(txn Transaction) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if txn.Type != ExpenseTxn {
		return // Only add expense transactions
	}
	p.Expenses[txn.Date] = txn
	p.MonthlyExpenses += txn.Amount
	p.NetIncome = p.MonthlyIncome - p.MonthlyExpenses
}
