#!/usr/bin/env sh

docker build -t neo9sas/mongodb-backup .

if [ $TRAVIS_BRANCH == 'master' ]; then
  echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_LOGIN" --password-stdin
  docker push neo9sas/mongodb-backup:latest
fi


if ! [ -z "$TRAVIS_TAG" ]; then
    echo "Release $TRAVIS_TAG"
    echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_LOGIN" --password-stdin
    docker tag neo9sas/mongodb-backup:latest neo9sas/mongodb-backup:$TRAVIS_TAG
    docker push neo9sas/mongodb-backup:$TRAVIS_TAG
fi

