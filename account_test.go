package rediscloud_api

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/RedisLabs/rediscloud-go-api/service/account"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccount_ListPayments(t *testing.T) {
	s := httptest.NewServer(testServer("/payment-methods", "apiKey", "secret", `{
  "accountId": 1,
  "paymentMethods": [
    {
      "id": 123,
      "type": "Visa",
      "creditCardEndsWith": 9876,
      "nameOnCard": "Guy Incognito",
      "expirationMonth": 1,
      "expirationYear": 2021
    },
    {
      "id": 654,
      "type": "Mastercard",
      "creditCardEndsWith": 4567,
      "nameOnCard": "Joey JoJo Junior Shabadoo",
      "expirationMonth": 2,
      "expirationYear": 2022
    }
  ],
  "_links": {
    "self": {
      "href": "https://example.com",
      "type": "GET"
    }
  }
}`))

	subject, err := NewClient(BaseUrl(s.URL), Auth("apiKey", "secret"), Transporter(s.Client().Transport))
	require.NoError(t, err)

	actual, err := subject.Account.ListPaymentMethods(context.TODO())
	require.NoError(t, err)

	assert.ElementsMatch(t, []account.PaymentMethod{
		{
			Id:                 123,
			Type:               "Visa",
			CreditCardEndsWith: 9876,
		},
		{
			Id:                 654,
			Type:               "Mastercard",
			CreditCardEndsWith: 4567,
		},
	}, actual)
}
