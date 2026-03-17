# Operations

Chains of commands to be executed in sequence, operations are user defined and when run all operations should be chained together into one command so the user can execute all oeprations together. operations should have custom names and should have their own creation ui, i'm thinking the operation name could be at the top and below there could be numbered list of commands to execute with an add new row/button at the bottom that starts the new command prompt when the user clicks it with space bar

### Operation tar file and backup

1. cd /path/to/file
2. tar file with options
3. sftp file with options
4. rm /path/to/archive.tar

### Operation restart docker

1. systemctl restart docker.service
2. systemctl is-active --wait docker.service