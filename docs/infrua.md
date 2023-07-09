# Bandwidth and Instance Requirements Estimation

Assuming we have 300,000 requests per day on our API endpoint, we need to calculate the requirements in terms of bandwidth and server instances. The API endpoint in question is the "Add Products" endpoint.

To calculate the number of requests per second, we can divide the total number of requests by the total number of seconds in a day:

    300,000 / (3600 * 8) = ~10.42 requests/second

Therefore, we need to account for roughly 11 requests per second.

Given that users try to add 100 products per request, and each product has the following structure:

    message Product {
      int64 chain_id = 1;
      int64 asset_id = 2;
      string barcode = 4;
      string owner = 3;
      string item_name = 5;
      string brand_name = 6;
      string additional_data = 7;
      string location = 8;
      int64 issued_date = 9;
    }

We need to estimate the size of this message:

- int64 types are 8 bytes each.
- string types can be complex. In UTF-8 encoding, one character generally takes 1 byte. However, if you're dealing with internationalized text, a character could take up to 4 bytes. We assume you are using simple ASCII characters, and therefore each character is 1 byte. Furthermore, the gRPC protocol uses length prefixing for strings, which means it includes the length of the string before the string itself, and this length prefix itself can take from 1 to 5 bytes. For our estimation, we use the worst-case scenario of 5 bytes.

Hence, the size of one Product message is:

- 3\*8 bytes (for the int64 fields)
- 6\*(256+5) bytes (for the string fields)

This sums up to 1590 bytes per Product message, or 0.151MB per request (since each request includes 100 Products).

Since we estimated about 11 requests per second, we therefore have a bandwidth requirement of about 1.661MB/s.

In terms of AWS EC2 instances, this requirement would be based on various factors such as the instance type, the specific workload characteristics, and the performance of the application itself. It's recommended to perform performance testing to determine the optimal number of instances required for your specific workload.


First, let's calculate the number of instances required to meet the bandwidth needs:

Given the calculated bandwidth requirement of approximately 0.604MB/s and the free tier bandwidth limit of 100GB/month, the number of instances would be determined by the total bandwidth consumption in a month.

First, let's convert 1.661MB/s to GB/month:

1.661MB/s * 60s/minute * 60min/hour * 24hour/day * 30day/month = ~4202GB/month

As this exceeds the 100GB/month free tier limit, the excess data transfer will be chargeable.

Note: Data transfer in (all data transferred into an instance from the Internet) is free, while Data transfer out has a cost associated with it, especially after the free tier limit.

Amazon EC2 in Ireland, the pricing for data transfer out after the first 1 GB / month was as follows:

`a1.medium	$0.0288	1(vcpu)	2 GiB(memory) EBS Only	Up to 10 Gigabit`

The excess data of ~4202GB -> 5TB will be chargeable as per these rates. Please verify the current prices, as they may have changed.

$0.0288 * 24 * 31 = $0.0288 * (24 * 31) = $0.0288 * 744  = $21.4272 / month
