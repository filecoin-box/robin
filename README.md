# Robin

Robin is the monitoring and alarm of win block and loss block. It now supports sending messages to slack.

**build**

- git clone https://github.com/luluup777/robin.git
- make all

**start**

```
robin -h
NAME:
   robin - mining monitoring and alarm

USAGE:
   robin [global options] command [command options] [arguments...]

VERSION:
   v0.1

COMMANDS:
   run      start robin
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)
```

- vim `robin/conf/robin.yaml`

  ```
  notify:
    platform: slcak
    token: xxx
    channel: xxx
  
  monitor:
    fullnode_api_info: xxx
    minerId: f011111,f022222
  ```

  - `notify`:
    - Fill in the `token` and `channel` of slack to allow status reporting
  - `monitor`:
    - `fullnode_api_info`:  it is lotus environment variables: `FULLNODE_API_INFO`, ps: `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBbGxvdyI6WyJyZWFkIiwid3JpdGUiLCJzaWduIiwiYWRtaW4iXX0.oeId55L_CiYhca8rSA7xKm6qyMquIGMAR5q_g-lgdFF:/ip4/127.0.0.1/tcp/1234/http`
    - `minerId`: If there are multiple miners to monitor, separate them with commas

- run
  - After filling in `robin.yaml`, you can start `robin`. please execute: `./robin run`
  - support to change `robin.yaml` in **runtime**.
    - Increase or decrease monitoring of miners
    - change daemon