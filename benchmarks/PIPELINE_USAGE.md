# Benchmark Pipeline

## Quick Start

```bash
# Run with defaults (up to k=1024, 360 min timeout)
python3 benchmark.py

# Test up to k=64 (fast, ~10 minutes)
python3 benchmark.py 30 64

# Full run up to k=1024 (~10 hours)  
python3 benchmark.py 360 1024

# Extended run up to k=2048 (~15 hours)
python3 benchmark.py 900 2048
```

## Usage

```bash
python3 benchmark.py [timeout_minutes] [max_k]
```

**Parameters:**
- `timeout_minutes`: Timeout in minutes (default: 360)
- `max_k`: Maximum k value to test (default: 1024)

## Adaptive Iterations

The pipeline automatically uses more iterations for smaller k values to ensure statistical reliability:

| k Value | Iterations | Approx. Time |
|---------|------------|--------------|
| k≤32    | 50         | ~2 minutes   |
| k≤64    | 30         | ~5 minutes   |
| k≤128   | 20         | ~20 minutes  |
| k≤256   | 15         | ~75 minutes  |
| k≤512   | 10         | ~200 minutes |
| k≤1024  | 5          | ~300 minutes |
| k>1024  | 3          | ~720 minutes |

## Output Files

- `results/benchmark_results_YYYYMMDD_HHMMSS.txt` - Raw benchmark data
- `results/benchmark_analysis_YYYYMMDD_HHMMSS.md` - Analysis report with performance tables and visualizations

## What It Measures

The pipeline benchmarks different tree construction approaches for Extended Data Square (EDS) generation:

- **NMT**: Namespaced Merkle Trees (baseline)
- **Stage1**: Merkle columns, NMT rows
- **Stage2**: Stage1 + Merkle for bottom half rows
- **Stage3**: Stage2 + hybrid row trees
- **MerkleTree**: Pure Merkle trees (fastest)

Each approach is measured for:
- Execution time (ms)
- Memory usage (MB)
- Speedup vs NMT baseline
- Statistical variance (CV%)