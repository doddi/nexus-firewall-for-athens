# nexus-firewall-for-athens

A Golang application to check for any vulnerable components when ingesting into Athens proxy server by using either Sonatype OssIndex or Sonatype Nexus IQ Server (soon!) 

## Introduction
To protect your applications from any already known vulnerabilities Firewall For Athens checks for any vulnerabilities at the time of ingestion.
If a package is known to be vulnerable the package is prevented from being fetched from upstream. 

This application can run as:
 - AWS lambda
 - Local server
 
## Environment Options
### AWS Lambda
To run as an AWS lambda no parameters are needed

### Localserver
Specify `--run=server` to run as a local server
    
When running as a local server you can specify a port to run under using:
- `--port`: to specify the port number to run on (default 8080)

## Vulnerability Source
### Sonatype OssIndex
This is the default behaviour, or ca specified using `--using=ossindex`

If *any* vulnerabilities are known for the package then the package is blocked from being fetched

### Sonatype Nexus IQ Server
Coming .....


## To Do
Any contributions would be great!
