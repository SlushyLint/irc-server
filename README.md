one fucing thing you have to run is this:
zsh```
openssl req -x509 -newkey rsa:2048 \                                                                                                                          N
    -keyout server.key \
    -out server.crt \
    -days 365 \
    -nodes \
    -subj "/CN=localhost"
```
only then can you run the file
