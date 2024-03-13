package user

type UserInfo struct {
	Username            string      `json:"username" example:"user1"`                 // User's username
	Email               string      `json:"email" example:"user1@example.com"`        // User's email address
	IsAdmin             bool        `json:"is_admin" example:"false"`                 // Indicates if user has admin privileges
	PhoneNumber         string      `json:"phone_number" example:"9876543210"`        // User's phone number
	AuthenticationLevel int         `json:"authentication_level" example:"0"`         // User's authentication level
	BlockedLevel        int         `json:"blocked_level" example:"0"`                // User's blocked level
	Balance             int         `json:"balance" example:"5000000000"`             // User's account balance
	IsPremium           bool        `json:"is_premium" example:"false"`               // Indicates if user has a premium account
	BanksNames          []string    `json:"banks_names" example:"['saman', 'sepah']"` // List of user's bank names
	OpenTrades          interface{} `json:"open_trades,omitempty"`
	ClosedTrades        interface{} `json:"closed_trades,omitempty"`
	Transactions        interface{} `json:"transactions,omitempty"`
	Payments            interface{} `json:"payments,omitempty"`
}

func NewUserInfo(user User) UserInfo {
	newUserInfo := UserInfo{}
	newUserInfo.Username = user.Username
	newUserInfo.Email = user.Email
	newUserInfo.PhoneNumber = user.Profile.PhoneNumber
	newUserInfo.AuthenticationLevel = user.Profile.AuthenticationLevel
	newUserInfo.BlockedLevel = user.Profile.BlockedLevel
	newUserInfo.Balance = user.Profile.Balance
	newUserInfo.IsPremium = user.Profile.IsPremium

	for _, bi := range user.BankingInfo {
		newUserInfo.BanksNames = append(newUserInfo.BanksNames, bi.BankName)
	}

	return newUserInfo
}
