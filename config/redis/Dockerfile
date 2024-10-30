FROM redis:7.2-alpine

WORKDIR /redis

COPY /config/redis.conf /usr/local/etc/redis/redis.conf

COPY /scripts/redis.sh ./

RUN chmod +x redis.sh