package domain

import (
	"errors"
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/transip/gotransip/v6"
	"github.com/transip/gotransip/v6/internal/testutil"
	"github.com/transip/gotransip/v6/rest"
)

const (
	error404Response   = `{ "error": "Domain with name 'example2.com' not found" }`
	domainsAPIResponse = `{ "domains": [
    {
      "name": "example.com",
      "authCode": "kJqfuOXNOYQKqh/jO4bYSn54YDqgAt1ksCe+ZG4Ud",
      "isTransferLocked": false,
      "registrationDate": "2016-01-01",
      "renewalDate": "2020-01-01",
      "isWhitelabel": false,
      "cancellationDate": "2020-01-01 12:00:00",
      "cancellationStatus": "signed",
      "isDnsOnly": false,
      "tags": [ "customTag", "anotherTag" ]
    }
  ] }`
	domainAPIResponse = `{ "domain": {
    "name": "example.com",
    "authCode": "kJqfuOXNOYQKqh/jO4bYSn54YDqgAt1ksCe+ZG4Ud",
    "isTransferLocked": false,
    "registrationDate": "2016-01-01",
    "renewalDate": "2020-01-01",
    "isWhitelabel": false,
    "cancellationDate": "2020-01-01 12:00:00",
    "cancellationStatus": "signed",
    "isDnsOnly": false,
    "tags": [ "customTag", "anotherTag" ]
  } } `
	brandingAPIResponse = `{
		"branding": {
		"companyName": "Example B.V.",
		"supportEmail": "admin@example.com",
		"companyUrl": "www.example.com",
		"termsOfUsageUrl": "www.example.com/tou",
		"bannerLine1": "Example B.V.",
		"bannerLine2": "Example",
		"bannerLine3": "http://www.example.com/products"
	} }`
	contactsAPIResponse = `{ "contacts": [ {
      "type": "registrant",
      "firstName": "John",
      "lastName": "Doe",
      "companyName": "Example B.V.",
      "companyKvk": "83057825",
      "companyType": "BV",
      "street": "Easy street",
      "number": "12",
      "postalCode": "1337 XD",
      "city": "Leiden",
      "phoneNumber": "+31 715241919",
      "faxNumber": "+31 715241919",
      "email": "example@example.com",
      "country": "nl"
    } ] }`
	dnsEntriesAPIResponse = `{ "dnsEntries": [
    { "name": "www", "expire": 86400, "type": "A", "content": "127.0.0.1" }
  ] }`
	dnsSecEntriesAPIResponseRequest = `{ "dnsSecEntries": [ {
      "keyTag": 67239,
      "flags": 1,
      "algorithm": 8,
      "publicKey": "kljlfkjsdfkjasdklf="
    } ] }`
)

func TestRepository_GetAll(t *testing.T) {
	server := testutil.MockServer{T: t, ExpectedMethod: "GET", ExpectedURL: "/domains", StatusCode: 200, Response: domainsAPIResponse}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := Repository{Client: *client}

	all, err := repo.GetAll()
	require.NoError(t, err)
	require.Equal(t, 1, len(all))
	assert.Equal(t, "example.com", all[0].Name)
	assert.Equal(t, "kJqfuOXNOYQKqh/jO4bYSn54YDqgAt1ksCe+ZG4Ud", all[0].AuthCode)
	assert.Equal(t, false, all[0].IsTransferLocked)
	assert.Equal(t, "2016-01-01 00:00:00", all[0].RegistrationDate.Format("2006-01-02 15:04:05"))
	assert.Equal(t, "2020-01-01 00:00:00", all[0].RenewalDate.Format("2006-01-02 15:04:05"))
	assert.Equal(t, false, all[0].IsWhitelabel)
	assert.Equal(t, "2020-01-01 12:00:00", all[0].CancellationDate.Format("2006-01-02 15:04:05"))
	assert.Equal(t, "signed", all[0].CancellationStatus)
	assert.Equal(t, false, all[0].IsDNSOnly)
	assert.Equal(t, []string{"customTag", "anotherTag"}, all[0].Tags)
}

func TestRepository_GetSelection(t *testing.T) {
	server := testutil.MockServer{T: t, ExpectedMethod: "GET", ExpectedURL: "/domains?page=1&pageSize=25", StatusCode: 200, Response: domainsAPIResponse}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := Repository{Client: *client}

	all, err := repo.GetSelection(1, 25)
	require.NoError(t, err)
	require.Equal(t, 1, len(all))
	assert.Equal(t, "example.com", all[0].Name)
	assert.Equal(t, "kJqfuOXNOYQKqh/jO4bYSn54YDqgAt1ksCe+ZG4Ud", all[0].AuthCode)
	assert.Equal(t, false, all[0].IsTransferLocked)
	assert.Equal(t, "2016-01-01 00:00:00", all[0].RegistrationDate.Format("2006-01-02 15:04:05"))
	assert.Equal(t, "2020-01-01 00:00:00", all[0].RenewalDate.Format("2006-01-02 15:04:05"))
	assert.Equal(t, false, all[0].IsWhitelabel)
	assert.Equal(t, "2020-01-01 12:00:00", all[0].CancellationDate.Format("2006-01-02 15:04:05"))
	assert.Equal(t, "signed", all[0].CancellationStatus)
	assert.Equal(t, false, all[0].IsDNSOnly)
	assert.Equal(t, []string{"customTag", "anotherTag"}, all[0].Tags)
}

func TestRepository_GetAllByTags(t *testing.T) {
	server := testutil.MockServer{T: t, ExpectedMethod: "GET", ExpectedURL: "/domains?tags=customTag", StatusCode: 200, Response: domainsAPIResponse}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := Repository{Client: *client}

	all, err := repo.GetAllByTags([]string{"customTag"})
	require.NoError(t, err)
	require.Equal(t, 1, len(all))

	assert.Equal(t, "example.com", all[0].Name)
	assert.Equal(t, "kJqfuOXNOYQKqh/jO4bYSn54YDqgAt1ksCe+ZG4Ud", all[0].AuthCode)
	assert.Equal(t, false, all[0].IsTransferLocked)
	assert.Equal(t, "2016-01-01 00:00:00", all[0].RegistrationDate.Format("2006-01-02 15:04:05"))
	assert.Equal(t, "2020-01-01 00:00:00", all[0].RenewalDate.Format("2006-01-02 15:04:05"))
	assert.Equal(t, false, all[0].IsWhitelabel)
	assert.Equal(t, "2020-01-01 12:00:00", all[0].CancellationDate.Format("2006-01-02 15:04:05"))
	assert.Equal(t, "signed", all[0].CancellationStatus)
	assert.Equal(t, false, all[0].IsDNSOnly)
	assert.Equal(t, []string{"customTag", "anotherTag"}, all[0].Tags)
}

func TestRepository_GetByDomainName(t *testing.T) {
	server := testutil.MockServer{T: t, ExpectedMethod: "GET", ExpectedURL: "/domains/example.com", StatusCode: 200, Response: domainAPIResponse}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := Repository{Client: *client}

	domain, err := repo.GetByDomainName("example.com")
	require.NoError(t, err)
	assert.Equal(t, "example.com", domain.Name)
	assert.Equal(t, "kJqfuOXNOYQKqh/jO4bYSn54YDqgAt1ksCe+ZG4Ud", domain.AuthCode)
	assert.Equal(t, false, domain.IsTransferLocked)
	assert.Equal(t, "2016-01-01 00:00:00", domain.RegistrationDate.Format("2006-01-02 15:04:05"))
	assert.Equal(t, "2020-01-01 00:00:00", domain.RenewalDate.Format("2006-01-02 15:04:05"))
	assert.Equal(t, false, domain.IsWhitelabel)
	assert.Equal(t, "2020-01-01 12:00:00", domain.CancellationDate.Format("2006-01-02 15:04:05"))
	assert.Equal(t, "signed", domain.CancellationStatus)
	assert.Equal(t, false, domain.IsDNSOnly)
	assert.Equal(t, []string{"customTag", "anotherTag"}, domain.Tags)
}

func TestRepository_GetByDomainNameError(t *testing.T) {
	domainName := "example2.com"
	server := testutil.MockServer{T: t, ExpectedMethod: "GET", ExpectedURL: fmt.Sprintf("/domains/%s", domainName), StatusCode: 404, Response: error404Response}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := Repository{Client: *client}

	domain, err := repo.GetByDomainName(domainName)
	if assert.Errorf(t, err, "getbydomainname server response error not returned") {
		require.Empty(t, domain.Name)
		assert.Equal(t, &rest.Error{Message: "Domain with name 'example2.com' not found", StatusCode: 404}, err)
	}
}

func TestRepository_Register(t *testing.T) {
	expectedRequest := `{"domainName":"example.com"}`
	server := testutil.MockServer{T: t, ExpectedMethod: "POST", ExpectedURL: "/domains", StatusCode: 201, ExpectedRequest: expectedRequest}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := Repository{Client: *client}

	register := Register{DomainName: "example.com"}
	err := repo.Register(register)
	require.NoError(t, err)
}

func TestRepository_RegisterError(t *testing.T) {
	errorResponse := `{"error":"The domain 'example.com' is not free and thus cannot be registered"}`
	server := testutil.MockServer{T: t, ExpectedMethod: "POST", ExpectedURL: "/domains", StatusCode: 406, SkipRequestBody: true, Response: errorResponse}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := Repository{Client: *client}

	register := Register{DomainName: "example.com"}
	err := repo.Register(register)
	if assert.Errorf(t, err, "register server response error not returned") {
		assert.Error(t, errors.New("The domain 'example.com' is not free and thus cannot be registered"), err)
	}
}

func TestRepository_Transfer(t *testing.T) {
	expectedRequest := `{"domainName":"example.com","authCode":"test123"}`
	server := testutil.MockServer{T: t, ExpectedMethod: "POST", ExpectedURL: "/domains", StatusCode: 201, ExpectedRequest: expectedRequest}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := Repository{Client: *client}

	transfer := Transfer{DomainName: "example.com", AuthCode: "test123"}

	err := repo.Transfer(transfer)
	require.NoError(t, err)
}

func TestRepository_TransferError(t *testing.T) {
	errorResponse := `{"error":"The domain 'example.com' is not registered and thus cannot be transferred"}`
	server := testutil.MockServer{T: t, ExpectedMethod: "POST", ExpectedURL: "/domains", StatusCode: 409, SkipRequestBody: true, Response: errorResponse}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := Repository{Client: *client}

	transfer := Transfer{DomainName: "example.com", AuthCode: "test123"}
	err := repo.Transfer(transfer)

	if assert.Errorf(t, err, "transfer server response error not returned") {
		assert.Error(t, errors.New("The domain 'example.com' is not registered and thus cannot be transferred"), err)
	}
}

func TestRepository_Update(t *testing.T) {
	expectedRequest := `{"domain":{"tags":["test123","test1234"],"cancellationDate":"0001-01-01T00:00:00Z","isTransferLocked":false,"isWhitelabel":false,"name":"example.com","registrationDate":"0001-01-01T00:00:00Z","renewalDate":"0001-01-01T00:00:00Z"}}`
	server := testutil.MockServer{T: t, ExpectedMethod: "PUT", ExpectedURL: "/domains/example.com", StatusCode: 204, ExpectedRequest: expectedRequest}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := Repository{Client: *client}

	domain := Domain{Tags: []string{"test123", "test1234"}, IsTransferLocked: false, IsWhitelabel: false, Name: "example.com"}

	err := repo.Update(domain)
	require.NoError(t, err)
}

func TestRepository_CancelEnd(t *testing.T) {
	expectedRequest := `{"endTime":"end"}`
	server := testutil.MockServer{T: t, ExpectedMethod: "DELETE", ExpectedURL: "/domains/example.com", StatusCode: 204, ExpectedRequest: expectedRequest}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := Repository{Client: *client}

	err := repo.Cancel("example.com", gotransip.CancellationTimeEnd)
	require.NoError(t, err)
}

func TestRepository_Cancel(t *testing.T) {
	expectedRequest := `{"endTime":"immediately"}`
	server := testutil.MockServer{T: t, ExpectedMethod: "DELETE", ExpectedURL: "/domains/example.com", StatusCode: 204, ExpectedRequest: expectedRequest}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := Repository{Client: *client}

	err := repo.Cancel("example.com", gotransip.CancellationTimeImmediately)
	require.NoError(t, err)
}

func TestRepository_GetDomainBranding(t *testing.T) {
	domainName := "example2.com"
	server := testutil.MockServer{T: t, ExpectedMethod: "GET", ExpectedURL: fmt.Sprintf("/domains/%s/branding", domainName), StatusCode: 200, Response: brandingAPIResponse}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := Repository{Client: *client}

	branding, err := repo.GetBranding(domainName)
	require.NoError(t, err)

	assert.Equal(t, "Example B.V.", branding.CompanyName)
	assert.Equal(t, "admin@example.com", branding.SupportEmail)
	assert.Equal(t, "www.example.com", branding.CompanyURL)
	assert.Equal(t, "www.example.com/tou", branding.TermsOfUsageURL)
	assert.Equal(t, "Example B.V.", branding.BannerLine1)
	assert.Equal(t, "Example", branding.BannerLine2)
	assert.Equal(t, "http://www.example.com/products", branding.BannerLine3)
}

func TestRepository_UpdateDomainBranding(t *testing.T) {
	expectedRequest := `{"branding":{"bannerLine1":"Example B.V.","bannerLine2":"admin@example.com","bannerLine3":"www.example.com","companyName":"www.example.com/tou","companyUrl":"Example B.V.","supportEmail":"Example","termsOfUsageUrl":"http://www.example.com/products"}}`
	server := testutil.MockServer{T: t, ExpectedMethod: "PUT", ExpectedURL: "/domains/example.com/branding", StatusCode: 204, ExpectedRequest: expectedRequest}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := Repository{Client: *client}

	branding := Branding{
		BannerLine1:     "Example B.V.",
		BannerLine2:     "admin@example.com",
		BannerLine3:     "www.example.com",
		CompanyName:     "www.example.com/tou",
		CompanyURL:      "Example B.V.",
		SupportEmail:    "Example",
		TermsOfUsageURL: "http://www.example.com/products",
	}

	err := repo.UpdateBranding("example.com", branding)
	require.NoError(t, err)
}

func TestRepository_GetContacts(t *testing.T) {
	domainName := "example.com"
	server := testutil.MockServer{T: t, ExpectedMethod: "GET", ExpectedURL: fmt.Sprintf("/domains/%s/contacts", domainName), StatusCode: 200, Response: contactsAPIResponse}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := Repository{Client: *client}

	contacts, err := repo.GetContacts(domainName)
	require.NoError(t, err)
	require.Equal(t, 1, len(contacts))

	assert.Equal(t, "registrant", contacts[0].Type)
	assert.Equal(t, "John", contacts[0].FirstName)
	assert.Equal(t, "Doe", contacts[0].LastName)
	assert.Equal(t, "Example B.V.", contacts[0].CompanyName)
	assert.Equal(t, "83057825", contacts[0].CompanyKvk)
	assert.Equal(t, "BV", contacts[0].CompanyType)
	assert.Equal(t, "Easy street", contacts[0].Street)
	assert.Equal(t, "12", contacts[0].Number)
	assert.Equal(t, "1337 XD", contacts[0].PostalCode)
	assert.Equal(t, "Leiden", contacts[0].City)
	assert.Equal(t, "+31 715241919", contacts[0].PhoneNumber)
	assert.Equal(t, "+31 715241919", contacts[0].FaxNumber)
	assert.Equal(t, "example@example.com", contacts[0].Email)
	assert.Equal(t, "nl", contacts[0].Country)
}

func TestRepository_UpdateContacts(t *testing.T) {
	expectedRequest := `{"contacts":[{"type":"registrant","firstName":"John","lastName":"Doe","companyName":"Example B.V.","companyKvk":"83057825","companyType":"BV","street":"Easy street","number":"12","postalCode":"1337 XD","city":"Leiden","phoneNumber":"+31 715241919","faxNumber":"+31 715241919","email":"example@example.com","country":"nl"}]}`
	server := testutil.MockServer{T: t, ExpectedMethod: "PUT", ExpectedURL: "/domains/example.com/contacts", StatusCode: 204, ExpectedRequest: expectedRequest}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := Repository{Client: *client}

	contacts := []WhoisContact{
		{
			Type:        "registrant",
			FirstName:   "John",
			LastName:    "Doe",
			CompanyName: "Example B.V.",
			CompanyKvk:  "83057825",
			CompanyType: "BV",
			Street:      "Easy street",
			Number:      "12",
			PostalCode:  "1337 XD",
			City:        "Leiden",
			PhoneNumber: "+31 715241919",
			FaxNumber:   "+31 715241919",
			Email:       "example@example.com",
			Country:     "nl",
		},
	}

	err := repo.UpdateContacts("example.com", contacts)
	require.NoError(t, err)
}

func TestRepository_GetDnsEntries(t *testing.T) {
	server := testutil.MockServer{T: t, ExpectedMethod: "GET", ExpectedURL: "/domains/example.com/dns", StatusCode: 200, Response: dnsEntriesAPIResponse}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := Repository{Client: *client}

	entries, err := repo.GetDNSEntries("example.com")
	require.Equal(t, 1, len(entries))
	require.NoError(t, err)
	assert.Equal(t, "www", entries[0].Name)
	assert.Equal(t, 86400, entries[0].Expire)
	assert.Equal(t, "A", entries[0].Type)
	assert.Equal(t, "127.0.0.1", entries[0].Content)

}

func TestRepository_AddDnsEntry(t *testing.T) {
	expectedRequest := `{"dnsEntry":{"name":"www","expire":1337,"type":"A","content":"127.0.0.1"}}`
	server := testutil.MockServer{T: t, ExpectedMethod: "POST", ExpectedURL: "/domains/example.com/dns", StatusCode: 201, ExpectedRequest: expectedRequest}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := Repository{Client: *client}

	dnsEntry := DNSEntry{Content: "127.0.0.1", Expire: 1337, Name: "www", Type: "A"}
	err := repo.AddDNSEntry("example.com", dnsEntry)
	require.NoError(t, err)
}

func TestRepository_UpdateDnsEntry(t *testing.T) {
	expectedRequest := `{"dnsEntry":{"name":"www","expire":1337,"type":"A","content":"127.0.0.1"}}`
	server := testutil.MockServer{T: t, ExpectedMethod: "PATCH", ExpectedURL: "/domains/example.com/dns", StatusCode: 204, ExpectedRequest: expectedRequest}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := Repository{Client: *client}

	dnsEntry := DNSEntry{Content: "127.0.0.1", Expire: 1337, Name: "www", Type: "A"}
	err := repo.UpdateDNSEntry("example.com", dnsEntry)
	require.NoError(t, err)
}

func TestRepository_ReplaceDnsEntries(t *testing.T) {
	expectedRequest := `{"dnsEntries":[{"name":"www","expire":1337,"type":"A","content":"127.0.0.1"}]}`
	server := testutil.MockServer{T: t, ExpectedMethod: "PUT", ExpectedURL: "/domains/example.com/dns", StatusCode: 204, ExpectedRequest: expectedRequest}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := Repository{Client: *client}

	dnsEntries := []DNSEntry{{Content: "127.0.0.1", Expire: 1337, Name: "www", Type: "A"}}
	err := repo.ReplaceDNSEntries("example.com", dnsEntries)
	require.NoError(t, err)
}

func TestRepository_RemoveDnsEntry(t *testing.T) {
	expectedRequest := `{"dnsEntry":{"name":"www","expire":1337,"type":"A","content":"127.0.0.1"}}`
	server := testutil.MockServer{T: t, ExpectedMethod: "DELETE", ExpectedURL: "/domains/example.com/dns", StatusCode: 204, ExpectedRequest: expectedRequest}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := Repository{Client: *client}

	dnsEntry := DNSEntry{Content: "127.0.0.1", Expire: 1337, Name: "www", Type: "A"}
	err := repo.RemoveDNSEntry("example.com", dnsEntry)
	require.NoError(t, err)
}

func TestRepository_GetDnsSecEntries(t *testing.T) {
	server := testutil.MockServer{T: t, ExpectedMethod: "GET", ExpectedURL: "/domains/example.com/dnssec", StatusCode: 200, Response: dnsSecEntriesAPIResponseRequest}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := Repository{Client: *client}

	entries, err := repo.GetDNSSecEntries("example.com")
	require.NoError(t, err)
	require.Equal(t, 1, len(entries))

	assert.Equal(t, 67239, entries[0].KeyTag)
	assert.Equal(t, 1, entries[0].Flags)
	assert.Equal(t, 8, entries[0].Algorithm)
	assert.Equal(t, "kljlfkjsdfkjasdklf=", entries[0].PublicKey)
}

func TestRepository_ReplaceDnsSecEntries(t *testing.T) {
	expectedRequestBody := `{"dnsSecEntries":[{"algorithm":8,"flags":1,"keyTag":67239,"publicKey":"test123"}]}`
	server := testutil.MockServer{T: t, ExpectedMethod: "PUT", ExpectedURL: "/domains/example.com/dnssec", StatusCode: 204, ExpectedRequest: expectedRequestBody}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := Repository{Client: *client}

	dnsSecEntries := []DNSSecEntry{{KeyTag: 67239, Flags: 1, Algorithm: 8, PublicKey: "test123"}}
	err := repo.ReplaceDNSSecEntries("example.com", dnsSecEntries)
	require.NoError(t, err)
}

func TestRepository_GetNameservers(t *testing.T) {
	apiResponse := `{"nameservers":[{"hostname":"ns0.transip.nl","ipv4":"127.0.0.1","ipv6":"2a01::1"}]}`
	server := testutil.MockServer{T: t, ExpectedMethod: "GET", ExpectedURL: "/domains/example.com/nameservers", StatusCode: 200, Response: apiResponse}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := Repository{Client: *client}

	nameservers, err := repo.GetNameservers("example.com")
	require.NoError(t, err)
	assert.Equal(t, 1, len(nameservers))
	assert.Equal(t, "ns0.transip.nl", nameservers[0].Hostname)
	assert.Equal(t, "127.0.0.1", nameservers[0].IPv4.String())
	assert.Equal(t, "2a01::1", nameservers[0].IPv6.String())
}

func TestRepository_UpdateNameservers(t *testing.T) {
	expectedRequest := `{"nameservers":[{"hostname":"ns0.transip.nl","ipv4":"127.0.0.1","ipv6":"2a01::1"}]}`
	server := testutil.MockServer{T: t, ExpectedMethod: "PUT", ExpectedURL: "/domains/example.com/nameservers", StatusCode: 204, ExpectedRequest: expectedRequest}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := Repository{Client: *client}

	nameservers := []Nameserver{{
		Hostname: "ns0.transip.nl",
		IPv4:     net.ParseIP("127.0.0.1"),
		IPv6:     net.ParseIP("2a01::1"),
	}}
	err := repo.UpdateNameservers("example.com", nameservers)
	require.NoError(t, err)
}

func TestRepository_GetDomainAction(t *testing.T) {
	apiResponse := `{"action":{"name":"changeNameservers","message":"success","hasFailed":false}}`
	server := testutil.MockServer{T: t, ExpectedMethod: "GET", ExpectedURL: "/domains/example.com/actions", StatusCode: 200, Response: apiResponse}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := Repository{Client: *client}

	action, err := repo.GetDomainAction("example.com")
	require.NoError(t, err)
	assert.Equal(t, "changeNameservers", action.Name)
	assert.Equal(t, "success", action.Message)
	assert.Equal(t, false, action.HasFailed)

}

func TestRepository_RetryDomainAction(t *testing.T) {
	expectedRequest := `{"authCode":"test","dnsEntries":[{"name":"www","expire":86400,"type":"A","content":"127.0.0.1"}],"nameservers":[{"hostname":"ns0.transip.nl","ipv4":"127.0.0.1","ipv6":"2a01::1"}],"contacts":[{"type":"registrant","firstName":"John","lastName":"Doe","companyName":"Example B.V.","companyKvk":"83057825","companyType":"BV","street":"Easy street","number":"12","postalCode":"1337 XD","city":"Leiden","phoneNumber":"+31 715241919","faxNumber":"+31 715241919","email":"example@example.com","country":"nl"}]}`
	server := testutil.MockServer{T: t, ExpectedMethod: "PATCH", ExpectedURL: "/domains/example.com/actions", StatusCode: 204, ExpectedRequest: expectedRequest}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := Repository{Client: *client}

	contacts := []WhoisContact{
		{
			Type:        "registrant",
			FirstName:   "John",
			LastName:    "Doe",
			CompanyName: "Example B.V.",
			CompanyKvk:  "83057825",
			CompanyType: "BV",
			Street:      "Easy street",
			Number:      "12",
			PostalCode:  "1337 XD",
			City:        "Leiden",
			PhoneNumber: "+31 715241919",
			FaxNumber:   "+31 715241919",
			Email:       "example@example.com",
			Country:     "nl",
		},
	}

	nameservers := []Nameserver{{
		Hostname: "ns0.transip.nl",
		IPv4:     net.ParseIP("127.0.0.1"),
		IPv6:     net.ParseIP("2a01::1"),
	}}

	dnsEntries := []DNSEntry{{Content: "127.0.0.1", Expire: 86400, Name: "www", Type: "A"}}

	err := repo.RetryDomainAction("example.com", "test", dnsEntries, nameservers, contacts)
	require.NoError(t, err)
}

func TestRepository_CancelDomainAction(t *testing.T) {
	server := testutil.MockServer{T: t, ExpectedMethod: "DELETE", ExpectedURL: "/domains/example.com/actions", StatusCode: 204}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := Repository{Client: *client}

	err := repo.CancelDomainAction("example.com")
	require.NoError(t, err)
}

func TestRepository_GetSSLCertificates(t *testing.T) {
	apiResponse := `{"certificates":[{"certificateId":12358,"commonName":"example.com","expirationDate":"2019-10-24 12:59:59","status":"active"}]}`
	server := testutil.MockServer{T: t, ExpectedMethod: "GET", ExpectedURL: "/domains/example.com/ssl", StatusCode: 200, Response: apiResponse}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := Repository{Client: *client}

	certificates, err := repo.GetSSLCertificates("example.com")
	require.NoError(t, err)
	require.Equal(t, 1, len(certificates))
	assert.Equal(t, 12358, certificates[0].CertificateID)
	assert.Equal(t, "example.com", certificates[0].CommonName)
	assert.Equal(t, "2019-10-24 12:59:59", certificates[0].ExpirationDate)
	assert.Equal(t, "active", certificates[0].Status)
}

func TestRepository_GetSSLCertificateByID(t *testing.T) {
	apiResponse := `{"certificate":{"certificateId":12358,"commonName":"example.com","expirationDate":"2019-10-24 12:59:59","status":"active"}}`
	server := testutil.MockServer{T: t, ExpectedMethod: "GET", ExpectedURL: "/domains/example.com/ssl/12358", StatusCode: 200, Response: apiResponse}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := Repository{Client: *client}

	certificates, err := repo.GetSSLCertificateByID("example.com", 12358)
	require.NoError(t, err)
	assert.Equal(t, 12358, certificates.CertificateID)
	assert.Equal(t, "example.com", certificates.CommonName)
	assert.Equal(t, "2019-10-24 12:59:59", certificates.ExpirationDate)
	assert.Equal(t, "active", certificates.Status)
}

func TestRepository_GetWHOIS(t *testing.T) {
	apiResponse := `{"whois":"test123"}`
	server := testutil.MockServer{T: t, ExpectedMethod: "GET", ExpectedURL: "/domains/example.com/whois", StatusCode: 200, Response: apiResponse}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := Repository{Client: *client}

	whoisInfo, err := repo.GetWHOIS("example.com")
	require.NoError(t, err)
	assert.Equal(t, "test123", whoisInfo)
}

func TestRepository_OrderWhitelabel(t *testing.T) {
	server := testutil.MockServer{T: t, ExpectedMethod: "POST", ExpectedURL: "/whitelabel", StatusCode: 201}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := Repository{Client: *client}

	err := repo.OrderWhitelabel()
	require.NoError(t, err)
}

func TestRepository_GetAvailability(t *testing.T) {
	apiResponse := `{"availability":{"domainName":"example.com","status":"free","actions":["register"]}}`
	server := testutil.MockServer{T: t, ExpectedMethod: "GET", ExpectedURL: "/domain-availability/example.com", StatusCode: 200, Response: apiResponse}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := Repository{Client: *client}

	availability, err := repo.GetAvailability("example.com")
	require.NoError(t, err)
	assert.Equal(t, "example.com", availability.DomainName)
	assert.EqualValues(t, "free", availability.Status)
	assert.Equal(t, []PerformAction{"register"}, availability.Actions)
}

func TestRepository_GetAvailabilityForMultipleDomains(t *testing.T) {
	apiResponse := `{"availability":[{"domainName":"example.com","status":"free","actions":["register"]}]}`
	expectedRequest := `{"domainNames":["example.com","example.nl"]}`
	server := testutil.MockServer{T: t, ExpectedMethod: "GET", ExpectedURL: "/domain-availability", StatusCode: 200, ExpectedRequest: expectedRequest, Response: apiResponse}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := Repository{Client: *client}

	availabilityList, err := repo.GetAvailabilityForMultipleDomains([]string{"example.com", "example.nl"})
	require.NoError(t, err)
	require.Equal(t, 1, len(availabilityList))
	assert.Equal(t, "example.com", availabilityList[0].DomainName)
	assert.EqualValues(t, "free", availabilityList[0].Status)
	assert.Equal(t, []PerformAction{"register"}, availabilityList[0].Actions)
}

func TestRepository_GetTLDs(t *testing.T) {
	apiResponse := `{"tlds":[{"name":".nl","price":399,"recurringPrice":749,"capabilities":["canRegister"],"minLength":2,"maxLength":63,"registrationPeriodLength":12,"cancelTimeFrame":1}]}`
	server := testutil.MockServer{T: t, ExpectedMethod: "GET", ExpectedURL: "/tlds", StatusCode: 200, Response: apiResponse}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := Repository{Client: *client}

	tlds, err := repo.GetTLDs()
	require.NoError(t, err)
	require.Equal(t, 1, len(tlds))
	assert.Equal(t, ".nl", tlds[0].Name)
	assert.Equal(t, 399, tlds[0].Price)
	assert.Equal(t, 749, tlds[0].RecurringPrice)
	assert.Equal(t, 2, tlds[0].MinLength)
	assert.Equal(t, 63, tlds[0].MaxLength)
	assert.Equal(t, 12, tlds[0].RegistrationPeriodLength)
	assert.Equal(t, 1, tlds[0].CancelTimeFrame)

	assert.Equal(t, []string{"canRegister"}, tlds[0].Capabilities)
}

func TestRepository_GetTldInfo(t *testing.T) {
	apiResponse := `{"tld":{"name":".nl","price":399,"recurringPrice":749,"capabilities":["canRegister"],"minLength":2,"maxLength":63,"registrationPeriodLength":12,"cancelTimeFrame":1}}`
	server := testutil.MockServer{T: t, ExpectedMethod: "GET", ExpectedURL: "/tlds/.nl", StatusCode: 200, Response: apiResponse}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := Repository{Client: *client}

	tld, err := repo.GetTLDByTLD(".nl")
	require.NoError(t, err)
	assert.Equal(t, ".nl", tld.Name)
	assert.Equal(t, 399, tld.Price)
	assert.Equal(t, 749, tld.RecurringPrice)
	assert.Equal(t, 2, tld.MinLength)
	assert.Equal(t, 63, tld.MaxLength)
	assert.Equal(t, 12, tld.RegistrationPeriodLength)
	assert.Equal(t, 1, tld.CancelTimeFrame)

	assert.Equal(t, []string{"canRegister"}, tld.Capabilities)
}
