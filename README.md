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

## Simulation Scope & Constraints
- Each auction runs with fixed attributes (20 by default).
- Bidders may or may not respond within the auction timeout.
- Each auction consumes 1 vCPU and 10 MB memory (configurable).
- Resource limits are enforced per auction using semaphore channels.
- No external queues (Kafka, Redis) are used — all concurrency is handled in-process.
- Bidders are fully simulated; there is no network communication.
- Results are printed to the console, and detailed JSON files are generated per auction.
- CPU concurrency is governed by runtime.GOMAXPROCS and a CPU semaphore.
- Memory concurrency is approximated using a semaphore channel (1 token ≈ 1 MB). 

---

## Directory Structure

```bash
Auction-Simulator/
├── app/
│ └── main.go # Entry point for the simulator
├── config/
│ └── config.go # Environment config loader and validator
├── auction_house/
│ └── auctionhouse.go # Auction logic and winner determination
├── bidders/
│ └── bidder.go # Bidder simulation
├── domain/
│ └── domain.go # Common structs (Bid, AuctionResult, etc.)
├── output/
│ └── auction_<id>.json # Results json files
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
| `NUM_AUCTIONS`         | Total number of auctions to simulate          | required         |
| `AUCTION_TIMEOUT_MS`   | Timeout per auction in milliseconds           | required         |
| `AUCTION_VCPU`         | VCPU cost per auction                         | required         |
| `AUCTION_MEMORY`       | Memory cost per auction (MB)                  | required         |
| `VCPU`                 | Total VCPU available for concurrency          | required         |
| `MEMORY`               | Total memory available (MB)                   | required         |

---

## Running the Simulator

1. Set environment variables:

    ```bash
    NUM_BIDDERS=100
    NUM_AUCTIONS=40
    AUCTION_TIMEOUT_MS=100
    AUCTION_VCPU=1
    AUCTION_MEMORY=10
    VCPU=4
    MEMORY=100
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

   - **CPU (vCPU):**
    A buffered channel (cpuSem) limits the number of concurrently running goroutines to the number of available vCPUs. The AUCTION_VCPU environment variable controls this value; if unset, Go defaults GOMAXPROCS to the number of available CPU cores.

    - **Memory (RAM):**
    Another buffered channel (memSem) approximates memory usage by requiring tokens per auction. The AUCTION_MEMORY environment variable sets this value. Dynamic RAM measurement is feasible (e.g., via runtime.MemStats), but this version standardizes memory per auction for reproducibility.


2. Configuration Validation

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