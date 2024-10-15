# telegram-budjet-bot

docker build -t firepand4/fortress:budget-bot .
docker push firepand4/fortress:budget-bot

docker run -dit --rm --name budget-bot-app firepand4/fortress:budget-bot

# 32 bit
docker buildx build --platform linux/arm/v7 -t firepand4/fortress:budget-bot .