# Kalah-AI Gigel

## Instruction

### Build

In order to build our agent, you must have [Go](https://golang.org) installed.
The build process is simple:

```shell
cd gigel
go build
```

Now you should have an executable called `gigel`. Which you can start by simply `./gigel`.

### Play

The game engine for our agent is provided by **The University of Manchester**,
for the **COMP34120** course. In order to run a game between two agents use:

```shell
java -jar ManKalah.jar "<AGENT 0 STRING>" "<AGENT 1 STRING>"
```

Some agents use a middleman to communicate with the game engine, eg: `netcat`.
Our agent uses `localhost:12345` for the agent and `localhost:12340` for a basic
user "interface" that allows humans to play against our agent.

```shell
java -jar ManKalah.jar "nc localhost 12345" "nc localhost 12340"
```

You can simply change this inside `gigel/main.go`. Another two aspects of our agent that can be configured are the heuristics weights and depth of min-max tree.

## Note

Please note that this is a **PUBLIC DEMONSTRATION** version of our agent. The development version is more advanced with analytics, better performance and overall improved code quality. The demo version is purposefully downgraded, as **COMP34120** is a competitive course (and it's also fun and satisfying to find your own successful strategy).

And yes, that's why lots of variables are named **hodor**. It's obfuscation-ish.

## Strategy

Gigel is a simple min-max search agent with alpha-beta pruning, with multiple optimizations to improve speed compared to other min-max agents. The tournament version can run at incredibly high depths, with quick decision times. Other optimizations involve low-level, compiler and multi-threading processing.

The main feature of Gigel, is the weight based heuristic. The tournament configuration was removed for this version.

## Performance

With the right weights and depth, our agent can defeat all agents found in `agents/` with a 100% win-rate, starting in either position, with a swift decision time of 0.5s / move.

During testing we noticed a 10-20% chance to lose and 20% chance to draw against some agents from other teams, when starting first. This can easily be fixed with more depth, increasing our win-rate to 100% and the response time to 3-5s / move.

## Other agents

All agents inside `agents/` are not owned, developed or in any way related to our team. They were either provided by the course staff for reference or salvaged from the tournament server. Thus our LICENSE does not extended over `agents/` or `ManKalah.jar`.

## Simulator

The kalah game simulator which can be found in `gigel/simulator.go` is copies the exact behaviour of the `ManKalah.jar` game engine, providing an interface in Go. The simulator can be extended to simulate the protocol (found in `docs/`).

## Releases

On the GitHub repo page, you can find the tournament version of **Gigel**, pre-compiled for **Darwin x64**, **Linux x64** and **Linux x86**. For windows users please check [here](https://distrowatch.com). The source code is not available to the public as explained above.
