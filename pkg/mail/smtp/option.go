package smtp

type Options struct {
	User           string
	From           string
	Password       string
	Host           string
	Port           int
	InsureSecurity bool
}

func NewOption() *Options {
	return &Options{
		User:           "",
		From:           "",
		Password:       "",
		Host:           "",
		Port:           0,
		InsureSecurity: false,
	}
}

func NewOptionWithParmas(
	user string,
	password string,
	host string,
	port int,
	insureSecurity bool,
	from string,

) *Options {
	return &Options{
		User:           user,
		From:           from,
		Password:       password,
		Host:           host,
		Port:           port,
		InsureSecurity: insureSecurity,
	}
}
