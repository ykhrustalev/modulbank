package modulbank

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

const DateFormat = "2006-01-02T15:04:05"

type AccountInfo struct {
	CompanyId    string        `json:"companyId"`
	CompanyName  string        `json:"companyName"`
	BankAccounts []BankAccount `json:"bankAccounts"`
}

type BankAccountCategory uint8

const (
	CheckingAccount       BankAccountCategory = iota
	DepositAccount
	TransitAccount
	CardAccount
	DepositRateAccount
	ReservationAccounting
)

func (category BankAccountCategory) String() string {
	switch category {
	case CheckingAccount:
		return "CheckingAccount"
	case DepositAccount:
		return "DepositAccount"
	case TransitAccount:
		return "TransitAccount"
	case CardAccount:
		return "CardAccount"
	case DepositRateAccount:
		return "DepositRateAccount"
	case ReservationAccounting:
		return "ReservationAccounting"
	default:
		return "unknown"
	}
}

func ParseBankAccountCategory(category string) (BankAccountCategory, error) {
	switch strings.ToLower(category) {
	case "checkingaccount":
		return CheckingAccount, nil
	case "depositaccount":
		return DepositAccount, nil
	case "transitaccount":
		return TransitAccount, nil
	case "cardaccount":
		return CardAccount, nil
	case "depositrateaccount":
		return DepositRateAccount, nil
	case "reservationaccounting":
		return ReservationAccounting, nil
	}

	var c BankAccountCategory
	return c, fmt.Errorf("invalid bank account category: %q", category)
}

type Currency uint8

const (
	CurrencyRUR Currency = iota
	CurrencyUSD
	CurrencyEUR
	CurrencyCNY
)

func (currency Currency) String() string {
	switch currency {
	case CurrencyRUR:
		return "RUR"
	case CurrencyUSD:
		return "USD"
	case CurrencyEUR:
		return "EUR"
	case CurrencyCNY:
		return "CNY"
	default:
		return "unknown"
	}
}

func ParseCurrency(currency string) (Currency, error) {
	switch strings.ToLower(currency) {
	case "rur":
		return CurrencyRUR, nil
	case "usd":
		return CurrencyUSD, nil
	case "eur":
		return CurrencyEUR, nil
	case "cny":
		return CurrencyCNY, nil
	}

	var c Currency
	return c, fmt.Errorf("invalid currency: %q", currency)
}

type BankAccountStatus uint8

const (
	NewAccount      BankAccountStatus = iota
	DeletedAccount
	ClosedAccount
	FreezedAccount
	ToClosedAccount
	ToOpenAccount
)

func (status BankAccountStatus) String() string {
	switch status {
	case NewAccount:
		return "New"
	case DeletedAccount:
		return "Deleted"
	case ClosedAccount:
		return "Closed"
	case FreezedAccount:
		return "Freezed"
	case ToClosedAccount:
		return "ToClosed"
	case ToOpenAccount:
		return "ToOpen"
	default:
		return "unknown"
	}
}

func ParseBankAccountStatus(status string) (BankAccountStatus, error) {
	switch strings.ToLower(status) {
	case "new":
		return NewAccount, nil
	case "deleted":
		return DeletedAccount, nil
	case "closed":
		return ClosedAccount, nil
	case "freezed":
		return FreezedAccount, nil
	case "toclosed":
		return ToClosedAccount, nil
	case "toopen":
		return ToOpenAccount, nil
	}

	var s BankAccountStatus
	return s, fmt.Errorf("invalid bank account status: %q", status)
}

type auxBankAccount struct {
	Id                       string  `json:"id"`
	AccountName              string  `json:"accountName"`
	Balance                  float32 `json:"balance"`
	BankBic                  string  `json:"bankBic"`
	BankInn                  string  `json:"bankInn"`
	BankKpp                  string  `json:"bankKpp"`
	BankCorrespondentAccount string  `json:"bankCorrespondentAccount"`
	BankName                 string  `json:"bankName"`
	BeginDate                string  `json:"beginDate"`
	Category                 string  `json:"category"`
	Currency                 string  `json:"currency"`
	Number                   string  `json:"number"`
	Status                   string  `json:"status"`
}

type BankAccount struct {
	Id                       string
	AccountName              string
	Balance                  float32
	BankBic                  string
	BankInn                  string
	BankKpp                  string
	BankCorrespondentAccount string
	BankName                 string
	BeginDate                time.Time
	Category                 BankAccountCategory
	Currency                 Currency
	Number                   string
	Status                   BankAccountStatus
}

func (acc BankAccount) UnmarshalJSON(b []byte) error {
	var aux auxBankAccount

	err := json.Unmarshal(b, &aux)
	if err != nil {
		return err
	}

	acc.Id = aux.Id
	acc.AccountName = aux.AccountName
	acc.Balance = aux.Balance
	acc.BankBic = aux.BankBic
	acc.BankInn = aux.BankInn
	acc.BankKpp = aux.BankKpp
	acc.BankCorrespondentAccount = aux.BankCorrespondentAccount
	acc.BankName = aux.BankName
	acc.BeginDate, err = time.Parse(DateFormat, aux.BeginDate)
	if err != nil {
		return err
	}
	acc.Category, err = ParseBankAccountCategory(aux.Category)
	if err != nil {
		return err
	}
	acc.Currency, err = ParseCurrency(aux.Currency)
	if err != nil {
		return err
	}
	acc.Number = aux.Number
	acc.Status, err = ParseBankAccountStatus(aux.Status)
	if err != nil {
		return err
	}

	return nil
}

func (acc *BankAccount) MarshalJSON() ([]byte, error) {
	return json.Marshal(&auxBankAccount{
		Id:                       acc.Id,
		AccountName:              acc.AccountName,
		Balance:                  acc.Balance,
		BankBic:                  acc.BankBic,
		BankInn:                  acc.BankInn,
		BankKpp:                  acc.BankKpp,
		BankCorrespondentAccount: acc.BankCorrespondentAccount,
		BankName:                 acc.BankName,
		BeginDate:                acc.BeginDate.Format(DateFormat),
		Category:                 acc.Category.String(),
		Currency:                 acc.Currency.String(),
		Number:                   acc.Number,
		Status:                   acc.Status.String(),
	})
}
