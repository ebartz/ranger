# Provisioning Configs

## Getting Started
Your GO suite should be set to `-run ^Test<enter_pkg_name_here>ProvisioningTestSuite$`. You can find the correct suite name in the below README links, or by checking the test file you plan to run.
In your config file, set the following:
```json
"ranger": { 
  "host": "ranger_server_address",
  "adminToken": "ranger_admin_token",
  "userToken": "your_ranger_user_token",
  "insecure": true/optional,
  "cleanup": false/optional,
}
```

From there, your config should contain the tests you want to run (provisioningInput), tokens and configuration for the provider(s) you will use, and any additional tests that you may want to run. Please use one of the following links to continue adding to your config for provisioning tests:

1. [RKE1 Provisioning](rke1/README.md)
2. [RKE2 Provisioning](rke2/README.md)
3. [Hosted Provider Provisioning](hosted/README.md)
