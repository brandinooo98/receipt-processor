# Receipt Processor

This is a simple webservice that takes in a receipt from a JSON and gives the receipt a score based on the given criteria:

- One point for every alphanumeric character in the retailer name.
- 50 points if the total is a round dollar amount with no cents.
- 25 points if the total is a multiple of `0.25`.
- 5 points for every two items on the receipt.
- If the trimmed length of the item description is a multiple of 3, multiply the price by `0.2` and round up to the nearest integer. The result is the number of points earned.
- 6 points if the day in the purchase date is odd.
- 10 points if the time of purchase is after 2:00pm and before 4:00pm.

The webservice consists of two endpoints:

- `POST /receipts/process` - This endpoint takes in a receipt from a JSON file, calculates it's score, and stores the score to the ID given in the response.
- `GET /receipts/:id/points` - This endpoint takes in the ID of a receipt and returns the score of the receipt with that ID.

# Running Instructions

To run this webservice, you will need to have Go installed. You can download it [here](https://golang.org/dl/). Once you have Go installed, you can run the following commands to start the webservice after navigating to the root directory of this project:

```bash
go run .
```

To make a request to the webservice, you can either use a tool like Postman or you can use the following `curl` command to make a request to the `/receipts/process` endpoint:

```bash
curl -X POST -H "Content-Type: application/json" -d @examples/morning-receipt.json http://localhost:8080/receipts/process
```

This will send the `morning-receipt.json` file to the webservice and return the score of the receipt. You can see the score of the receipt by making a request to the `/receipts/:id/points` endpoint. The id of the receipt will be returned in the response of the initial request. To query the score of the receipt, you can use the following `curl` command:

```bash
curl -X GET http://localhost:8080/receipts/0/points
```

This will return the score of the receipt with the id `0`, which should be the id of the first receipt that was processed.
