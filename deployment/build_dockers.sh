#!/usr/bin/env bash

##########################
# Build watchtopus-server
##########################

# Delete files that were copied to docker dir on previous builds (if any)
rm -rf watchtopus/infra
rm -rf watchtopus/orm
rm -rf watchtopus/server

# Copy source files to docker
cp -r ../infra watchtopus/
cp -r ../orm watchtopus/.
cp -r ../server watchtopus/.

# Override the server's active conf file to be docker-compose's config file
cp -f watchtopus/server/conf/config_docker-compose.json watchtopus/server/conf/config.json

# Build
sudo docker build -t watchtopus watchtopus/


##########################################
# Build elasticsearch with index mappings
# ready for watchtopus-server
##########################################

# Build elastics docker
sudo docker build -t watchtopus-elastics watchtopus-elastics/


##########################################
# Push dockers to docker registry
##########################################

# Set AWS profile to "edi" in order use that profile credentials to push to AWS account
# Prior configuration should be done by editing the ~/.aws/credentials file as explained in the link below
# https://stackoverflow.com/questions/44243368/how-to-login-with-aws-cli-using-credentials-profiles
export AWS_PROFILE=edi

# Login to AWS ECS docker registry
OUTPUT="$(aws ecr get-login --no-include-email --region eu-central-1)"
cmd="sudo $OUTPUT"
eval $cmd

# Tag our docker local images that we've just build with "latest" tag
sudo docker tag watchtopus:latest 312452674585.dkr.ecr.eu-central-1.amazonaws.com/watchtopus:latest
sudo docker tag watchtopus-elastics:latest 312452674585.dkr.ecr.eu-central-1.amazonaws.com/watchtopus-elastics:latest

# Push image to docker registry
sudo docker push 312452674585.dkr.ecr.eu-central-1.amazonaws.com/watchtopus:latest
sudo docker push 312452674585.dkr.ecr.eu-central-1.amazonaws.com/watchtopus-elastics:latest