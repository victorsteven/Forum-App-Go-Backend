package errormsg

type ErrorMessage struct {
	Required_email    string
	Required_nickname string
	Invalid_email     string
	Required_password string
	Invalid_body      string
	Unmarshal_error   string
}
