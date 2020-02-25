FROM redis
RUN apt-get update && apt-get install -y stunnel && apt-get purge && apt-get autoremove
COPY stunnel.conf /etc/stunnel/redis-cli.conf
