package errors

import "fmt"

const loginMessage = `Please login to your server: 

astra login https://mycluster.mydomain.com
`

type Unauthorized struct {
}

func (u *Unauthorized) Error() string {
	return fmt.Sprintf("Unauthorized to access the cluster\n%s", loginMessage)
}
