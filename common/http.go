package common

import tls_client "github.com/bogdanfinn/tls-client"
import "os"

func NewClient() tls_client.HttpClient {
	jar := tls_client.NewCookieJar()
	options := []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(360),
		tls_client.WithClientProfile(tls_client.Okhttp4Android13),
		tls_client.WithNotFollowRedirects(),
		tls_client.WithCookieJar(jar),
		tls_client.WithInsecureSkipVerify(),
	}

	if proxyURL, exists := os.LookupEnv("PROXY_URL"); exists && proxyURL != "" {
		options = append(options, tls_client.WithProxyUrl(proxyURL))
	}

	client, _ := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)

	return client
}
