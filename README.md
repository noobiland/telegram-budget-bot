# telegram-budjet-bot

docker build -t firepand4/fortress:budget-bot .
docker push firepand4/fortress:budget-bot

docker run -dit --rm --name budget-bot-app -v "$(Get-Location)/output:/app/output" firepand4/fortress:budget-bot
docker run -dit --name budget-bot-app -v "$(Get-Location)/output:/output" firepand4/fortress:budget-bot


# 32 bit
docker buildx build --platform linux/arm/v7 -t firepand4/fortress:budget-bot .
docker push firepand4/fortress:budget-bot

docker run -dit --name budget-bot-app -v "/home/pi/temp/telegram:/output" firepand4/fortress:budget-bot

-v $(pwd):/app

/usr/src/app


# TODO
- add logging
- reorganize code
- add confirmation step
- add db versioning (migrations)
    - https://github.com/golang-migrate/migrate/blob/master/GETTING_STARTED.md