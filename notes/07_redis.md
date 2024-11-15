## Detailed Parameter Explanation (in docker-compose.yml)

- **./redis.sh &**: Executes a custom shell script in the background before starting the Redis server. This script is intended for any preparatory steps such as loading Redis modules or performing initial configurations.
    
- **redis-server /usr/local/etc/redis/redis.conf**: Initiates the Redis server using a configuration file that includes customized settings tailored for the project.
    
- **--loadmodule /usr/lib/redis/modules/redisbloom.so**: This command specifies the path to the RedisBloom module, instructing Redis to load it upon startup.

- **--requirepass**: Protects the Redis server with a password, preventing unauthorized access. The password is managed via an environment variable to keep sensitive data out of the version control.
    
- **--loglevel warning**: Limits log output to warnings and above, which helps in reducing the volume of log data, especially useful in production environments.
    
- **--bind 0.0.0.0**: Configures Redis to accept connections on all network interfaces, facilitating connectivity within the Docker network and, if necessary, from the host.
    
- **--appendonly no**: Disables the append-only mode to prevent Redis from logging every write operation to disk. This can be beneficial for reducing disk IO in environments where persistence is not a priority.
    
- **--save 300 1**: Instructs Redis to save the dataset to disk every 300 seconds if at least one change has been made. This setting provides a basic level of data persistence.
    
- **--dir /data**: Designates `/data` as the directory where Redis will store its persistence files. This directory should be mounted as a Docker volume to ensure data persists across container restarts.

## Verify RedisBloom Installation

After starting the container, connect to Redis and verify that RedisBloom is installed and working:

```bash
docker exec -it iot-parking-gateway_redis redis-cli -a your_redis_password
```
Then, run the following command in the Redis CLI to check if the Bloom filter commands are available

```bash
BF.RESERVE test_filter 0.01 1000
```

