# abuseipdb-batch
Batch submission towards the Abuse IP DB API.

## How to use this script
- Either compile the code using `go build` or run the code using `go run`
- Provide a file with FQDN's or IP addresses (single or CIDR annotated). Each line is interpreted as individual entry that is sent to the API.
- `go run main.go sample.txt APIKEY` and replace APIKEY with your API key obtained from the Abuse IP DB account settings.
- The results are saved in `results` directory, each line item will be saved in its own individual file.
