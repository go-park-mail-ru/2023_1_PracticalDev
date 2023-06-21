#!/bin/bash

logs_directory="./logs/postgres"
reports_directory="./logs/reports"
json_files=$(find "$logs_directory" -name "*.json")
for file in $json_files; do
  log_filename=$(basename "$file" .json)
  report_filename="${log_filename}.html"
  
  pgbadger "${file}" -O "./logs/reports" -o "${report_filename}"
done
