package zarinpal

import (
	"bytes"
	"encoding/json"
	"net/http"
)

const (
	requestEndpoint = "https://sandbox.zarinpal.com/pg/rest/WebGate/PaymentRequest.json"
	verifyEndpoint  = "https://sandbox.zarinpal.com/pg/rest/WebGate/PaymentVerification.json"
)

type RequestData struct {
	MerchantID  string
	CallbackURL string
	Amount      int
	Description string
}

type Request struct {
	endpoint string
	data     *RequestData
}

type requestResponse struct {
	Status    int
	Authority string
	Errors    *map[string][]string `json:errors`
}

func (r *Request) Request() (*requestResponse, error) {
	result := &requestResponse{}
	if err := postData(r.endpoint, &r.data, result); err != nil {
		return nil, err
	}

	// if result.Status != 100 {
	//     return result, errors.New("An error occurred.")
	// }

	return result, nil
}

func NewRequest(merchantID, callbackURL string, amount int, description string) Request {
	data := &RequestData{merchantID, callbackURL, amount, description}
	return Request{requestEndpoint, data}
}

type VerifyData struct {
	MerchantID string
	Authority  string
	Amount     int
}

type Verify struct {
	endpoint string
	data     *VerifyData
}

type verifyResponse struct {
	Status    int
	Authority string
	Errors    *map[string][]string `json:errors`
}

func (v *Verify) Verify() (*verifyResponse, error) {
	result := &verifyResponse{}
	if err := postData(v.endpoint, &v.data, result); err != nil {
		return nil, err
	}

	// if result.Status != 100 {
	//     return result, errors.New("An error occurred.")
	// }

	return result, nil
}

func NewVerify(merchantID, authority string, amount int) Verify {
	data := &VerifyData{merchantID, authority, amount}
	return Verify{verifyEndpoint, data}
}

func postData(url string, data interface{}, result interface{}) error {
	d, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(d))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(resp.Body); err != nil {
		return err
	}

	return json.Unmarshal(buf.Bytes(), result)
}
