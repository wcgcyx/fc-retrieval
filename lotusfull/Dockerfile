# Copyright (C) 2020 ConsenSys Software Inc
FROM wcgcyx/lotusbase
COPY lotusfull/daemon/config.toml /root/.lotus/
COPY lotusfull/miner/config.toml /root/.lotusminer/
WORKDIR /app/lotus
ADD lotusfull/start-lotus-full-node.sh .
RUN chmod +x start-lotus-full-node.sh
ENTRYPOINT ./start-lotus-full-node.sh
EXPOSE 1234 2345