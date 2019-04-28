set -ex

# SET THE FOLLOWING VARIABLES
# docker hub username
USERNAME=davidxchen
# image name
IMAGE=cartrackerapi

docker build -f Dockerfile -t $USERNAME/$IMAGE:latest .