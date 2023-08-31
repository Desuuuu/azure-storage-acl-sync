# azure-storage-acl-sync

[![Tag Badge]][Tag] [![Go Version Badge]][Go Version] [![Go Report Card Badge]][Go Report Card]

Synchronize Azure storage account IP ACL with Azure service IPs.

## Authentication

Azure authentication is handled by the [azidentity](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity) package with `DefaultAzureCredential`. The easiest way to authenticate is using the following environment variables:

<details>
  <summary><b>Service principal with secret</b></summary>

  `AZURE_TENANT_ID`, `AZURE_CLIENT_ID` and `AZURE_CLIENT_SECRET`.
</details>

<details>
  <summary><b>Service principal with certificate</b></summary>

  `AZURE_TENANT_ID`, `AZURE_CLIENT_ID`, `AZURE_CLIENT_CERTIFICATE_PATH` and `AZURE_CLIENT_CERTIFICATE_PASSWORD`.
</details>

<details>
  <summary><b>Username and password</b></summary>

  `AZURE_CLIENT_ID`, `AZURE_USERNAME` and `AZURE_PASSWORD`.
</details>

## Permissions

* `Microsoft.Network/locations/*/serviceTags/read` action on the subscription to retrieve the service IPs.
* Writing properties on the configured storage account to update its IP ACL.

<details>
  <summary><b>Custom role for reading service tags</b></summary>

  ```json
{
    "Name": "Service Tag Reader",
    "IsCustom": true,
    "Description": "List service tags and their respective IPs.",
    "Actions": [
        "Microsoft.Network/locations/*/serviceTags/read"
    ],
    "NotActions": [],
    "DataActions": [],
    "NotDataActions": [],
    "AssignableScopes": [
        "/subscriptions/{subscriptionId}"
    ]
}
  ```
</details>

## Options

| Flag                  | Environment variable    | Default                                  | Description                                                                                                                                  |
|:----------------------|:------------------------|:-----------------------------------------|:----------------------------------------------------------------|
| `--subscription-id`   | `AZURE_SUBSCRIPTION_ID` | -                                        | Azure subscription ID.                                          |
| `--services`          | `AZURE_SERVICES`        | `AzureFrontDoor.Backend`                 | Azure [services](https://learn.microsoft.com/en-us/azure/virtual-network/service-tags-overview#available-service-tags) to retrieve IPs from. |
| `--location`          | `AZURE_LOCATION`        | `westus`                                 | Azure location to retrieve IPs for.                             |
| `--resource-group`    | `AZURE_RESOURCE_GROUP`  | -                                        | Storage account resource group.                                 |
| `--storage-account`   | `AZURE_STORAGE_ACCOUNT` | -                                        | Storage account name.                                           |
| `--extra-allow-rules` | `EXTRA_ALLOW_RULES`     | `168.63.129.16` <br /> `169.254.169.254` | Additional allow IP rules.                                      |
| `--extra-deny-rules`  | `EXTRA_DENY_RULES`      | -                                        | Additional deny IP rules.                                       |
| `--dry-run`           | `DRY_RUN`               | `false`                                  | Only print the IP rules that would be applied.                  |

The two IP addresses allowed by default are documented [here](https://learn.microsoft.com/en-us/azure/virtual-network/network-security-groups-overview#azure-platform-considerations).

[Tag]: https://github.com/Desuuuu/azure-storage-acl-sync/tags
[Tag Badge]: https://img.shields.io/github/v/tag/Desuuuu/azure-storage-acl-sync?sort=semver
[Go Version]: /go.mod
[Go Version Badge]: https://img.shields.io/github/go-mod/go-version/Desuuuu/azure-storage-acl-sync
[Go Report Card]: https://goreportcard.com/report/github.com/Desuuuu/azure-storage-acl-sync
[Go Report Card Badge]: https://goreportcard.com/badge/github.com/Desuuuu/azure-storage-acl-sync
