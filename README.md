# cartracker-api
go-lang


# build docker image
docker build -t davidxchen/cartrackerapi -f Dockerfile .

# run docker image
docker run --rm -it -p 8100:8000 -v /opt/docker_data_vol/cartrackerapi:/app/settings davidxchen/cartrackerapi