# webservice-proxy

This project is still just a dream. It's under development, but time well tell how far I get!

## Roadmap

* Create proxy that handles POSTs
* Stats page that shows all URLs seen
* Clustering of servers (ring?)
* Stats page stats and config options
    * min,max,mean,median successful request times in past 10 minutes
    * timed-out requests in past 10 minutes
    * error requests in past 10 minutes
    * successful requests in past 10 minutes
    * warn-if-more-than-X-slow-connections: 200 default
    * fail-if-more-than-X-slow-connections: 0 default (never fail)
    * max-timeout: 45 seconds default