# GDS Supported Platform Continuous Deployment

<!--__[edit draw.io diagram](https://www.draw.io/?state=%7B%22ids%22:%5B%221p4qkP2-fsMpc42ZCsKIVlGwJ4blmnu_q%22%5D,%22action%22:%22open%22,%22userId%22:%22104206899246339571570%22%7D#G1p4qkP2-fsMpc42ZCsKIVlGwJ4blmnu_q)__
-->

![GDS Supported Platform Continuous Delivery](diagrams/gsp-architecture-continuous-delivery.svg)

|Seq|Step|Description|
|-|-|-|
|1.|Raise pull request|The Developer raises a pull request to review changes to the source code repository|
|2.|Merge pull request to master|Require a minimum of two developers to approve the change before allowing the merge to the master branch|
|3.|Check change to master|Determine if the master branch has changed and needs building|
|4.|Check for 2-eyes| Ensure that the change has been signed by two different developers|
|5.|Check signers on approved list|Check the change has been signed by approved developers|
|6.|Trigger build task|Start the pipeline to build from the master branch|
|7.|Build images|Build the Docker images|
|8.|Run unit tests|Unit test the Docker images|
|9.|Push images |Push successfully built Docker images to the Docker Registry|
|10.|Tag and sign Docker images|Tag the newly built Docker images and sign them|
|11.|Bump image tag|Increment the tag on the Docker image|
|12.|Check for deploy|Determine if a deployment needs to be applied to the cluster|
|13.|Trigger deploy|Start the deployment|
|14.|Render deploy template |Merge the secrets and variables with the deployment template|
|15.|Apply deploy|Apply the deployment to the cluster|
|16.|Pull images|Pull images referred to in the cluster from the private Docker Registry|
|17.|Return image|Registry returned the Docker image|
|18.|Wait for healthy state|Wait until the deployment has been applied to the cluster and the deployed application is in a healthy state|
|19.|Run acceptance tests|Run any post installation tests|
|20.|Tag and sign deployment|Tag and sign the successfully applied deployment to mark it OK for promotion to next step (e.g. production)|
