# Using GSP Service Operator

## How to use it
GSP Service Operator is a tool we use to allow GSP users to write kubeyaml resources that will generate SQS Queues or Databases (via Principal objects for access control). Here's an example of an SQS Queue:
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
- kind: Principal
  apiVersion: access.govsvc.uk/v1beta1
  metadata:
    name: alexs-test-princ
    namespace: sandbox-gsp-service-operator-test
    labels:
      group.access.govsvc.uk: alexs-test-principal
```
This will create an SQS Queue on AWS named alexs-test-queue, with a message retention period of 1 hour, and a maximum message size of 1KiB. It will also ensure you can get access to the created queue. It will store the queue URL in a secret named `alexs-test-queue-secret` that we will use below.

Here's an example of a Postgres database:

```
apiVersion: v1
kind: List
items:
- kind: Postgres
  apiVersion: database.govsvc.uk/v1beta1
  metadata:
    name: alexs-test-db
    namespace: sandbox-gsp-service-operator-test
    labels:
      group.access.govsvc.uk: alexs-test-principal
  spec:
    aws:
      instanceType: db.t3.medium
    secret: alexs-test-db-secret
- kind: Principal
  apiVersion: access.govsvc.uk/v1beta1
  metadata:
    name: alexs-test-princ
    namespace: sandbox-gsp-service-operator-test
    labels:
      group.access.govsvc.uk: alexs-test-principal
```

This will create a Postgres database on AWS including the name alexs-test-db, with an instance type of db.t3.medium. It will ensure you can get access to the created queue via the details written into the secret whose name you specify (it will create the secret for you if it does not already exist). It will store details such as the hostname, port, username, and password in this secret.

## How to connect to a created queue

The URL of the Queue will be stored inside the `secret` you specified as `QueueURL`. If you make a pod like:
```
apiVersion: v1
kind: Pod
metadata:
  name: alexs-test-pod
  annotations:
    iam.amazonaws.com/role: svcop-sandbox-gsp-service-operator-test-alexs-test-princ
spec:
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

When the Principal creation is handled a role like svcop-sandbox-sandbox-gsp-service-operator-test-alexs-test-princ will have been created - in the form of svcop-{cluster}-{namespace}-{resourcename}. Your namespace will have an annotation that allows it to access such roles, so you will just need to annotate your pod to assume the role and then you can simply do this inside:
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

## How to connect to a created Postgres database

If you make a pod like the one above:

```
apiVersion: v1
kind: Pod
metadata:
    name: alexs-test-pod
    annotations:
        iam.amazonaws.com/role: svcop-sandbox-sandbox-gsp-service-operator-test-alexs-test-princ
spec:
    containers:
    -   name: myapp-container
        image: governmentpaas/psql
        command: ['sleep', '1000000']
        volumeMounts:
        -   name: secrets
            mountPath: /secrets
    volumes:
    -   name: secrets
        secret:
            secretName: alexs-test-db-secret
```

You will be able to access the login details under /secrets/Endpoint, /secrets/Port, /secrets/Username and /secrets/Password:

```
/ # cat /secrets/Password
[redacted]
/ # psql -h$(cat /secrets/Endpoint) -p$(cat /secrets/Port) -U$(cat /secrets/Username) postgres
Password for user [redacted]:
psql (11.5, server 10.7)
SSL connection (protocol: TLSv1.2, cipher: ECDHE-RSA-AES256-GCM-SHA384, bits: 256, compression: off)
Type "help" for help.

postgres=>
```

## How it works
You don't need to know this to use it, this information is for cluster operators.
GSP Service Operator consists of a container that runs essentially a daemon, and a kubeyaml config that sets up the container, provides a bunch of custom resource definitions (e.g., there is a definition in there for SQS Queues), etc. - it also gives the container access to interact with the cluster.
The daemon monitors the k8s cluster for such custom resources being created and will create the requested SQS/Database resources.
