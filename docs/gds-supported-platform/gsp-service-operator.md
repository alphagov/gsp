# Using GSP Service Operator

## How to use it
GSP Service Operator is a tool we use to allow GSP users to write kubeyaml resources that will generate SQS Queues or S3 Buckets (access control via Service accounts), or RDS databases (Postgres, access control via credentials we provide).

Here's an example of an SQS Queue:
```yaml
apiVersion: v1
kind: List
items:
- kind: SQS
  apiVersion: queue.govsvc.uk/v1beta1
  metadata:
    name: alexs-test-queue
    namespace: sandbox-gsp-service-operator-test
    labels:
      group.access.govsvc.uk: alexs-test-principal
  spec:
    aws:
      messageRetentionPeriod: 3600
      maximumMessageSize: 1024
    secret: alexs-test-queue-secret
- kind: ServiceAccount
  apiVersion: v1
  metadata:
    name: alexs-test-sa
    namespace: sandbox-gsp-service-operator-test
    labels:
      group.access.govsvc.uk: alexs-test-principal
```
This will create an SQS Queue on AWS named `alexs-test-queue`, with a message retention period of 1 hour, and a maximum message size of 1KiB. It will also ensure you can get access to the created queue. It will store the queue URL in a secret named `alexs-test-queue-secret` that we will use below.

Here's an example of an S3 Bucket:

```
apiVersion: v1
kind: List
items:
- apiVersion: storage.govsvc.uk/v1beta1
  kind: S3Bucket
  metadata:
    labels:
      group.access.govsvc.uk: alexs-test-principal
    name: alexs-test-bucket
    namespace: sandbox-gsp-service-operator-test
  spec:
    aws: {}
    secret: alexs-test-bucket-secret
- kind: ServiceAccount
  apiVersion: v1
  metadata:
    name: alexs-test-sa
    namespace: sandbox-gsp-service-operator-test
    labels:
      group.access.govsvc.uk: alexs-test-principal
```

This will create an S3 Bucket on AWS including the name `alexs-test-bucket`. It will ensure you can get access to the created bucket via an IAM Role available to the Service Account. It will store the created bucket name and URL inside the specified secret.

Here's an example of a Postgres database:

```
kind: Postgres
apiVersion: database.govsvc.uk/v1beta1
metadata:
  name: alexs-test-db
  namespace: sandbox-gsp-service-operator-test
spec:
  aws:
    instanceType: db.t3.medium
  secret: alexs-test-db-secret
```

This will create a Postgres database on AWS including the name `alexs-test-db`, with an instance type of `db.t3.medium`. It will ensure you can get access to the created database via the details written into the secret whose name you specify (it will create the secret for you if it does not already exist). It will store details such as the hostname, port, username, and password in this secret.

Here's an example of a ElastiCache Cluster:

```
kind: ElasticacheCluster
apiVersion: cache.govsvc.uk/v1beta1
metadata:
  name: alexs-test
  namespace: ondemand-0-main
spec:
  aws:
    nodeType: cache.t3.micro
    engineVersion: 5.0.6
    numCacheClusters: 2
  secret: alexs-test-cluster-secret
```

This will create an ElastiCache redis cluster on AWS including the name `alexs-test`, with an instance type of `cache.t3.micro`. It will ensure you can get access to the created cluster via the details written into the secret whose name you specify (it will create the secret for you if it does not already exist). It will store details such as the hostname, port, username, and password in this secret.
The combination of the following must fit into 40 characters, be made of alphanumeric characters and hyphens, and cannot contain two consecutive hyphens:
* your cluster's name (e.g. `sandbox`)
* a hyphen (`-`)
* your namespace's name (e.g. `sandbox-main`)
* a hyphen (`-`)
* your ElasticacheCluster resource's name

## How to connect to a created queue

The URL of the Queue will be stored inside the `secret` you specified as `QueueURL` (in addition, if you specified the `redriveMaxReceiveCount` parameter in the spec a redrive policy will have been configured with it pointing at the queue URL stored in key `DLQueueURL`). If you make a pod like:
```
apiVersion: v1
kind: Pod
metadata:
  name: alexs-test-pod
spec:
  serviceAccountName: alexs-test-sa
  containers:
  - name: myapp-container
    image: governmentpaas/awscli
    command: ['sleep', '1000000']
    volumeMounts:
    - name: secrets
      mountPath: /secrets
  volumes:
  - name: secrets
    secret:
      secretName: alexs-test-queue-secret
```

You will be able to access the URL of the Queue from inside your pod using `cat /secrets/QueueURL`.

The service account that the pod runs with will have access to a role that has access to the SQS queue. It should be possible to do the following (for this example we will use `gds sandbox kubectl exec -n sandbox-gsp-service-operator-test alexs-test-pod -it /bin/ash`):

```
/ # aws sqs send-message --queue-url $(cat /secrets/QueueURL) --message-body sup --region eu-west-2
{
    "MD5OfMessageBody": "2eeecd72c567401e6988624b179d0b14",
    "MessageId": "ac0f61ca-29a7-4eef-b998-831c7ed37ff3"
}
/ # aws sqs receive-message --queue-url $(cat /secrets/QueueURL) --region eu-west-2
{
    "Messages": [
        {
            "Body": "sup",
            "ReceiptHandle": "AQEBwFCRxEEt8T0NdFTy+F53zdQsVenKd6ZrMQyvsheq78rzsJOOr6255u8h4aAUxkRsXo9DKBxM3jI+fNcRPVEtNtfQqacdaJfYcxBs9rp0ogmHUpvfMCO27tjfUl5jqK3EEQ8fUG1SlDrOR22OwTZ73w2piZP7w6AWwEU5ohujVmcC6O/q44gI651lXP1HNHW9ZCMPtQdy0rdGtqpa/gcW8E2WYk7IvHD3SgSSzGkhMldT+VoPswNO1KEonjvP2DsJpiqlxacvmE4WHoxMlEufqNgYdxgSntKdAsig5/mRfjdqKuA39xe5X6gW7C9/8p7+I1UklO1rbPGcqNlssWqEuovHfqS+bpOGz2RvUxRNBTEkAKT+k7JRIPU9fJ5Y4OhmhbrMqQuv5x3U7jofTxTkESfmeibASfF1kAOx3+QmT6Mz0PF6C84vTy+lsMDZkKod+y4f8YYuvZvTJGOSMwP1fg==",
            "MD5OfBody": "2eeecd72c567401e6988624b179d0b14",
            "MessageId": "ac0f61ca-29a7-4eef-b998-831c7ed37ff3"
        }
    ]
}
```

## How to connect to a created bucket

The URL of the Bucket will be stored inside the `secret` you specified as `S3BucketURL`, and the name will be in the `S3BucketName` key. If you make a pod like:
```
apiVersion: v1
kind: Pod
metadata:
  name: alexs-test-pod
spec:
  serviceAccountName: alexs-test-sa
  containers:
  - name: myapp-container
    image: governmentpaas/awscli
    command: ['sleep', '1000000']
    env:
    - name: S3_BUCKET_NAME
      valueFrom:
        secretKeyRef:
          name: alexs-test-bucket-secret
          key: S3BucketName
```

You will be able to access the URL of the Bucket from inside your pod using `cat /secrets/S3BucketURL`.

The service account that the pod runs with will have access to a role that has access to the S3 bucket. It should be possible to do the following (for this example we will use `gds sandbox kubectl exec -n sandbox-gsp-service-operator-test alexs-test-pod -it /bin/ash`):

```
/ # echo hello > world
/ # aws s3 cp ./world s3://$S3_BUCKET_NAME/world --region eu-west-2
upload: ./world to s3://sandbox-sandbox-gsp-service-operator-test-alexs-test-bucket/world
/ # aws s3 cp s3://$S3_BUCKET_NAME/world ./downloaded --region eu-west-2
download: s3://sandbox-sandbox-gsp-service-operator-test-alexs-test-bucket/world to ./downloaded
/ # cat downloaded
hello
```

## How to connect to a created Postgres database

If you make a pod like the one above:

```
apiVersion: v1
kind: Pod
metadata:
    name: alexs-test-pod
spec:
    containers:
    -   name: myapp-container
        image: governmentpaas/psql
        command: ['sleep', '1000000']
        env:
        - name: PGHOST
          valueFrom:
            secretKeyRef:
              name: alexs-test-db-secret
              key: Endpoint
        - name: PGPORT
          valueFrom:
            secretKeyRef:
              name: alexs-test-db-secret
              key: Port
        - name: PGUSER
          valueFrom:
            secretKeyRef:
              name: alexs-test-db-secret
              key: Username
        - name: PGPASSWORD
          valueFrom:
            secretKeyRef:
              name: alexs-test-db-secret
              key: Password
```

You will be able to exec into this pod and get a PostgreSQL prompt.

```
$ gds sandbox kubectl exec -n sandbox-gsp-service-operator-test alexs-test-pod -c myapp-container -it /usr/bin/psql postgres
psql (11.5, server 10.7)
SSL connection (protocol: TLSv1.2, cipher: ECDHE-RSA-AES256-GCM-SHA384, bits: 256, compression: off)
Type "help" for help.

postgres=>
```

You could also get the read endpoint using the ReadEndpoint key.

## How to connect to a created ElastiCache Cluster

If you make a pod like the one above:

```
apiVersion: v1
kind: Pod
metadata:
  name: alexs-test-pod
spec:
  containers:
  - name: myapp-container
    image: alexmonk/redisandstunnel
    command: ['sleep', '1000000']
    env:
    - name: STUNNEL_HOSTNAME
      valueFrom:
        secretKeyRef:
          name: alexs-test-cluster-secret
          key: ClusterPrimaryRedisHostname
    - name: STUNNEL_PORT
      valueFrom:
        secretKeyRef:
          name: alexs-test-cluster-secret
          key: ClusterPrimaryRedisPort
    - name: REDISCLI_AUTH
      valueFrom:
        secretKeyRef:
          name: alexs-test-cluster-secret
          key: SecretAuthToken
```

You will be able to exec into this pod and connect to Redis after configuring and running stunnel:
```
$ gds sandbox kubectl exec -n sandbox-gsp-service-operator-test alexs-test-pod -c myapp-container -it /bin/bash
root@alexs-test-pod:/data# echo "  connect = $STUNNEL_HOSTNAME:$STUNNEL_PORT" >> /etc/stunnel/redis-cli.conf
root@alexs-test-pod:/data# stunnel /etc/stunnel/redis-cli.conf
root@alexs-test-pod:/data# redis-cli -h 127.0.0.1
127.0.0.1:6379> hgetall *
(empty list or set)
127.0.0.1:6379>
```

## How it works
You don't need to know this to use it, this information is for cluster operators.
GSP Service Operator consists of a container that runs essentially a daemon, and a kubeyaml config that sets up the container, provides a bunch of custom resource definitions (e.g., there is a definition in there for SQS Queues), etc. - it also gives the container access to interact with the cluster.
The daemon monitors the k8s cluster for such custom resources being created and will create the requested SQS/Database resources.
