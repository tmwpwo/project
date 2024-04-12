
go build .


executable_directory=$(pwd)


if [[ ":$PATH:" != *":$executable_directory:"* ]]; then
    echo "export PATH=\$PATH:$executable_directory" >> ~/.bashrc
    echo "Directory added to PATH."
else
    echo "Directory is already in PATH."
fi


source ~/.bashrc

echo "Build complete and directory added to PATH."