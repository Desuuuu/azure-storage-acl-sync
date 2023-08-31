package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v4"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
)

// GetIPRulesForServices returns a list of allow IP rules for the specified services.
func GetIPRulesForServices(ctx context.Context, credential azcore.TokenCredential, subscriptionID string, location string, serviceTags ...string) ([]*armstorage.IPRule, error) {
	if len(serviceTags) < 1 {
		return nil, nil
	}

	factory, err := armnetwork.NewClientFactory(subscriptionID, credential, nil)
	if err != nil {
		return nil, err
	}

	client := factory.NewServiceTagsClient()

	ids := make(map[string]struct{}, len(serviceTags))
	for _, serviceTag := range serviceTags {
		ids[serviceTag] = struct{}{}
	}

	res, err := client.List(ctx, location, nil)
	if err != nil {
		return nil, err
	}

	var rules []*armstorage.IPRule

	allow := "Allow"

	for _, serviceTag := range res.Values {
		if serviceTag == nil || serviceTag.ID == nil || serviceTag.Properties == nil {
			continue
		}

		if _, ok := ids[*serviceTag.ID]; !ok {
			continue
		}

		for _, addressPrefix := range serviceTag.Properties.AddressPrefixes {
			if addressPrefix == nil {
				continue
			}

			rules = append(rules, &armstorage.IPRule{
				IPAddressOrRange: addressPrefix,
				Action:           &allow,
			})
		}
	}

	return rules, nil
}
