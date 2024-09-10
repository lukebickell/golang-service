# Define the function to build the Go binary and package it into a zip
function Build-GoLambda {

      # Set the environment variables
    Write-Host "Setting environment variables..."
    $env:CGO_ENABLED = "0"
    $env:GOOS = "linux"
    $env:GOARCH = "amd64"

    # Build the Go binary with the same flags
    Write-Host "Building Go binary..."
    $buildCommand = "go build -ldflags='-s -w' -o main main.go"
    Invoke-Expression $buildCommand

    # Check if the build was successful
    if (-Not (Test-Path "./main")) {
        Write-Host "Build failed: binary not found." -ForegroundColor Red
        exit 1
    }

    # Create bin directory if it doesn't exist
    if (-Not (Test-Path "./bin")) {
        Write-Host "Creating bin/ directory..."
        New-Item -ItemType Directory -Path "./bin"
    }

    # Create a zip file with the binary
    Write-Host "Zipping the binary..."
    Compress-Archive -Path "./main" -DestinationPath "./bin/deployment.zip" -Force

    # Clean up by removing the binary
    Write-Host "Cleaning up..."
    Remove-Item "./main" -Force

    Write-Host "Build and package process complete!" -ForegroundColor Green
}

# Run the build function
Build-GoLambda
