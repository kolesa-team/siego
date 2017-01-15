Siego
=====

Siego is regression test and benchmark utility. Golang port of [Siege](https://github.com/JoeDog/siege) with extended stats output.

## Requirements

 * Go 1.6+
 * make
 * build-essential
 * fpm (`apt-get install ruby-dev && gem install fpm`) to generate *.deb

## Installation

    # make dep-install
	# make
	# make install

Optionally to create Debian package

	# make deb
	
## Usage

    $ siego -h

    GLOBAL OPTIONS:
       --concurrent value, -c value    CONCURRENT users. (default: 10)
       --delay value, -d value         Time DELAY, random delay before each requst between 1 and NUM. (NOT COUNTED IN STATS) (default: 1)
       --reps value, -r value          REPS, number of times to run the test. (default: 0)
       --url value, -u value           URL to test.
       --file value, -f value          FILE, select a specific URLS FILE.
       --log value, -l value           LOG to FILE. (default: "/var/siege.log")
       --time value, -t value          TIMED testing where "m" is modifier s, m, or h. Ex: --time=1h, one hour test.
       --header value, -H value        Add a header to request (can be many)
       --user-agent value, -A value    Sets User-Agent in request
       --content-type value, -T value  Sets Content-Type in request
       --get, -g                       Use GET method.
       --post, -p                      Use POST method.
       --internet, -i                  INTERNET user simulation, hits URLs randomly.
       --benchmark, -b                 BENCHMARK: no delays between requests.
       --xml, -x                       Use XML output.
       --timeout value                 Request timeout in seconds. (default: 1)
       --help, -h                      show help
       --version, -v                   print the version

	
## Statistics output example

                  Transactions: 14
                  Availability: 82.35%
                  Elapsed time: 1.5609s
              Data transferred: 0.0000Mb
                 Response time: 0.1115s
              Transaction rate: 8.9691/s
                    Throughput: 0.0000Mb/s
                   Concurrency: 0.0141
       Successful transactions: 11
           Failed transactions: 3
           Longest transaction: 0.0038s
          Shortest transaction: 0.0001s


                Response codes: 
                      HTTP_200: 11

     Response time percentiles: 
                           10%: 0.0003s
                           20%: 0.0006s
                           30%: 0.0007s
                           40%: 0.0008s
                           50%: 0.0010s
                           60%: 0.0016s
                           70%: 0.0021s
                           80%: 0.0033s
                           90%: 0.0035s
