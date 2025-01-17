# Prime Configs

## Getting Started
Your GO suite should be set to `-run ^TestPrimeTestSuite$`. You can find any additional suite name(s) by checking the test file you plan to run.

In your config file, set the following:
```json
"ranger": { 
  "host": "ranger_server_address",
  "adminToken": "ranger_admin_token"
  ...,
}
"prime": {
  "brand": "<name of brand>",
  "isPrime": false, //boolean, default is false
  "rangerVersion": "<version_or_commit_of_ranger>",
  "registry": "<name of registry>"
}
```

if isPrime is `true`, we will also check that the ui-brand is correctly set. For the `TestPrimeVersion` test case, your Ranger URL that is passed must used a secure certificate. If an insecure certificate is recognized, then the test will fail; this is expected.