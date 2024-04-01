![Repository Top Language](https://img.shields.io/github/languages/top/JavaHutt/hashcash)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/JavaHutt/hashcash)
![License](https://img.shields.io/badge/license-MIT-green)
![Coding all night)](https://img.shields.io/badge/coding-all%20night%20-purple)

# Proof of Work server and client

## Task description
<img align="right" width="30%" src="./assets/image.png">

Design and implement “Word of Wisdom” tcp server.
<ul>
  <li>TCP server should be protected from DDOS attacks with the Prof of Work (<a href="https://en.wikipedia.org/wiki/Proof_of_work">https://en.wikipedia.org/wiki/Proof_of_work</a>), the challenge-response protocol should be used.</li>
  <li>The choice of the POW algorithm should be explained.</li>
  <li>After Prof Of Work verification, server should send one of the quotes from “word of wisdom” book or any other collection of the quotes.</li>
  <li>Docker file should be provided both for the server and for the client that solves the POW challenge</li>
</ul>

## How to run
To run locally, put your `config.yaml` file in `configs` directory. Default config file is already presented in repo. You can run server and clint separately with `make server` and `make client` commands respectively. You'll also need to run Redis as a store.
For deployment, Dockerfiles for both the server, the client and redis is also provided. Run containers with command `make up`.

## Algorithm choice explained
For the "Word of Wisdom" TCP server that employs a Proof of Work (PoW) mechanism to mitigate the risk of DDoS attacks, an efficient and suitable choice of PoW algorithm is crucial. The chosen algorithm is Hashcash. This decision is based on several factors that make Hashcash particularly apt for this application:
<ul>
  <li>Simplicity and Efficiency</li>
  <li>Widely Used and Tested</li>
  <li>Adjustable Difficulty</li>
  <li>Statelessness</li>
  <li>Used by Bitcoin 💰</li>
</ul>
</ul>