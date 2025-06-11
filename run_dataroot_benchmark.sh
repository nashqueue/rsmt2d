#!/bin/bash

# Script to run Dataroot benchmarks and generate visualization

set -e

echo "Running Dataroot benchmarks..."

# Run the benchmark and save results to a file
go test -bench=BenchmarkDatarootGeneration -benchmem -count=5 -timeout=30m > benchmark_results.txt 2>&1

echo "Benchmark results saved to benchmark_results.txt"

# Extract benchmark data and create CSV for visualization
echo "Extracting benchmark data..."

# Create CSV file with headers
echo "benchmark_name,k_value,eds_size,ns_per_op,bytes_per_op,allocs_per_op" > benchmark_data.csv

# Parse benchmark results and extract data
grep "BenchmarkDatarootGeneration" benchmark_results.txt | while read line; do
    # Extract full benchmark name
    benchmark_name=$(echo "$line" | awk '{print $1}')
    
    # Extract k value from benchmark name (e.g., k=32_EDS=64x64)
    k_value=$(echo "$line" | sed -n 's/.*k=\([0-9]*\)_.*/\1/p')
    
    # Extract EDS size (e.g., 64x64)
    eds_size=$(echo "$line" | sed -n 's/.*EDS=\([0-9]*x[0-9]*\).*/\1/p')
    
    # Extract performance metrics
    ns_per_op=$(echo "$line" | awk '{print $3}')
    bytes_per_op=$(echo "$line" | awk '{print $5}')
    allocs_per_op=$(echo "$line" | awk '{print $7}')
    
    # Only add if we have valid data
    if [[ -n "$k_value" && -n "$eds_size" && -n "$ns_per_op" ]]; then
        echo "$benchmark_name,$k_value,$eds_size,$ns_per_op,$bytes_per_op,$allocs_per_op" >> benchmark_data.csv
    fi
done

echo "Benchmark data extracted to benchmark_data.csv"

# Create Python script for visualization
echo "Creating visualization script..."

cat > create_benchmark_plot.py << 'EOF'
#!/usr/bin/env python3

import pandas as pd
import matplotlib.pyplot as plt
import numpy as np
import seaborn as sns

# Set style
plt.style.use('seaborn-v0_8')
sns.set_palette("husl")

# Read the benchmark data
try:
    df = pd.read_csv('benchmark_data.csv')
    
    # Convert ns_per_op to milliseconds for better readability
    df['ms_per_op'] = df['ns_per_op'] / 1_000_000
    
    # Convert bytes to MB
    df['mb_per_op'] = df['bytes_per_op'] / (1024 * 1024)
    
    # Extract EDS width from eds_size (e.g., "64x64" -> 64)
    df['eds_width'] = df['eds_size'].str.split('x').str[0].astype(int)
    
    # Calculate total shares (EDS width squared)
    df['total_shares'] = df['eds_width'] ** 2
    
    # Sort by k_value for consistent plotting
    df = df.sort_values('k_value')
    
    print("Benchmark Data Summary:")
    print(df[['k_value', 'eds_size', 'ms_per_op', 'mb_per_op', 'allocs_per_op']])
    
    # Create subplots
    fig, ((ax1, ax2), (ax3, ax4)) = plt.subplots(2, 2, figsize=(15, 12))
    fig.suptitle('Dataroot Generation Benchmark Results\n(EDS with Random Data)', fontsize=16, fontweight='bold')
    
    # Plot 1: Time vs K value
    ax1.plot(df['k_value'], df['ms_per_op'], 'bo-', linewidth=2, markersize=8)
    ax1.set_xlabel('K Value (Original Data Square Size)', fontsize=12)
    ax1.set_ylabel('Time (ms)', fontsize=12)
    ax1.set_title('Dataroot Generation Time vs K Value', fontsize=14)
    ax1.grid(True, alpha=0.3)
    ax1.set_yscale('log')
    
    # Add annotations for EDS sizes
    for i, row in df.iterrows():
        ax1.annotate(f'EDS: {row["eds_size"]}', 
                    (row['k_value'], row['ms_per_op']), 
                    textcoords="offset points", 
                    xytext=(0,10), 
                    ha='center', fontsize=10)
    
    # Plot 2: Memory usage vs K value
    ax2.plot(df['k_value'], df['mb_per_op'], 'ro-', linewidth=2, markersize=8)
    ax2.set_xlabel('K Value (Original Data Square Size)', fontsize=12)
    ax2.set_ylabel('Memory Usage (MB)', fontsize=12)
    ax2.set_title('Memory Usage vs K Value', fontsize=14)
    ax2.grid(True, alpha=0.3)
    ax2.set_yscale('log')
    
    # Plot 3: Time vs Total Shares
    ax3.plot(df['total_shares'], df['ms_per_op'], 'go-', linewidth=2, markersize=8)
    ax3.set_xlabel('Total EDS Shares', fontsize=12)
    ax3.set_ylabel('Time (ms)', fontsize=12)
    ax3.set_title('Dataroot Generation Time vs Total EDS Shares', fontsize=14)
    ax3.grid(True, alpha=0.3)
    ax3.set_yscale('log')
    ax3.set_xscale('log')
    
    # Plot 4: Allocations vs K value
    ax4.plot(df['k_value'], df['allocs_per_op'], 'mo-', linewidth=2, markersize=8)
    ax4.set_xlabel('K Value (Original Data Square Size)', fontsize=12)
    ax4.set_ylabel('Allocations per Operation', fontsize=12)
    ax4.set_title('Memory Allocations vs K Value', fontsize=14)
    ax4.grid(True, alpha=0.3)
    ax4.set_yscale('log')
    
    plt.tight_layout()
    
    # Save the plot
    plt.savefig('dataroot_benchmark_results.png', dpi=300, bbox_inches='tight')
    print("\nVisualization saved as 'dataroot_benchmark_results.png'")
    
    # Show the plot
    plt.show()
    
    # Create a summary table
    summary_df = df[['k_value', 'eds_size', 'total_shares', 'ms_per_op', 'mb_per_op', 'allocs_per_op']].copy()
    summary_df.columns = ['K Value', 'EDS Size', 'Total Shares', 'Time (ms)', 'Memory (MB)', 'Allocations']
    
    print("\n" + "="*80)
    print("BENCHMARK SUMMARY")
    print("="*80)
    print(summary_df.to_string(index=False, float_format='%.2f'))
    
    # Performance scaling analysis
    if len(df) > 1:
        print("\n" + "="*80)
        print("PERFORMANCE SCALING ANALYSIS")
        print("="*80)
        
        # Calculate scaling factors
        base_k = df.iloc[0]['k_value']
        base_time = df.iloc[0]['ms_per_op']
        
        for i, row in df.iterrows():
            k_ratio = row['k_value'] / base_k
            time_ratio = row['ms_per_op'] / base_time
            theoretical_ratio = k_ratio ** 2  # O(n^2) for EDS
            
            print(f"K={row['k_value']:3d}: {k_ratio:4.1f}x larger → {time_ratio:6.1f}x slower "
                  f"(theoretical O(n²): {theoretical_ratio:4.1f}x)")

except FileNotFoundError:
    print("Error: benchmark_data.csv not found. Please run the benchmark first.")
except Exception as e:
    print(f"Error creating visualization: {e}")
EOF

# Check if Python and required packages are available
echo "Checking Python dependencies..."
if command -v python3 >/dev/null 2>&1; then
    # Run the comparison analysis (no dependencies required)
    echo "Creating benchmark comparison analysis..."
    python3 analyze_benchmark_simple.py
    
    # Try to run the advanced visualization if packages are available
    echo ""
    echo "Attempting advanced visualization..."
    python3 create_benchmark_plot.py 2>/dev/null || echo "Advanced visualization requires: pip install pandas matplotlib seaborn"
else
    echo "Python3 not found. Please install Python3 to run analysis."
    echo "The benchmark data is available in benchmark_data.csv for manual analysis."
fi

echo "Benchmark complete!"
echo ""
echo "Files generated:"
echo "  - benchmark_results.txt: Raw benchmark output"
echo "  - benchmark_data.csv: Processed benchmark data"
echo "  - create_benchmark_plot.py: Visualization script"
echo "  - dataroot_benchmark_results.png: Benchmark visualization (if Python is available)"