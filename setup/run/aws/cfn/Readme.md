# FILE system:

Created File system with ID  fs-d049bae4

 Mounted inise the Linux mahcine: https://docs.aws.amazon.com/efs/latest/ug/wt1-test.html
sudo mount -t nfs4 -o nfsvers=4.1,rsize=1048576,wsize=1048576,hard,timeo=600,retrans=2,noresvport fs-d049bae4.efs.eu-west-1.amazonaws.com:/ ~/nfs-tp-unicorn/ 




# Data Gen API

aws servicediscovery create-service --name data-gen-api --dns-config 'NamespaceId="ns-nvo7dvldqg7zlyl5",DnsRecords=[{Type="A",TTL="300"}]' --health-check-custom-config FailureThreshold=1 --region  eu-west-1

aws ecs register-task-definition --cli-input-json file:///home/nxf42609/projects/tpte-platform/tp-unicorn/setup/run/aws/cfn/data-gen-api-task.json --region eu-west-1

aws ecs create-service --cli-input-json file:///home/nxf42609/projects/tpte-platform/tp-unicorn/setup/run/aws/cfn/data-gen-api-serdis.json --region eu-west-1

# Data Saver

aws servicediscovery create-service --name data-saver --dns-config 'NamespaceId="ns-nvo7dvldqg7zlyl5",DnsRecords=[{Type="A",TTL="300"}]' --health-check-custom-config FailureThreshold=1 --region  eu-west-1

aws ecs register-task-definition --cli-input-json file:///home/nxf42609/projects/tpte-platform/tp-unicorn/setup/run/aws/cfn/data-saver-task.json --region eu-west-1

aws ecs create-service --cli-input-json file:///home/nxf42609/projects/tpte-platform/tp-unicorn/setup/run/aws/cfn/data-saver-serdis.json --region eu-west-1


# Data Retrieve

aws servicediscovery create-service --name data-retrieve --dns-config 'NamespaceId="ns-nvo7dvldqg7zlyl5",DnsRecords=[{Type="A",TTL="300"}]' --health-check-custom-config FailureThreshold=1 --region  eu-west-1

aws ecs register-task-definition --cli-input-json file:///home/nxf42609/projects/tpte-platform/tp-unicorn/setup/run/aws/cfn/data-retrieve-task.json --region eu-west-1

aws ecs create-service --cli-input-json file:///home/nxf42609/projects/tpte-platform/tp-unicorn/setup/run/aws/cfn/data-retrieve-serdis.json --region eu-west-1

# MSSQL
aws servicediscovery create-service --name sql-server-db --dns-config 'NamespaceId="ns-nvo7dvldqg7zlyl5",DnsRecords=[{Type="A",TTL="300"}]' --health-check-custom-config FailureThreshold=1 --region  eu-west-1

   >>> "registryArn": "arn:aws:servicediscovery:eu-west-1:485458181338:service/srv-onghvpx3ww47uwh4"


aws ecs register-task-definition --cli-input-json file:///home/nxf42609/projects/tpte-platform/tp-unicorn/setup/run/aws/cfn/mssql-task.json --region eu-west-1

aws ecs create-service --cli-input-json file:///home/nxf42609/projects/tpte-platform/tp-unicorn/setup/run/aws/cfn/mssql-serdis.json --region eu-west-1



## ATTENTION: 
Change registryArn": "arn:aws:servicediscovery:eu-west-1:485458181338:service/srv-dbxnp5o4um474ud7" from creating the disconvery endpoint

For MSSQL file permission look at here: https://docs.microsoft.com/en-us/sql/linux/sql-server-linux-docker-container-security?view=sql-server-ver15