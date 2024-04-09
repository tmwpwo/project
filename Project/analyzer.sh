# Build the Go program
go build .

# Get the directory containing the executable
executable_directory=$(pwd)

# Add the directory to the PATH
if [[ ":$PATH:" != *":$executable_directory:"* ]]; then
    echo "export PATH=\$PATH:$executable_directory" >> ~/.bashrc
    echo "Directory added to PATH."
else
    echo "Directory is already in PATH."
fi

# Reload the .bashrc file to apply changes
source ~/.bashrc

echo "Build complete and directory added to PATH."