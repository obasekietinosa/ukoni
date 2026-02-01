#!/bin/bash
set -e

echo "Starting setup for Ukoni..."

# 1. System Updates & Dependencies
echo "[1/5] Installing system dependencies..."
# Update apt
sudo apt-get update
# Install curl, git, build-essential, postgres
sudo apt-get install -y curl git build-essential postgresql postgresql-contrib

# 2. Go Installation
echo "[2/5] Checking Go installation..."
if ! command -v go &> /dev/null; then
    echo "Go not found. Installing Go 1.25.5..."
    # Download Go
    wget https://go.dev/dl/go1.25.5.linux-amd64.tar.gz
    # Remove old installation and install new
    sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.25.5.linux-amd64.tar.gz
    # Clean up
    rm go1.25.5.linux-amd64.tar.gz

    # Add to PATH for this session
    export PATH=$PATH:/usr/local/go/bin

    # Add to PATH for future sessions if not already present
    if ! grep -q "/usr/local/go/bin" ~/.profile; then
        echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.profile
        echo "Go added to ~/.profile."
    fi
else
    echo "Go is already installed."
fi

# 3. Node.js Installation
echo "[3/5] Checking Node.js installation..."
if ! command -v node &> /dev/null; then
    echo "Node.js not found. Installing Node.js 22 (LTS)..."
    curl -fsSL https://deb.nodesource.com/setup_22.x | sudo -E bash -
    sudo apt-get install -y nodejs
else
    echo "Node.js is already installed."
fi

# 4. Database Setup
echo "[4/5] Configuring PostgreSQL..."
sudo systemctl start postgresql

# Create user 'etin' if it doesn't exist
if ! sudo -u postgres psql -tAc "SELECT 1 FROM pg_roles WHERE rolname='etin'" | grep -q 1; then
    echo "Creating database user 'etin'..."
    sudo -u postgres psql -c "CREATE USER etin WITH PASSWORD 'etin';"
else
    echo "User 'etin' already exists."
fi

# Create database 'ukoni' if it doesn't exist
if ! sudo -u postgres psql -tAc "SELECT 1 FROM pg_database WHERE datname='ukoni'" | grep -q 1; then
    echo "Creating database 'ukoni'..."
    sudo -u postgres psql -c "CREATE DATABASE ukoni OWNER etin;"
else
    echo "Database 'ukoni' already exists."
fi

# Ensure user has DB creation rights (optional but good practice)
sudo -u postgres psql -c "ALTER USER etin CREATEDB;"

# 5. Project Setup
echo "[5/5] Setting up project..."

# Setup API
echo "--- Setting up API ---"
if [ -d "api" ]; then
    cd api
    if [ -f "go.mod" ]; then
        echo "Downloading Go modules..."
        go mod download

        echo "Running migrations..."
        # Set DATABASE_URL for migration
        export DATABASE_URL="postgres://etin:etin@localhost:5432/ukoni?sslmode=disable"
        go run cmd/migrate/main.go up
    else
        echo "Warning: api/go.mod not found. Skipping API setup."
    fi
    cd ..
else
    echo "Warning: Directory 'api' not found. Skipping API setup."
fi

# Setup Web
echo "--- Setting up Web ---"
if [ -d "web" ]; then
    cd web
    if [ -f "package.json" ]; then
        echo "Installing Node dependencies..."
        npm install
    else
        echo "Warning: web/package.json not found. Skipping Web setup."
    fi
    cd ..
else
    echo "Warning: Directory 'web' not found. Skipping Web setup."
fi

echo "----------------------------------------------------------------"
echo "Setup Complete!"
echo ""
echo "To run the API:"
echo "  cd api"
echo "  export DATABASE_URL=\"postgres://etin:etin@localhost:5432/ukoni?sslmode=disable\""
echo "  go run cmd/api/main.go"
echo ""
echo "To run the Web App:"
echo "  cd web"
echo "  npm run dev"
echo "----------------------------------------------------------------"
