package zarinpal

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

const (
	baseurl         = "https://www.zarinpal.com/pg/rest/WebGate"
	requestEndpoint = baseurl + "PaymentRequest.json"
	verifyEndpoint  = baseurl + "PaymentVerification.json"
)

// RequestData holds request POST data.
type RequestData struct {
	MerchantID  string
	CallbackURL string
	Amount      int
	Description string
}

// Request implements sending request to create new transaction.
type Request struct {
	endpoint string
	data     RequestData
}

// RequestResponse holds request response data.
type RequestResponse struct {
	Status    int
	Authority string
	Errors    *map[string][]string `json:"errors"`
}

// Request send HTTP request to endpoint.
func (r *Request) Request() (*RequestResponse, error) {
	result := &RequestResponse{}
	if err := postData(r.endpoint, &r.data, result); err != nil {
		return nil, err
	}

	if result.Status != 100 {
		return result, errors.New("An error occurred")
	}

	return result, nil
}

// NewRequest creates new instance of Request.
func NewRequest(merchantID, callbackURL string, amount int, description string) Request {
	data := RequestData{merchantID, callbackURL, amount, description}
	return Request{requestEndpoint, data}
}

// VerifyData holds verify POST data.
type VerifyData struct {
	MerchantID string
	Authority  string
	Amount     int
}

// Verify implements sending request to verify a transaction.
type Verify struct {
	endpoint string
	data     VerifyData
}

// VerifyResponse holds verify response data.
type VerifyResponse struct {
	Status    int
	Authority string
	Errors    *map[string][]string `json:"errors"`
}

// Verify send HTTP request to endpoint.
func (v *Verify) Verify() (*VerifyResponse, error) {
	result := &VerifyResponse{}
	if err := postData(v.endpoint, &v.data, result); err != nil {
		return nil, err
	}

	if result.Status != 100 {
		return result, errors.New("An error occurred")
	}

	return result, nil
}

// NewVerify returns a new instance of Verify.
func NewVerify(merchantID, authority string, amount int) Verify {
	data := VerifyData{merchantID, authority, amount}
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
