# SKYPE BOT

## Example
See working example in [example_test.go](example_test.go).
Just set before the following env variables:
- SKYPE_APP_ID `// bot id`
- SKYPE_APP_SECRET `// bot token`
- PORT `// web server port (used for hook)`

## How to find bot app credentials?
- open Azure Portal https://portal.azure.com
- go "App Registrations" https://portal.azure.com/#blade/Microsoft_AAD_RegisteredApps/ApplicationsListBlade
- find your bot in the list (my bot app was located inside "Applications from personal account") and open it
- see your SKYPE_APP_ID in "Application (client) ID"
- see your SKYPE_APP_SECRET in the section "Client secrets" of "Manage -> Certificates & secrets"
