package symbols

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/tidwall/gjson"
)

var modules string = "assetProfile,balanceSheetHistory,balanceSheetHistoryQuarterly,calendarEvents,cashflowStatementHistory,cashflowStatementHistoryQuarterly,defaultKeyStatistics,earnings,earningsHistory,earningsTrend,financialData,fundOwnership,incomeStatementHistory,incomeStatementHistoryQuarterly,indexTrend,industryTrend,insiderHolders,insiderTransactions,institutionOwnership,majorDirectHolders,majorHoldersBreakdown,netSharePurchaseActivity,price,quoteType,recommendationTrend,secFilings,sectorTrend,summaryDetail,summaryProfile,symbol,upgradeDowngradeHistory,fundProfile,topHoldings,fundPerformance"

type Client struct {
	baseUrl string
}

type SymbolData struct {
	// Symbol              string `json:"symbol"`
	LongBusinessSummary     string              `json:"longBusinessSummary"`
	CashflowStatements      []CashFlowStatement `json:"cashflowStatements"`
	ProfitMargins           float64             `json:"profitMargins"`
	SharesOutstanding       int                 `json:"sharesOutstanding"`
	Beta                    float64             `json:"beta"`
	BookValue               float64             `json:"bookValue"`
	PriceToBook             float64             `json:"priceToBook"`
	EarningsQuarterlyGrowth float64             `json:"earningsQuarterlyGrowth"`
}

type CashFlowStatement struct {
	EndDate                          string `json:"endDate"`
	NetIncome                        int    `json:"netIncome"`
	TotalCashFromOperatingActivities int    `json:"totalCashFromOperatingActivities"`
	NetBorrowings                    int    `json:"netBorrowings"`
}

func NewClient(silent bool) *Client {
	return &Client{
		baseUrl: "https://query2.finance.yahoo.com/v10/finance/quoteSummary",
	}
}

func (c *Client) MakeRequest(symbol string) (*SymbolData, error) {
	url := fmt.Sprintf(
		"%s/%s?modules=%s",
		c.baseUrl,
		strings.ToUpper(symbol),
		modules,
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return &SymbolData{}, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return &SymbolData{}, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &SymbolData{}, err
	}

	symbolData := SymbolData{}
	longBusinessSummary := gjson.Get(string(body), "quoteSummary.result.0.assetProfile.longBusinessSummary")
	if longBusinessSummary.Exists() {
		symbolData.LongBusinessSummary = longBusinessSummary.String()
	}

	cashflowStatements := gjson.Get(string(body), "quoteSummary.result.0.cashflowStatementHistory.cashflowStatements")
	if cashflowStatements.Exists() {
		for _, cs := range cashflowStatements.Array() {

			cashflowStatement := CashFlowStatement{}

			endDate := gjson.Get(cs.String(), "endDate.fmt")
			if endDate.Exists() {
				cashflowStatement.EndDate = endDate.String()
			}

			netIncome := gjson.Get(cs.String(), "netIncome.raw")
			if endDate.Exists() {
				cashflowStatement.NetIncome = int(netIncome.Int())
			}

			totalCashFromOperatingActivities := gjson.Get(cs.String(), "totalCashFromOperatingActivities.raw")
			if endDate.Exists() {
				cashflowStatement.TotalCashFromOperatingActivities = int(totalCashFromOperatingActivities.Int())
			}

			netBorrowings := gjson.Get(cs.String(), "netBorrowings.raw")
			if endDate.Exists() {
				cashflowStatement.NetBorrowings = int(netBorrowings.Int())
			}

			symbolData.CashflowStatements = append(symbolData.CashflowStatements, cashflowStatement)
		}
	}

	profitMargins := gjson.Get(string(body), "quoteSummary.result.0.defaultKeyStatistics.profitMargins.raw")
	if profitMargins.Exists() {
		symbolData.ProfitMargins = profitMargins.Float()
	}

	sharesOutstanding := gjson.Get(string(body), "quoteSummary.result.0.defaultKeyStatistics.sharesOutstanding.raw")
	if profitMargins.Exists() {
		symbolData.SharesOutstanding = int(sharesOutstanding.Int())
	}

	beta := gjson.Get(string(body), "quoteSummary.result.0.defaultKeyStatistics.beta.raw")
	if beta.Exists() {
		symbolData.Beta = beta.Float()
	}

	bookValue := gjson.Get(string(body), "quoteSummary.result.0.defaultKeyStatistics.bookValue.raw")
	if bookValue.Exists() {
		symbolData.BookValue = bookValue.Float()
	}

	priceToBook := gjson.Get(string(body), "quoteSummary.result.0.defaultKeyStatistics.priceToBook.raw")
	if priceToBook.Exists() {
		symbolData.PriceToBook = priceToBook.Float()
	}

	earningsQuarterlyGrowth := gjson.Get(string(body), "quoteSummary.result.0.defaultKeyStatistics.earningsQuarterlyGrowth.raw")
	if earningsQuarterlyGrowth.Exists() {
		symbolData.EarningsQuarterlyGrowth = earningsQuarterlyGrowth.Float()
	}

	return &symbolData, nil
}
