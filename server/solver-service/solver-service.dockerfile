# Use Python 3.11 as base image (to support numpy 2.0.0)
FROM python:3.11-slim

# Set working directory
WORKDIR /app

# Install system dependencies required for scientific packages
RUN apt-get update && apt-get install -y \
    build-essential \
    clang \
    gringo \
    && rm -rf /var/lib/apt/lists/*

# Copy requirements first to leverage Docker cache
COPY requirements.txt .

# Install Python dependencies
RUN pip install --upgrade pip && \
    pip install --no-cache-dir -r requirements.txt

# Copy the rest of the application code
COPY . .

# Create a directory for ASP program files
RUN mkdir -p asp_program

# Copy ASP program files
COPY asp_program/program.lp asp_program/
COPY asp_program/colors.lp asp_program/

# Expose the port the app runs on
EXPOSE 80 

# Set environment variables
ENV PYTHONUNBUFFERED=1
ENV PORT=80

# Command to run the application
CMD ["gunicorn", "--bind", "0.0.0.0:80", "app:app"]