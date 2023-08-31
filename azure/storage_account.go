package azure

import (
	"context"
	"net"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
)

// UpdateStorageAccountIPRules updates the network IP rules for the specified storage account.
func UpdateStorageAccountIPRules(ctx context.Context, credential azcore.TokenCredential, subscriptionID string, resourceGroup string, storageAccount string, rules []*armstorage.IPRule) error {
	factory, err := armstorage.NewClientFactory(subscriptionID, credential, nil)
	if err != nil {
		return err
	}

	client := factory.NewAccountsClient()

	_, err = client.Update(ctx, resourceGroup, storageAccount, armstorage.AccountUpdateParameters{
		Properties: &armstorage.AccountPropertiesUpdateParameters{
			NetworkRuleSet: &armstorage.NetworkRuleSet{
				IPRules: rules,
			},
		},
	}, nil)

	return err
}

// FilterRules transforms the specified rules into a list that is acceptable by
// Azure. See https://learn.microsoft.com/en-us/azure/storage/common/storage-network-security?tabs=azure-portal#restrictions-for-ip-network-rules.
func FilterRules(rules []*armstorage.IPRule) []*armstorage.IPRule {
	res := make([]*armstorage.IPRule, 0, len(rules))

	for _, rule := range rules {
		value := strings.TrimSpace(*rule.IPAddressOrRange)
		if value == "" {
			continue
		}

		ip := net.ParseIP(value)
		if ip != nil {
			ip = ip.To4()
			if ip == nil {
				continue
			}

			if ip.IsUnspecified() || ip.Equal(net.IPv4bcast) || ip.IsLoopback() || ip.IsMulticast() {
				continue
			}

			res = append(res, &armstorage.IPRule{
				IPAddressOrRange: to.Ptr(ip.String()),
				Action:           rule.Action,
			})

			continue
		}

		_, ipNet, err := net.ParseCIDR(value)
		if err != nil {
			continue
		}

		mask, total := ipNet.Mask.Size()
		if total != 32 {
			continue
		}

		if mask <= 30 {
			res = append(res, &armstorage.IPRule{
				IPAddressOrRange: to.Ptr(ipNet.String()),
				Action:           rule.Action,
			})

			continue
		}

		for ip = ipNet.IP.Mask(ipNet.Mask); ipNet.Contains(ip); nextIP(ip) {
			res = append(res, &armstorage.IPRule{
				IPAddressOrRange: to.Ptr(ip.String()),
				Action:           rule.Action,
			})
		}
	}

	return res
}

func nextIP(ip net.IP) {
	for i := len(ip) - 1; i >= 0; i-- {
		ip[i]++

		if ip[i] > 0 {
			break
		}
	}
}
