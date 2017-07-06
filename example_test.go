package zarinpal_test

import (
	"bufio"
	"fmt"
	"os"

	"github.com/sijad/zarinpal-go"
)

func ExampleSimple() {
	merchantID := "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
	callbackurl := "http://yoursite.com/callbackurl"
	amount := 100
	description := "Description"

	// new request
	r := zarinpal.NewRequest(merchantID, callbackurl, amount, description)
	requestResponse, err := r.Request()

	if err != nil {
		// An error occured durring request.
		panic(err)
	}

	fmt.Println("Open folloing url and pay:")
	fmt.Println("https://www.zarinpal.com/pg/StartPay/" + requestResponse.Authority)

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Hit Enter after you've paid.")
	reader.ReadString('\n')

	// do verify
	v := zarinpal.NewVerify(merchantID, requestResponse.Authority, amount)
	verifyResponse, err := v.Verify()

	if err != nil {
		if verifyResponse != nil {
			// status code was not 100, so the verification has been faild.
			panic("Unsuccessful paymnet.")
		}

		// An error occured durring verification.
		panic(err)
	}

	fmt.Println("Successful paymnet.")
}
