<img align="right" width="30%" src="./assets/image.png">

# Proof of Work server and client

## Task description

Design and implement “Word of Wisdom” tcp server.
<ul>
  <li>TCP server should be protected from DDOS attacks with the Prof of Work (<a href="https://en.wikipedia.org/wiki/Proof_of_work">https://en.wikipedia.org/wiki/Proof_of_work</a>), the challenge-response protocol should be used.</li>
  <li>The choice of the POW algorithm should be explained.</li>
  <li>After Prof Of Work verification, server should send one of the quotes from “word of wisdom” book or any other collection of the quotes.</li>
  <li>• Docker file should be provided both for the server and for the client that solves the POW challenge</li>
</ul>

## How to run
To run locally, put your `config.yaml` file in `configs` directory. Default config file is already presented in repo. You can run server and clint separately with `make server` and `make client` commands respectively.
For deployment, Dockerfiles for both the server and the client is also provided. Run containers with command `make up`
