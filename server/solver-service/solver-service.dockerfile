# Use Python 3.11.5 as base image
FROM python:3.11.5-slim

# Set working directory
WORKDIR /app

# Install system dependencies required for scientific packages and protobuf compiler
RUN apt-get update && apt-get install -y \
    build-essential \
    clang \
    gringo \
    protobuf-compiler \
    && rm -rf /var/lib/apt/lists/*

# Copy requirements first to leverage Docker cache
COPY requirements.txt .

# Install Python dependencies
RUN pip install --upgrade pip && \
    pip install --no-cache-dir -r requirements.txt

# Create necessary directories
RUN mkdir -p asp_program proto/logs

# Copy proto files first
COPY proto/logs/logger.proto proto/logs/

# Create __init__.py files to make proto a proper Python package
RUN touch proto/__init__.py proto/logs/__init__.py

# Generate Python code from proto files
RUN python -m grpc_tools.protoc \
    -I. \
    --python_out=. \
    --grpc_python_out=. \
    proto/logs/logger.proto

# Copy ASP program files
COPY asp_program/program.lp asp_program/
COPY asp_program/colors.lp asp_program/

# Copy the rest of the application code
COPY . .

# Expose the ports the app runs on
EXPOSE 80 

# Set environment variables
ENV PYTHONUNBUFFERED=1
ENV PORT=80
ENV PYTHONPATH=/app

# Command to run the application
CMD ["gunicorn", "--bind", "0.0.0.0:80", "app:app"]