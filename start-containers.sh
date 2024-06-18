docker pull scylladb/scylla

docker container stop nodex nodey nodez
docker container rm nodex nodey nodez

docker run --name nodex -d -p 9002:9042 scylladb/scylla:latest --reactor-backend=epoll --overprovisioned 1 --smp 1
sleep 60
docker run --name nodey -d -p 9003:9042 scylladb/scylla:latest --reactor-backend=epoll --seeds="$(docker inspect --format='{{ .NetworkSettings.IPAddress }}' nodex)" --overprovisioned 1 --smp 1
sleep 60
docker run --name nodez -d -p 9004:9042 scylladb/scylla:latest --reactor-backend=epoll --seeds="$(docker inspect --format='{{ .NetworkSettings.IPAddress }}' nodex)" --overprovisioned 1 --smp 1
sleep 60
docker container exec -it nodex nodetool status
echo "Ready!"
