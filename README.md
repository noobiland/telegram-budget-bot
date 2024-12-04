# telegram-budget-bot

To run locally: generate dbs with db-handler and put files into output

## Docker commands for debugging purposes
```
docker build -t firepand4/fortress:budget-bot .
docker buildx build --platform linux/arm/v7 -t firepand4/fortress:budget-bot .
docker push firepand4/fortress:budget-bot

docker run -dit --rm --name budget-bot-app -v "$(Get-Location)/output:/app/output" firepand4/fortress:budget-bot
docker run -dit --name budget-bot-app -v "$(Get-Location)/output:/output" firepand4/fortress:budget-bot
```


# TODO
* [x] get users from db
* [x] write data to db, which was prepared by db-handler
* [ ] add confirmation step
* [ ] discard doesn't return default keyboard
* [X] ~~*add total reporting for previous month*~~ [2024-12-03]
* [X] ~~*total reporting for current month*~~ [2024-12-05]
* [ ] total reporting for week
* [ ] add options mechanism for reports
* [ ] add tests
* [ ] log clean up