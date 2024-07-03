FROM devopsfaith/krakend:2.6.1
EXPOSE APP_PORT
COPY krakend.json /etc/krakend/krakend.json
ENTRYPOINT ["/usr/bin/krakend", "run", "-d", "-c", "/etc/krakend/krakend.json", "-p", "APP_PORT"]