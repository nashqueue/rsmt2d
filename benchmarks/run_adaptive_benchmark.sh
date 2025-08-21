#!/bin/bash

# Adaptive benchmark runner with configurable iterations per k value
# Usage: ./run_adaptive_benchmark.sh [max_k] [timeout_minutes]

MAX_K=${1:-1024}
TIMEOUT=${2:-360}

echo "================================================"
echo "Adaptive Iteration Benchmark Runner"
echo "================================================"
echo ""
echo "This will run benchmarks with optimized iteration counts:"
echo "  k=32:  50 iterations (~2 min)"
echo "  k=64:  30 iterations (~3 min)" 
echo "  k=128: 20 iterations (~20 min)"
echo "  k=256: 15 iterations (~75 min)"
echo "  k=512: 10 iterations (~200 min)"
echo "  k=1024: 5 iterations (~300 min)"
echo ""
echo "Configuration:"
echo "  Max k: $MAX_K"
echo "  Timeout: $TIMEOUT minutes"
echo ""
echo "Starting in 3 seconds..."
sleep 3

# Run the adaptive benchmark
python3 generate_analysis_report.py 1 $TIMEOUT $MAX_K

echo ""
echo "Benchmark complete! Check the results/ directory for:"
echo "  - benchmark_results_*.txt (raw data)"
echo "  - benchmark_analysis_*.md (analysis report)"