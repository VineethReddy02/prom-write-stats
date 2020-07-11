# Prometheus Remote Write Stats

Long term storage in prometheus needs an external metric process system integration like cortex, thanos, etc... This tiny application aims to provide stats on metrics by consuming metrics from prometheus remote write api.

All you need is point the remote write api to this application along with your external storage system. Consider you have multiple prometheus series scraping different targets this application can provide overall series information with last timestamp the sample has been received. 

Running this application:

```
$ cd bin 
$ ./prom-write-stats
```

Now accessing the URL http://localhost:2233/stats will provide the below information:

![Screenshot from 2020-07-11 15-07-45](https://user-images.githubusercontent.com/25104868/87221296-886c7800-c388-11ea-9b55-945486110d1d.png)

Feel free to suggest improvements & thoughts.

Cheers! 