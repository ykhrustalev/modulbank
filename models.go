package modulbank

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

const DateFormat = "2006-01-02T15:04:05"

// API returns all dates in MSK timezone without offset
var MskLocation *time.Location

func init() {
	var err error
	MskLocation, err = time.LoadLocation("Europe/Moscow")
	if err != nil {
		log.Printf("Can't parse Moscow timezone, %v", err)
	}
}

type AccountInfo struct {
	CompanyId    string        `json:"companyId"`
	CompanyName  string        `json:"companyName"`
	BankAccounts []BankAccount `json:"bankAccounts"`
}

type BankAccountCategory uint8

const (
	CheckingAccount BankAccountCategory = iota
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
	return c, fmt.Errorf("invalid bank account category: %s", category)
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
	return c, fmt.Errorf("invalid currency: %s", currency)
}

type BankAccountStatus uint8

const (
	NewAccount BankAccountStatus = iota
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
	return s, fmt.Errorf("invalid bank account status: %s", status)
}

type auxBankAccount struct {
	Id                       string  `json:"id"`
	AccountName              string  `json:"accountName"`
	Balance                  float64 `json:"balance"`
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
	Balance                  float64
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

func (acc *BankAccount) UnmarshalJSON(b []byte) error {
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
	acc.BeginDate, err = time.ParseInLocation(DateFormat, aux.BeginDate, MskLocation)
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
		BeginDate:                acc.BeginDate.In(MskLocation).Format(DateFormat),
		Category:                 acc.Category.String(),
		Currency:                 acc.Currency.String(),
		Number:                   acc.Number,
		Status:                   acc.Status.String(),
	})
}

type OperationCategory uint8

const (
	OperationCategoryNone OperationCategory = iota
	OperationCategoryDebet
	OperationCategoryCredit
)

func (status OperationCategory) String() string {
	switch status {
	case OperationCategoryDebet:
		return "Debet"
	case OperationCategoryCredit:
		return "Credit"
	default:
		return "unknown"
	}
}

func ParseOperationCategory(status string) (OperationCategory, error) {
	switch strings.ToLower(status) {
	case "debet":
		return OperationCategoryDebet, nil
	case "credit":
		return OperationCategoryCredit, nil
	}

	var o OperationCategory
	return o, fmt.Errorf("invalid operation category: %s", status)
}

//
// Operations
//

type OperationHistorySearch struct {
	Category OperationCategory
	From     *time.Time
	Till     *time.Time
	Skip     int // offset
	Records  int // limit
}

type auxOperationHistorySearch struct {
	Category string `json:"category,omitifempty"`
	From     string `json:"from,omitifempty"`
	Till     string `json:"till,omitifempty"`
	Skip     int    `json:"skip,omitifempty"`
	Records  int    `json:"records,omitifempty"`
}

func (search *OperationHistorySearch) MarshalJSON() ([]byte, error) {
	aux := &auxOperationHistorySearch{}

	if search.Category != OperationCategoryNone {
		aux.Category = search.Category.String()
	}
	if search.From != nil {
		aux.Category = search.From.In(MskLocation).Format(DateFormat)
	}
	if search.Till != nil {
		aux.Category = search.Till.In(MskLocation).Format(DateFormat)
	}
	if search.Skip != 0 {
		aux.Skip = search.Skip
	}
	if search.Records != 0 {
		aux.Records = search.Records
	}

	return json.Marshal(aux)
}

type OperationStatus uint8

const (
	OperationStatusSendToBank OperationStatus = iota
	OperationStatusExecuted
	OperationStatusRejectByBank
	OperationStatusCanceled
	OperationStatusReceived
)

func (status OperationStatus) String() string {
	switch status {
	case OperationStatusSendToBank:
		return "SendToBank"
	case OperationStatusExecuted:
		return "Executed"
	case OperationStatusRejectByBank:
		return "RejectByBank"
	case OperationStatusCanceled:
		return "Canceled"
	case OperationStatusReceived:
		return "Received"
	default:
		return "unknown"
	}
}

func ParseOperationStatus(status string) (OperationStatus, error) {
	switch strings.ToLower(status) {
	case "sendtobank":
		return OperationStatusSendToBank, nil
	case "executed":
		return OperationStatusExecuted, nil
	case "rejectbybank":
		return OperationStatusRejectByBank, nil
	case "canceled":
		return OperationStatusCanceled, nil
	case "received":
		return OperationStatusReceived, nil
	}

	var s OperationStatus
	return s, fmt.Errorf("invalid operation status: %s", status)
}

type Operation struct {
	Id                          string
	CompanyId                   string
	Status                      OperationStatus
	Category                    OperationCategory
	ContragentName              string
	ContragentInn               string
	ContragentKpp               string
	ContragentBankAccountNumber string
	ContragentBankName          string
	ContragentBankBic           string
	Currency                    Currency
	Amount                      float64
	AmountWithCommission        float64
	BankAccountNumber           string
	PaymentPurpose              string
	Executed                    time.Time
	Created                     time.Time
	DocNumber                   string
	Kbk                         string // (104)
	Oktmo                       string // (105)
	PaymentBasis                string // (106)
	TaxCode                     string
	TaxDocNum                   string // (108)
	TaxDocDate                  string // (109)
	PayerStatus                 string // (101)
	Uin                         string
}

type auxOperation struct {
	Id                          string  `json:"id"`
	CompanyId                   string  `json:"companyId"`
	Status                      string  `json:"status"`
	Category                    string  `json:"category"`
	ContragentName              string  `json:"contragentName"`
	ContragentInn               string  `json:"contragentInn"`
	ContragentKpp               string  `json:"contragentKpp"`
	ContragentBankAccountNumber string  `json:"contragentBankAccountNumber"`
	ContragentBankName          string  `json:"contragentBankName"`
	ContragentBankBic           string  `json:"contragentBankBic"`
	Currency                    string  `json:"currency"`
	Amount                      float64 `json:"amount"`
	AmountWithCommission        float64 `json:"amountWithCommission"`
	BankAccountNumber           string  `json:"bankAccountNumber"`
	PaymentPurpose              string  `json:"paymentPurpose"`
	Executed                    string  `json:"executed"`
	Created                     string  `json:"created"`
	DocNumber                   string  `json:"docNumber"`
	Kbk                         string  `json:"kbk"`
	Oktmo                       string  `json:"oktmo"`
	PaymentBasis                string  `json:"paymentBasis"`
	TaxCode                     string  `json:"taxCode"`
	TaxDocNum                   string  `json:"taxDocNum"`
	TaxDocDate                  string  `json:"taxDocDate"`
	PayerStatus                 string  `json:"payerStatus"`
	Uin                         string  `json:"uin"`
}

func (operation *Operation) UnmarshalJSON(b []byte) error {
	var aux auxOperation

	err := json.Unmarshal(b, &aux)
	if err != nil {
		return err
	}

	operation.Id = aux.Id
	operation.CompanyId = aux.CompanyId
	operation.Status, err = ParseOperationStatus(aux.Status)
	if err != nil {
		return err
	}
	operation.Category, err = ParseOperationCategory(aux.Category)
	if err != nil {
		return err
	}
	operation.ContragentName = aux.ContragentName
	operation.ContragentInn = aux.ContragentInn
	operation.ContragentKpp = aux.ContragentKpp
	operation.ContragentBankAccountNumber = aux.ContragentBankAccountNumber
	operation.ContragentBankName = aux.ContragentBankName
	operation.ContragentBankBic = aux.ContragentBankBic
	operation.Currency, err = ParseCurrency(aux.Currency)
	if err != nil {
		return err
	}
	operation.Amount = aux.Amount
	operation.AmountWithCommission = aux.AmountWithCommission
	operation.BankAccountNumber = aux.BankAccountNumber
	operation.PaymentPurpose = aux.PaymentPurpose
	operation.Executed, err = time.ParseInLocation(DateFormat, aux.Executed, MskLocation)
	if err != nil {
		return err
	}
	operation.Created, err = time.ParseInLocation(DateFormat, aux.Created, MskLocation)
	if err != nil {
		return err
	}
	operation.DocNumber = aux.DocNumber
	operation.Kbk = aux.Kbk
	operation.Oktmo = aux.Oktmo
	operation.PaymentBasis = aux.PaymentBasis
	operation.TaxCode = aux.TaxCode
	operation.TaxDocNum = aux.TaxDocNum
	operation.TaxDocDate = aux.TaxDocDate
	operation.PayerStatus = aux.PayerStatus
	operation.Uin = aux.Uin

	return nil
}

func (operation *Operation) MarshalJSON() ([]byte, error) {
	return json.Marshal(&auxOperation{
		Id:                          operation.Id,
		CompanyId:                   operation.CompanyId,
		Status:                      operation.Status.String(),
		Category:                    operation.Category.String(),
		ContragentName:              operation.ContragentName,
		ContragentInn:               operation.ContragentInn,
		ContragentKpp:               operation.ContragentKpp,
		ContragentBankAccountNumber: operation.ContragentBankAccountNumber,
		ContragentBankName:          operation.ContragentBankName,
		ContragentBankBic:           operation.ContragentBankBic,
		Currency:                    operation.Currency.String(),
		Amount:                      operation.Amount,
		AmountWithCommission:        operation.AmountWithCommission,
		BankAccountNumber:           operation.BankAccountNumber,
		PaymentPurpose:              operation.PaymentPurpose,
		Executed:                    operation.Executed.In(MskLocation).Format(DateFormat),
		Created:                     operation.Created.In(MskLocation).Format(DateFormat),
		DocNumber:                   operation.DocNumber,
		Kbk:                         operation.Kbk,
		Oktmo:                       operation.Oktmo,
		PaymentBasis:                operation.PaymentBasis,
		TaxCode:                     operation.TaxCode,
		TaxDocNum:                   operation.TaxDocNum,
		TaxDocDate:                  operation.TaxDocDate,
		PayerStatus:                 operation.PayerStatus,
		Uin:                         operation.Uin,
	})
}
