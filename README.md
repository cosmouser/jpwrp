# jpwrp
Reverse proxy server for forwarding Jamf Pro webhooks

## disclaimer
**CONSIDER USING NGINX INSTEAD OF THIS**. It's got way more power and flexibility.

## background
In order to understand what jpwrp does it is necessary to first understand the problem it was intended to solve. The Jamf Pro server can be configured to send webhooks to a webhook server to report on specific events, such as a computer's inventory report being completed. The Jamf Pro server will often be in a secured network with significantly monitored traffic. It makes sense to limit the number of hosts and ports that the server communicates with. It also makes sense to give the JSS a secure platform (TLS) for transporting all of the requests that it sends to your webhook servers. Using jpwrp you can point each of your webhooks to the same host and port and specify unique web servers using a URI.

## setup
1. Download and compile
2. Put the webhook paths and ports to your servers in the config.toml file
3. Start it up
