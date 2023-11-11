# topedcrawler

Entity: This layer would define the structure of the data that you intend to scrape. For instance, if you're scraping product information, you might have a Product entity with fields like Name, Price, Description, etc.

Repo (Repository): The web scraper itself would most likely reside within the repo directory. The scraper would be responsible for the implementation of interfaces defined in the repo to fetch data from the web. The repository would also handle the conversion of the scraped data into the format expected by your domain entities.

Usecase (Business Logic): The usecase layer would define interfaces that your application's service layer will implement. It could invoke the web scraper through the repository layer interface to obtain data and then perform any necessary business logic on this data, such as filtering, validation, or aggregation.

Handler (Transport or Presentation Layer): The handler layer would not typically contain the scraper itself but would call the usecase layer to initiate scraping based on a user request (like an API call or a cron job initiation). This layer is responsible for dealing with client-side operations, such as handling HTTP requests and responses, and it would translate the results from the usecase layer into a format suitable for the client.

## How to run app
1. Create env.sh in files/etc/env.sh
```sh
export DB_HOST=localhost
export DB_PORT=5432
export DB_USERNAME=postgres
export DB_PASSWORD=postgres
export DB_NAME=postgres
export NUM_WORKERS=2
export NUM_PRODUCTS=100
```

2. Source the environment variable to terminal session
```
source files/etc/env.sh
```

3. Run posgresql in docker (OPTIONAL)

```
docker compose up -d
```

4. Run the application
```
go run cmd/app/main.go  
```

## Extra
Csv file stored in `data.csv`
Known issue, can't be solved because had no time:
- bug `failed to scroll page: Execution context was destroyed, most likely because of a navigation`
- retry mechanism
- Making Url Unique
- automated test/unit test