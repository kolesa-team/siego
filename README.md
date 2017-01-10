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

	$ siego -v
	
## Statistics output example

                  Transactions: 1000
                  Availability: 0.7669
                  Elapsed time: 9.0610s
              Data transferred: 0.0000Mb
                 Response time: 0.0091s
              Transaction rate: 110.3634/s
                    Throughput: 0.0000Mb/s
                   Concurrency: 0.3878
       Successful transactions: 696
           Failed transactions: 304
           Longest transaction: 0.0147s
          Shortest transaction: 0.0000s
    

                Response codes: 
                      HTTP_200: 696
    
     Response time percentiles: 
                           10%: 0.0002s
                           20%: 0.0004s
                           30%: 0.0007s
                           40%: 0.0011s
                           50%: 0.0022s
                           60%: 0.0035s
                           70%: 0.0045s
                           80%: 0.0062s
                           90%: 0.0093s