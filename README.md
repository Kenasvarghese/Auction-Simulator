# Auction Simulator

## Overview

The **Auction Simulator** is a Go-based system designed to simulate multiple concurrent auctions with multiple bidders. Each auction runs with a timeout, collects bids from bidders, determines the highest bidder, and tracks execution metrics such as duration and resources used.  

This project demonstrates Go concurrency patterns, resource management, and real-time auction evaluation.

---

## Features

- Simulates multiple auctions concurrently.  
- Supports configurable number of bidders and auctions.  
- Each auction is executed with a timeout.  
- Bidders respond with randomized bid values.  
- Determines the highest bidder per auction.  
- Supports resource-limiting using **semaphores** (CPU and memory).  
- Tracks and prints execution duration for all auctions.  
- Configurable via environment variables.  
- Clean separation of components:  
  - **Auction House**: Manages auctions and determines winners.  
  - **Bidder Simulation**: Generates bidder behavior.  

---

## Assumptions

1. Each auction has **fixed attributes** (20 by default).  
2. Each bidder may or may not respond within the auction timeout.  
3. Each auction uses **1 VCPU and 10 MB memory** (configurable).  
4. Resource limits are enforced **per auction** using semaphores.  
5. The system does not use external queues (Kafka, Redis) — all concurrency is handled in-process.  
6. Bidders are simulated — there is no network communication.  
7. Output is printed to console; JSON output files is also added per auction.  
8. CPU concurrency is limited using **runtime.GOMAXPROCS** and channel semaphores.  
9. Memory concurrency is simulated using a semaphore channel (1 token ≈ 1 MB).  

---

## Directory Structure

```bash
Auction-Simulator/
├── main.go # Entry point for the simulator
├── config/
│ └── config.go # Environment config loader and validator
├── auction_house/
│ └── auctionhouse.go # Auction logic and winner determination
├── bidders/
│ └── bidder.go # Bidder simulation
├── domain/
│ └── domain.go # Common structs (Bid, AuctionResult, etc.)
├── utils/
│ └── utils.go # Helper functions (attribute generator, etc.)
└── README.md
```

---

## Configuration

All parameters are loaded from **environment variables**.  

| Variable               | Description                                    | Default / Required |
|------------------------|-----------------------------------------------|------------------|
| `NUM_BIDDERS`          | Number of bidders per auction                 | required         |
| `NUM_ATTRIBUTES`       | Number of attributes per auction item         | required         |
| `NUM_AUCTIONS`         | Total number of auctions to simulate          | required         |
| `AUCTION_TIMEOUT_MS`   | Timeout per auction in milliseconds           | required         |
| `AUCTION_VCPU`         | VCPU cost per auction                          | required         |
| `AUCTION_MEMORY`       | Memory cost per auction (MB)                  | required         |
| `VCPU`                 | Total VCPU available for concurrency          | required         |
| `MEMORY`               | Total memory available (MB)                   | required         |

---

## Running the Simulator

1. Set environment variables:

    ```bash
    export NUM_BIDDERS=100
    export NUM_ATTRIBUTES=20
    export NUM_AUCTIONS=40
    export AUCTION_TIMEOUT_MS=100
    export AUCTION_VCPU=1
    export AUCTION_MEMORY=10
    export VCPU=4
    export MEMORY=100
    ```
2.  Run the Simulator (Go):

    ```bash
    go run ./app/main.go
    ```

3. Run the Simulator (Docker):

    ```bash
        docker-compose up
    ```


    **Note:**

        Docker runs the same binary as go run.
        Output files appear in the host output/ folder after execution.

4. Expected output:

    ```bash
    Starting 40 auctions with 100 bidders
    Auction 1 winner: BidderID=5 Price=87.23
    Auction 2 winner: BidderID=12 Price=92.14
    ...
    Completed 40 auctions in 152 ms
    ```
    **Note:**

        The exact timing values and winning bidder IDs will vary with each run because bidder responses are randomized and auctions run concurrently.

        Results are not deterministic across runs.

5. Output Folder

    - After execution, an output/ folder will be created containing one JSON file per auction:
    - Each file is named `auction_<id>.json` (e.g., `auction_001.json`).
    - The file contains the auction’s attributes, bids received, and winning bidder details.
    - These files allow post-run analysis of auction performance and bidder behavior.
    - Together, the output files serve as an audit trail for the simulator run.

## Design Decisions
1. Semaphores with Channels

    - CPU: A buffered channel (cpuSem) ensures no more goroutines run than available vCPUs.
    - Memory: Another buffered channel (memSem) approximates memory by requiring tokens per auction.

2. Validation at Startup

    - Config values are validated at the start; the program panics if required values are missing or invalid.

3. Simplicity vs Accuracy

    - Resource usage is approximated, not enforced at the OS level.
    - The design favors simplicity, concurrency control, and reproducibility.
---

## Future Improvements
   - Replace simulated memory control with real memory monitoring (runtime.MemStats).
   - Add persistence (e.g., auction history in a database).
   - Extend to distributed setup with multiple nodes handling auctions.
---