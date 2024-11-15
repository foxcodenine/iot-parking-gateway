FROM redislabs/rebloom:latest  

WORKDIR /redis

COPY /config/redis.conf /usr/local/etc/redis/redis.conf

COPY /scripts/redis.sh ./

RUN chmod +x redis.sh