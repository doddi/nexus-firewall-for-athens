# nexus-firewall-for-athens

A Golang application to check for any vulnerable components when ingesting into Athens proxy server by using either Sonatype OssIndex or Sonatype Nexus IQ Server (soon!) 

## Introduction
To protect your applications from any already known vulnerabilities Firewall For Athens checks for any vulnerabilities at the time of ingestion.
If a package is known to be vulnerable the package is prevented from being fetched from upstream. 

## Usage
To make use of this application, your Athens server needs to be configured to point to this running server.
To do so, provide a url for ValidatorHook parameter in your configuration file or override using the environment variable `ATHENS_PROXY_VALIDATOR`. For example, `ATHENS_PROXY_VALIDATOR=http://localhost:8080`

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

If you would like the builds to still pass (respond with http OK) then you can force the build to pass but Athens will still log the warning message
- `--failBuild`: defaults to true to block the fetching of the component

## Vulnerability Source
### Sonatype OssIndex
This is the default behaviour, or can specified using `--using=ossindex`

If *any* vulnerabilities are known for the package then the package is blocked from being fetched

### Sonatype Nexus IQ Server
Specify `--using=nexusiq`

To use basic authentication supply `--username` and `--password`
An applicationId is required to evaluate against policies, use `--appid` to supply the correct Application Id.
To determine the application id follow the instructions supplied here https://help.sonatype.com/iqserver/automating/rest-apis/application-rest-apis---v2#ApplicationRESTAPIs-v2-Step5-UpdateApplicationInformation and use the `id` field.

The results using Nexus Iq will provide policy violations, security vulnerabilities and license information.

## To Do
Any contributions would be great!
