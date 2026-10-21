package networks

import "testing"

func TestNetworkUpdateOpts_ToNetworkUpdateMapIncludesPrivateDNSDomain(t *testing.T) {
	const domain = "test.domain."

	body, err := (NetworkUpdateOpts{PrivateDNSDomain: domain}).ToNetworkUpdateMap()
	if err != nil {
		t.Fatalf("ToNetworkUpdateMap() error = %v", err)
	}

	network, ok := body["network"].(map[string]interface{})

	if !ok {
		t.Fatalf("body[\"network\"] has type %T, want map[string]interface{}", body["network"])
	}

	if got := network["private_dns_domain"]; got != domain {
		t.Errorf("private_dns_domain = %v, want %q", got, domain)
	}
}
