# nexus-firewall-for-athens

A Golang application to check for any vulnerable components when ingesting into Athens proxy server

This application can run as:
 - AWS lambda
 - Local server
 - Standalone cli
 
## Options
### AWS Lambda
To run as an AWS lambda no parameters are needed

### Localserver
Specify `--run=server` to run as a local server
    
Other configurable parameters when running as a local server are
- port: to specify the port number to run on (default 8080)
