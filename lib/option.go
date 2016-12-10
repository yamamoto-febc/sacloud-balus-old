package lib

//import "fmt"

type Option struct {
	Path   string
	Port   int
	Secret string
	Debug  bool

	AccessToken       string
	AccessTokenSecret string
	Zone              []string
	TraceMode         bool

	AzureSubscriptionKey string

	MagicalSpel string
}

const MAGICAL_SPEL = "バルス"

func (o *Option) Validate() []error {
	var errors []error
	//if o.AccessToken == "" {
	//	errors = append(errors, fmt.Errorf("[%s] is required", "sakuracloud-access-token"))
	//}
	//if o.AccessTokenSecret == "" {
	//	errors = append(errors, fmt.Errorf("[%s] is required", "sakuracloud-access-token-secret"))
	//}
	//if o.AzureSubscriptionID == "" {
	//	errors = append(errors, fmt.Errorf("[%s] is required", "azure-subscription-id"))
	//}

	return errors
}

func NewOption() *Option {
	return &Option{
		MagicalSpel: MAGICAL_SPEL,
	}
}
