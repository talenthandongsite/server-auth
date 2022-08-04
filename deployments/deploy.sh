NAME=$(git config --get remote.origin.url | sed 's/.*\/\([^.]*\).*/\1/')
TAG=latest

SERVER_ENV=`echo ${SERVER_ENV} | xargs`
ENV_LINE=${SERVER_ENV//" "/" -e "}

docker build -f deployments/dockerfile -t ${NAME}:${TAG} .

docker ps -f name=${NAME} -q | xargs -r docker container stop
docker container ls -a -fname=${NAME} -q | xargs -r docker container rm -f
docker images --no-trunc --all --quiet --filter='dangling=true' | xargs -r docker rmi
docker run -e ${ENV_LINE} --network ${DOCKER_NETWORK} -p 8080:8080 -d --name ${NAME} ${NAME}:${TAG}
