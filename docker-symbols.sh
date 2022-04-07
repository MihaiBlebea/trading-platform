#!/bin/bash

docker exec -it app /bin/bash -c "./wait-for-db.sh ./trading-platform symbols ./assets/test.csv"