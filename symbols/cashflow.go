package symbols

type CashFlowStatement struct {
	EndDate                          string `json:"endDate"`
	NetIncome                        int    `json:"netIncome"`
	TotalCashFromOperatingActivities int    `json:"totalCashFromOperatingActivities"`
	NetBorrowings                    int    `json:"netBorrowings"`
}
