package regru

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//const defaultApiEndpoint = "https://api.reg.ru/api/regru2/zone/add_alias"

//func Provider() *schema.Provider {
//	return &schema.Provider{
//		Schema: map[string]*schema.Schema{
//			"api_username": {
//				Type:        schema.TypeString,
//				Required:    true,
//				DefaultFunc: schema.EnvDefaultFunc("REGRU_API_USERNAME", nil),
//				Description: "API username for reg.ru",
//			},
//			"api_password": {
//				Type:        schema.TypeString,
//				Required:    true,
//				Sensitive:   true,
//				DefaultFunc: schema.EnvDefaultFunc("REGRU_API_PASSWORD", nil),
//				Description: "API password for reg.ru",
//			},
//			"api_endpoint": {
//				Type:        schema.TypeString,
//				Optional:    true,
//				Default:     defaultApiEndpoint,
//				Description: "reg.ru API endpoint",
//			},
//			"cert_file": {
//				Type:        schema.TypeString,
//				Optional:    true,
//				DefaultFunc: schema.EnvDefaultFunc("REGRU_CERT_FILE", nil),
//				Description: "Path to the client SSL certificate file",
//			},
//			"key_file": {
//				Type:        schema.TypeString,
//				Optional:    true,
//				DefaultFunc: schema.EnvDefaultFunc("REGRU_KEY_FILE", nil),
//				Description: "Path to the client SSL key file",
//			},
//		},
//		ResourcesMap: map[string]*schema.Resource{
//			"regru_dns_record": resourceRegruDNSRecord(),
//		},
//		ConfigureFunc: configureProvider,
//	}
//}

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_username": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("REGRU_API_USERNAME", nil),
				Description: "API username for reg.ru",
			},
			"api_password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("REGRU_API_PASSWORD", nil),
				Description: "API password for reg.ru",
			},
			"api_endpoint": {
				Type:        schema.TypeString,
				Optional:    true, // Оставляем Optional, но убираем Default
				Description: "reg.ru API endpoint",
			},
			"cert_file": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("REGRU_CERT_FILE", nil),
				Description: "Path to the client SSL certificate file",
			},
			"key_file": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("REGRU_KEY_FILE", nil),
				Description: "Path to the client SSL key file",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"regru_dns_record": resourceRegruDNSRecord(),
		},
		ConfigureFunc: configureProvider,
	}
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	// Получаем значения из конфигурации
	username := d.Get("api_username").(string)
	password := d.Get("api_password").(string)
	endpoint := d.Get("api_endpoint").(string)
	certFile := d.Get("cert_file").(string)
	keyFile := d.Get("key_file").(string)

	// Проверяем обязательные поля
	if username == "" || password == "" {
		return nil, fmt.Errorf("api_username and api_password are required")
	}
	if endpoint == "" {
		return nil, fmt.Errorf("api_endpoint is required")
	}

	// Создаем клиент
	client, err := NewClient(username, password, endpoint, certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	// Логируем успешное создание клиента
	log.Printf("Regru client created for user: %s", username)

	return client, nil
}
