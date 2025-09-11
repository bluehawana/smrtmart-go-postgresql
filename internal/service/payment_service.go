package service

import (
	"encoding/json"
	"fmt"
	"log"

	"smrtmart-go-postgresql/internal/config"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
	"github.com/stripe/stripe-go/v76/webhook"
)

type PaymentService interface {
	CreateCheckoutSession(items []CheckoutItem, customerEmail string, successURL, cancelURL string) (*stripe.CheckoutSession, error)
	CreateCheckoutSessionWithFullInfo(items []CheckoutItem, customerInfo CustomerInfo, shippingAddress Address, billingAddress *Address, successURL, cancelURL string) (*stripe.CheckoutSession, error)
	HandleWebhook(payload []byte, signature string) error
}

type CustomerInfo struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
}

type Address struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Company      string `json:"company,omitempty"`
	AddressLine1 string `json:"address_line1"`
	AddressLine2 string `json:"address_line2,omitempty"`
	City         string `json:"city"`
	State        string `json:"state"`
	PostalCode   string `json:"postal_code"`
	Country      string `json:"country"`
	Phone        string `json:"phone,omitempty"`
}

type CheckoutItem struct {
	ProductID   string  `json:"product_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	Images      []string `json:"images"`
}

type paymentService struct {
	stripeConfig config.StripeConfig
}

func NewPaymentService(stripeConfig config.StripeConfig) PaymentService {
	// Initialize Stripe
	stripe.Key = stripeConfig.SecretKey
	return &paymentService{stripeConfig: stripeConfig}
}

func (s *paymentService) CreateCheckoutSession(items []CheckoutItem, customerEmail string, successURL, cancelURL string) (*stripe.CheckoutSession, error) {
	var lineItems []*stripe.CheckoutSessionLineItemParams

	for _, item := range items {
		// Convert price to cents (Stripe uses cents)
		priceInCents := int64(item.Price * 100)

    lineItem := &stripe.CheckoutSessionLineItemParams{
        PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
            Currency: stripe.String("sek"),
				ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
					Name:        stripe.String(item.Name),
					Description: stripe.String(item.Description),
				},
				UnitAmount: stripe.Int64(priceInCents),
			},
			Quantity: stripe.Int64(int64(item.Quantity)),
		}

		// Add product images if available
		if len(item.Images) > 0 {
			var images []*string
			for _, img := range item.Images {
				// Construct full image URL
				imageURL := fmt.Sprintf("http://localhost:8080/uploads/%s", img)
				images = append(images, stripe.String(imageURL))
			}
			lineItem.PriceData.ProductData.Images = images
		}

		lineItems = append(lineItems, lineItem)
	}

    params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems:          lineItems,
		Mode:               stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL:         stripe.String(successURL),
		CancelURL:          stripe.String(cancelURL),
		CustomerEmail:      stripe.String(customerEmail),
		
		// Enable shipping address collection
        ShippingAddressCollection: &stripe.CheckoutSessionShippingAddressCollectionParams{
            AllowedCountries: stripe.StringSlice([]string{"SE", "NO", "DK", "FI", "GB", "DE", "FR", "ES", "IT", "NL", "BE", "AT", "CH", "US", "CA"}),
        },
		
		// Add metadata for order tracking
		Metadata: map[string]string{
			"source": "smrtmart_website",
		},
	}

	sess, err := session.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create checkout session: %w", err)
	}

	return sess, nil
}

func (s *paymentService) CreateCheckoutSessionWithFullInfo(items []CheckoutItem, customerInfo CustomerInfo, shippingAddress Address, billingAddress *Address, successURL, cancelURL string) (*stripe.CheckoutSession, error) {
	var lineItems []*stripe.CheckoutSessionLineItemParams

	for _, item := range items {
		// Convert price to cents (Stripe uses cents)
		priceInCents := int64(item.Price * 100)

    lineItem := &stripe.CheckoutSessionLineItemParams{
        PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
            Currency: stripe.String("sek"),
				ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
					Name:        stripe.String(item.Name),
					Description: stripe.String(item.Description),
				},
				UnitAmount: stripe.Int64(priceInCents),
			},
			Quantity: stripe.Int64(int64(item.Quantity)),
		}

		// Add product images if available
		if len(item.Images) > 0 {
			var images []*string
			for _, img := range item.Images {
				// Use Supabase image URL for production
				imageURL := fmt.Sprintf("https://mqkoydypybxgcwxioqzc.supabase.co/storage/v1/object/public/products/%s", img)
				images = append(images, stripe.String(imageURL))
			}
			lineItem.PriceData.ProductData.Images = images
		}

		lineItems = append(lineItems, lineItem)
	}

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems:          lineItems,
		Mode:               stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL:         stripe.String(successURL),
		CancelURL:          stripe.String(cancelURL),
		CustomerEmail:      stripe.String(customerInfo.Email),
		
		// Enable shipping address collection
        ShippingAddressCollection: &stripe.CheckoutSessionShippingAddressCollectionParams{
            AllowedCountries: stripe.StringSlice([]string{"SE", "NO", "DK", "FI", "GB", "DE", "FR", "ES", "IT", "NL", "BE", "AT", "CH", "US", "CA"}),
        },
		
		// Pre-fill shipping address
		ShippingOptions: []*stripe.CheckoutSessionShippingOptionParams{
			{
				ShippingRateData: &stripe.CheckoutSessionShippingOptionShippingRateDataParams{
					Type: stripe.String("fixed_amount"),
                FixedAmount: &stripe.CheckoutSessionShippingOptionShippingRateDataFixedAmountParams{
                    Amount:   stripe.Int64(0), // Free shipping
                    Currency: stripe.String("sek"),
                },
					DisplayName: stripe.String("Free Shipping"),
				},
			},
			{
				ShippingRateData: &stripe.CheckoutSessionShippingOptionShippingRateDataParams{
					Type: stripe.String("fixed_amount"),
                FixedAmount: &stripe.CheckoutSessionShippingOptionShippingRateDataFixedAmountParams{
                    Amount:   stripe.Int64(9900), // 99.00 SEK
                    Currency: stripe.String("sek"),
                },
					DisplayName: stripe.String("Express Shipping"),
					DeliveryEstimate: &stripe.CheckoutSessionShippingOptionShippingRateDataDeliveryEstimateParams{
						Minimum: &stripe.CheckoutSessionShippingOptionShippingRateDataDeliveryEstimateMinimumParams{
							Unit:  stripe.String("business_day"),
							Value: stripe.Int64(1),
						},
						Maximum: &stripe.CheckoutSessionShippingOptionShippingRateDataDeliveryEstimateMaximumParams{
							Unit:  stripe.String("business_day"),
							Value: stripe.Int64(3),
						},
					},
				},
			},
		},
		
		// Add metadata for order tracking
		Metadata: map[string]string{
			"source":                "smrtmart_website_full_info",
			"customer_first_name":   customerInfo.FirstName,
			"customer_last_name":    customerInfo.LastName,
			"customer_phone":        customerInfo.Phone,
			"shipping_address_line1": shippingAddress.AddressLine1,
			"shipping_city":         shippingAddress.City,
			"shipping_state":        shippingAddress.State,
			"shipping_country":      shippingAddress.Country,
		},
	}

	// Add billing address metadata if different from shipping
	if billingAddress != nil {
		params.Metadata["billing_address_line1"] = billingAddress.AddressLine1
		params.Metadata["billing_city"] = billingAddress.City
		params.Metadata["billing_state"] = billingAddress.State
		params.Metadata["billing_country"] = billingAddress.Country
	}

	sess, err := session.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create checkout session: %w", err)
	}

	return sess, nil
}

func (s *paymentService) HandleWebhook(payload []byte, signature string) error {
	event, err := webhook.ConstructEvent(payload, signature, s.stripeConfig.WebhookSecret)
	if err != nil {
		return fmt.Errorf("failed to verify webhook signature: %w", err)
	}

	switch event.Type {
	case "checkout.session.completed":
		var session stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
			return fmt.Errorf("failed to unmarshal session: %w", err)
		}
		
		log.Printf("Payment successful for session: %s", session.ID)
		// TODO: Create order in database, update inventory, send confirmation email
		
	case "payment_intent.succeeded":
		var paymentIntent stripe.PaymentIntent
		if err := json.Unmarshal(event.Data.Raw, &paymentIntent); err != nil {
			return fmt.Errorf("failed to unmarshal payment intent: %w", err)
		}
		
		log.Printf("Payment intent succeeded: %s", paymentIntent.ID)
		
	case "payment_intent.payment_failed":
		var paymentIntent stripe.PaymentIntent
		if err := json.Unmarshal(event.Data.Raw, &paymentIntent); err != nil {
			return fmt.Errorf("failed to unmarshal payment intent: %w", err)
		}
		
		log.Printf("Payment failed: %s", paymentIntent.ID)
		
	default:
		log.Printf("Unhandled event type: %s", event.Type)
	}

	return nil
}
