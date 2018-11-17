# Watchtopus
Easily fully-monitor a network of servers with a quick installation of the agents and the central server.
<br>
Watchtopus is an educational project and not suitable (yet) for use in production.
<br><br>
A DR document that describes the Watchtopus project, can be found [here](https://docs.google.com/document/d/1jAmNmHwWiGXkTauNhiiRrf9f3IvOTqSMdzgIomw__r0).
<br>
Usage instructions can be found [here](https://docs.google.com/document/d/1MmbKV-CGezTQLdcqOLPWkzRsNAsPzM1vPB_Zkm-BQWk).


## Quick Installation of central server
1. Configure aws or any other docker registry provider that holds the docker images:
```bash
aws configure
```
2. Login to ECR (docker registry)
```bash
export AWS_PROFILE=edi #(Optional)
OUTPUT="$(aws ecr get-login --no-include-email --region eu-central-1)"
eval "sudo $OUTPUT"
```
3. Start all servers with docker-compose
```bash
cd ~/go/src/watchtopus/deployment
sudo docker-compose up -d
```

## Quick Installation of agent
1. Install Go
```bash
wget https://dl.google.com/go/go1.11.1.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.11.1.linux-amd64.tar.gz
echo "export PATH=\$PATH:/usr/local/go/bin" >> ~/.bashrc
```
2. Download watchtopus repository
```bash
sudo yum install git
cd ~/go/src/
git clone https://github.com/edibusl/watchtopus.git
cd ~/go/src/watchtopus
bash deployment/install_go_packages.sh
```
3. Set ping  permissions
```bash
sudo sysctl -w net.ipv4.ping_group_range="0   2147483647"
```
4. Compile & run the agent
```bash
cd ~/go/src/watchtopus
go run ./agent/main.go &
```


## Developers - contributing

### Running local separate services

1. The following servers should be started before running locally the watchtopus server:
```bash
sudo docker rm -f elastics; sudo docker run -d --name elastics -p 9200:9200 -p 9300:9300 -e "discovery.type=single-node" 312452674585.dkr.ecr.eu-central-1.amazonaws.com/watchtopus-elastics
sudo docker rm -f kibana;sudo docker run -d --name kibana -p 5601:5601 -e ELASTICSEARCH_URL="http://172.17.0.1:9200" docker.elastic.co/kibana/kibana-oss:6.4.2
```

2. Manually configuring ElasticSearch mappings of hosts index (can be run using Kibana DevTools)
```
PUT hosts
PUT hosts/_mapping/_doc
{
  "properties": {}
}
```

3. Manually configuring ElasticSearch mappings of metrics index (can be run using Kibana DevTools)
```
PUT metrics
PUT metrics/_mapping/_doc
{
    "properties":{
        "key":{"type":"keyword"},
        "val":{"type":"float"},
        "category":{"type":"keyword"},
        "subcategory":{"type":"keyword"},
        "component":{"type":"keyword"},
        "timestamp":{"type": "date"},
        "hostId": {"type": "keyword"},
        "hostIp": {"type": "keyword"}
    }
}
```

### Running unit tests
```bash
cd ~/go/src/watchtopus
go test server/tests/* -v
go test agent/tests/* -v
```

### Building & deploying to docker registry
```bash
cd ~/go/src/watchtopus/deployment
bash build_dockers.sh
```
