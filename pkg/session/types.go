package session

type Session struct {
	Id                  int
	UserId              string `db:"user_id"`
	Tid                 string `db:"tid"`
	TppId               string `db:"tpp_id"`
	InternalAccessToken string `db:"internal_access_token"`
	ReferenceId         string `db:"reference_id"`
	AspspId             string `db:"aspsp_id"`
	CreateDateTime      string `db:"create_date_time"`
	UpdateDateTime      string `db:"update_date_time"`
}
