#!/bin/bash

docker tag trusch/pkid:latest trusch/pkid:$(git describe)
docker push trusch/pkid:latest
docker push trusch/pkid:$(git describe)

exit $?
