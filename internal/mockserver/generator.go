// Package mockserver provides a mock YNAB API server for demo/testing.
//
//nolint:gosec // G404: math/rand is acceptable for mock data generation
package mockserver

import (
	"math/rand"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// MockConfig defines the configuration for generating mock data.
type MockConfig struct {
	BudgetCount                    int // Number of budgets (6-10)
	AccountsPerBudget              int // Accounts per budget (3-5)
	CategoryGroups                 int // Category groups per budget (4-6)
	CategoriesPerGroup             int // Categories per group (3-5)
	MonthsOfHistory                int // Months of transaction history (6-36)
	TransactionsPerMonthPerAccount int // Transactions per month per account (5+)
}

// DefaultMockConfig returns the default mock configuration.
func DefaultMockConfig() MockConfig {
	return MockConfig{
		BudgetCount:                    8,
		AccountsPerBudget:              4,
		CategoryGroups:                 5,
		CategoriesPerGroup:             4,
		MonthsOfHistory:                18,
		TransactionsPerMonthPerAccount: 8,
	}
}

// Generator creates mock YNAB data.
type Generator struct {
	details map[string]*BudgetDetail
	budgets []BudgetSummary
	config  MockConfig
}

// NewGenerator creates a new mock data generator.
func NewGenerator(config MockConfig) *Generator {
	return &Generator{
		config:  config,
		details: make(map[string]*BudgetDetail),
	}
}

// Budget names for realistic demo data.
var budgetNames = []string{
	"Personal Budget",
	"Family Finances",
	"Household Expenses",
	"Monthly Budget",
	"Joint Expenses",
	"Side Hustle",
	"Travel Planning",
	"Home Renovation",
	"Wedding Budget",
	"New Baby Fund",
}

// Account names and types.
var accountTemplates = []struct {
	Name     string
	Type     AccountType
	OnBudget bool
}{
	{"Checking Account", Checking, true},
	{"Savings Account", Savings, true},
	{"Credit Card", CreditCard, true},
	{"Cash", Cash, true},
	{"Investment Account", OtherAsset, false},
	{"Car Loan", AutoLoan, false},
	{"Mortgage", Mortgage, false},
	{"Student Loans", StudentLoan, false},
}

// Category group and category templates.
var categoryGroupTemplates = []struct {
	Name       string
	Categories []string
}{
	{"Fixed Expenses", []string{"Rent/Mortgage", "Utilities", "Insurance", "Phone", "Internet"}},
	{"Flexible Spending", []string{"Groceries", "Dining Out", "Entertainment", "Shopping", "Personal Care"}},
	{"Transportation", []string{"Gas", "Car Maintenance", "Public Transit", "Parking", "Rideshare"}},
	{"Health & Wellness", []string{"Medical", "Gym Membership", "Pharmacy", "Mental Health"}},
	{"Savings Goals", []string{"Emergency Fund", "Vacation", "New Car", "Home Down Payment"}},
	{"Debt Payments", []string{"Credit Card Payment", "Student Loan", "Car Payment"}},
	{"Subscriptions", []string{"Streaming Services", "Software", "News & Magazines", "Cloud Storage"}},
}

// Payee names for transactions.
var payeeNames = []string{
	"Whole Foods Market",
	"Amazon.com",
	"Netflix",
	"Spotify",
	"Target",
	"Costco",
	"Trader Joe's",
	"Shell Gas Station",
	"Chevron",
	"Starbucks",
	"Uber",
	"Lyft",
	"AT&T",
	"Verizon",
	"Electric Company",
	"Water & Sewer",
	"City Parking",
	"Planet Fitness",
	"CVS Pharmacy",
	"Walgreens",
	"Home Depot",
	"Lowe's",
	"Apple Store",
	"Best Buy",
	"Restaurant XYZ",
	"Coffee Shop",
	"Local Grocery",
	"Gas Station",
	"Insurance Co",
	"Medical Center",
}

// GenerateBudgets generates mock budget summaries.
func (g *Generator) GenerateBudgets() []BudgetSummary {
	if g.budgets != nil {
		return g.budgets
	}

	g.budgets = make([]BudgetSummary, g.config.BudgetCount)
	now := time.Now()
	cumulativeMonths := 0

	for i := 0; i < g.config.BudgetCount; i++ {
		budgetID := uuid.New()
		// Stagger last modified dates by 2-8 months each, plus random day jitter
		var jitterDays int
		if i == 0 {
			jitterDays = rand.Intn(7) // 0-6 days for latest budget (within a week)
		} else {
			cumulativeMonths += rand.Intn(7) + 2 // 2-8 months
			jitterDays = rand.Intn(28) + 3       // 3-30 days
		}
		lastModified := now.AddDate(0, -cumulativeMonths, -jitterDays)

		// Calculate first and last month
		firstMonth := now.AddDate(0, -g.config.MonthsOfHistory, 0)
		firstMonthDate := openapi_types.Date{Time: time.Date(firstMonth.Year(), firstMonth.Month(), 1, 0, 0, 0, 0, time.UTC)}
		lastMonthDate := openapi_types.Date{Time: time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)}

		g.budgets[i] = BudgetSummary{
			Id:             budgetID,
			Name:           budgetNames[i%len(budgetNames)],
			LastModifiedOn: &lastModified,
			FirstMonth:     &firstMonthDate,
			LastMonth:      &lastMonthDate,
			CurrencyFormat: &CurrencyFormat{
				IsoCode:          "USD",
				CurrencySymbol:   "$",
				DecimalDigits:    2,
				DecimalSeparator: ".",
				DisplaySymbol:    true,
				GroupSeparator:   ",",
				SymbolFirst:      true,
			},
			DateFormat: &DateFormat{
				Format: "MM/DD/YYYY",
			},
		}
	}

	return g.budgets
}

// GenerateBudgetDetail generates a full budget detail for a given budget ID.
func (g *Generator) GenerateBudgetDetail(budgetID string) *BudgetDetail {
	// Check if we've already generated this budget
	if detail, ok := g.details[budgetID]; ok {
		return detail
	}

	// Find the budget summary
	var summary *BudgetSummary
	for i := range g.budgets {
		if g.budgets[i].Id.String() == budgetID {
			summary = &g.budgets[i]
			break
		}
	}

	if summary == nil {
		return nil
	}

	// Generate accounts
	accounts := g.generateAccounts()

	// Generate category groups and categories
	categoryGroups, categories := g.generateCategories()

	// Generate payees
	payees := g.generatePayees(accounts)

	// Generate transactions
	transactions := g.generateTransactions(accounts, categories, payees)

	detail := &BudgetDetail{
		Id:             summary.Id,
		Name:           summary.Name,
		LastModifiedOn: summary.LastModifiedOn,
		FirstMonth:     summary.FirstMonth,
		LastMonth:      summary.LastMonth,
		CurrencyFormat: summary.CurrencyFormat,
		DateFormat:     summary.DateFormat,
		Accounts:       &accounts,
		CategoryGroups: &categoryGroups,
		Categories:     &categories,
		Payees:         &payees,
		Transactions:   &transactions,
	}

	g.details[budgetID] = detail
	return detail
}

func (g *Generator) generateAccounts() []Account {
	count := g.config.AccountsPerBudget
	accounts := make([]Account, count)

	for i := 0; i < count; i++ {
		template := accountTemplates[i%len(accountTemplates)]
		accountID := uuid.New()
		transferPayeeID := uuid.New()

		// Random balance between -5000 and 50000 dollars (in milliunits)
		balance := int64((rand.Intn(55000) - 5000) * 1000)
		clearedBalance := balance - int64(rand.Intn(500)*1000)

		var note *string
		if rand.Float32() < 0.3 {
			n := faker.Sentence()
			note = &n
		}

		accounts[i] = Account{
			Id:               accountID,
			Name:             template.Name,
			Type:             template.Type,
			OnBudget:         template.OnBudget,
			Closed:           i == count-1 && rand.Float32() < 0.3, // Maybe close the last one
			Deleted:          false,
			Balance:          balance,
			ClearedBalance:   clearedBalance,
			UnclearedBalance: balance - clearedBalance,
			TransferPayeeId:  &transferPayeeID,
			Note:             note,
		}
	}

	return accounts
}

func (g *Generator) generateCategories() ([]CategoryGroup, []Category) {
	var categoryGroups []CategoryGroup
	var categories []Category

	numGroups := g.config.CategoryGroups
	if numGroups > len(categoryGroupTemplates) {
		numGroups = len(categoryGroupTemplates)
	}

	for i := 0; i < numGroups; i++ {
		template := categoryGroupTemplates[i]
		groupID := uuid.New()

		categoryGroups = append(categoryGroups, CategoryGroup{
			Id:      groupID,
			Name:    template.Name,
			Hidden:  false,
			Deleted: false,
		})

		numCats := g.config.CategoriesPerGroup
		if numCats > len(template.Categories) {
			numCats = len(template.Categories)
		}

		for j := 0; j < numCats; j++ {
			catID := uuid.New()
			budgeted := int64(rand.Intn(100000) * 10) // 0-1000 dollars in milliunits
			activity := int64(-rand.Intn(80000) * 10) // Negative (spending)
			balance := budgeted + activity

			categories = append(categories, Category{
				Id:              catID,
				CategoryGroupId: groupID,
				Name:            template.Categories[j],
				Hidden:          rand.Float32() < 0.05, // 5% hidden
				Deleted:         rand.Float32() < 0.02, // 2% deleted
				Budgeted:        budgeted,
				Activity:        activity,
				Balance:         balance,
			})
		}
	}

	return categoryGroups, categories
}

func (g *Generator) generatePayees(accounts []Account) []Payee {
	payees := make([]Payee, 0, len(payeeNames)+len(accounts))

	// Regular payees
	for _, name := range payeeNames {
		payees = append(payees, Payee{
			Id:      uuid.New(),
			Name:    name,
			Deleted: false,
		})
	}

	// Transfer payees for each account
	for i := range accounts {
		transferAccountID := accounts[i].Id.String()
		payees = append(payees, Payee{
			Id:                uuid.New(),
			Name:              "Transfer : " + accounts[i].Name,
			TransferAccountId: &transferAccountID,
			Deleted:           false,
		})
	}

	return payees
}

//nolint:gocognit // Complex function generating realistic transaction data
func (g *Generator) generateTransactions(accounts []Account, categories []Category, payees []Payee) []TransactionSummary {
	var transactions []TransactionSummary
	now := time.Now()

	// Filter to only non-transfer payees for regular transactions
	regularPayees := make([]Payee, 0)
	for _, p := range payees {
		if p.TransferAccountId == nil {
			regularPayees = append(regularPayees, p)
		}
	}

	// Filter to only active categories
	activeCategories := make([]Category, 0)
	for i := range categories {
		if !categories[i].Deleted && !categories[i].Hidden {
			activeCategories = append(activeCategories, categories[i])
		}
	}

	for accIdx := range accounts {
		acc := &accounts[accIdx]
		if acc.Closed || acc.Deleted {
			continue
		}

		for month := 0; month < g.config.MonthsOfHistory; month++ {
			monthDate := now.AddDate(0, -month, 0)

			for t := 0; t < g.config.TransactionsPerMonthPerAccount; t++ {
				txnID := uuid.New().String()

				// Random day in the month
				day := rand.Intn(28) + 1
				txnDate := time.Date(monthDate.Year(), monthDate.Month(), day, 0, 0, 0, 0, time.UTC)

				// Random amount: mostly negative (expenses), occasionally positive (income)
				var amount int64
				if rand.Float32() < 0.15 { // 15% chance of income
					amount = int64(rand.Intn(300000) + 10000) // $10-$3010 income
				} else {
					amount = int64(-rand.Intn(50000) - 100) // $0.10-$500 expense
				}

				// Pick random payee and category
				payee := regularPayees[rand.Intn(len(regularPayees))]
				category := activeCategories[rand.Intn(len(activeCategories))]

				// Cleared status
				cleared := Cleared
				if month == 0 && rand.Float32() < 0.2 {
					cleared = Uncleared
				}

				payeeID := payee.Id
				categoryID := category.Id

				var memo *string
				if rand.Float32() < 0.2 {
					m := faker.Sentence()
					memo = &m
				}

				transactions = append(transactions, TransactionSummary{
					Id:         txnID,
					AccountId:  acc.Id,
					Date:       openapi_types.Date{Time: txnDate},
					Amount:     amount,
					PayeeId:    &payeeID,
					CategoryId: &categoryID,
					Cleared:    cleared,
					Approved:   true,
					Deleted:    false,
					Memo:       memo,
				})
			}
		}
	}

	return transactions
}

// GenerateUser generates a mock user.
func (g *Generator) GenerateUser() User {
	return User{
		Id: uuid.New(),
	}
}
