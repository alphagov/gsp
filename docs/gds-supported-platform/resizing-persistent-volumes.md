# Resizing persistent volumes in StatefulSets

You may need to resize your persistent volumes as your needs change.

This document contains a cookbook for how to do this, and some context
on why the process is the way it is.

## Cookbook: resizing persistent volumes in StatefulSets

As an example, we're going to double the space on the Concourse Worker
persistent volumes - from 64GiB to 128GiB. Replace
`gsp-concourse-worker` with the name of the resource you want to
resize the volume for.

1. Make a change to `charts/gsp-system/values.yml` with your desired
   sizing, like [in PR 343](https://github.com/alphagov/gsp/pull/343):

```yaml
...
persistence:
  worker:
    size: 128Gi
```

1. When merged, the PR will be released to sandbox. Changes must pass
   through sandbox via the full release process to make sure we haven't
   broken anything before going to the Verify cluster.

1. You should destroy the StatefulSet when the Concourse release
   pipeline starts to fail. You will see an error message similar to:

    > The StatefulSet “gsp-concourse-worker” is invalid: spec:
    > Forbidden: updates to statefulset spec for fields other than
    > ‘replicas’, ‘template’, and ‘updateStrategy’ are forbidden

   To destroy the StatefulSet, run:

   ```sh
   $ gds sandbox k -n gsp-system delete statefulset --cascade=false gsp-concourse-worker
   ```

   The `--cascade=false` flag makes sure that you do not delete the pods inside the StatefulSet.

   The `-deployer` pipeline will re-apply the StatefulSet. This will
   intentionally roll all the associated pods.

1. Once the pods are in a `Running` state, edit the
   `PersisentVolumeClaim`s manually.

    - Find the name of the claims you want to delete. In this example,
      we want the three `concourse-worker`s.

      ```sh
      $ gds sandbox k -n gsp-system get persistentvolumeclaim
      NAME                                        STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
      concourse-work-dir-gsp-concourse-worker-0   Bound    pvc-a633fb0a-76fe-11e9-bbf3-0ae9bf7a1ac0   64Gi      RWO            gp2            84d
      concourse-work-dir-gsp-concourse-worker-1   Bound    pvc-a6366add-76fe-11e9-bbf3-0ae9bf7a1ac0   64Gi      RWO            gp2            84d
      concourse-work-dir-gsp-concourse-worker-2   Bound    pvc-a6386683-76fe-11e9-bbf3-0ae9bf7a1ac0   64Gi      RWO            gp2            84d
      ```

    - For each claim, run:

      ```sh
      $ gds sandbox k -n gsp-system edit persistentvolumeclaim concourse-work-dir-gsp-concourse-worker-0
      $ gds sandbox k -n gsp-system edit persistentvolumeclaim concourse-work-dir-gsp-concourse-worker-1
      $ gds sandbox k -n gsp-system edit persistentvolumeclaim concourse-work-dir-gsp-concourse-worker-2
      ```

      This will open `$EDITOR`. Update
      `.spec.resources.requests.storage` to the desired size, then
      save.

1. Run `gds sandbox k -n gsp-system get persistentvolumeclaim`
   to check that your volumes have resized as expected. There should
   be the same number of volumes as you resized.

1. Finally, delete each pod in turn (it will be recreated by the
   StatefulSet), allowing the volume resize in AWS to proceed.

   ```sh
   $ gds sandbox k -n gsp-system get pods
   NAME                                                READY   STATUS      RESTARTS   AGE
   gsp-concourse-worker-0                              1/1     Running     0          3m13s
   gsp-concourse-worker-1                              1/1     Running     0          4m20s
   gsp-concourse-worker-2                              1/1     Running     0          5m11s
   ...
   $ gds sandbox k -n gsp-system get persistentvolumeclaims
   NAME                                        STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
   concourse-work-dir-gsp-concourse-worker-0   Bound    pvc-a633fb0a-76fe-11e9-bbf3-0ae9bf7a1ac0   128Gi      RWO            gp2            84d
   concourse-work-dir-gsp-concourse-worker-1   Bound    pvc-a6366add-76fe-11e9-bbf3-0ae9bf7a1ac0   128Gi      RWO            gp2            84d
   concourse-work-dir-gsp-concourse-worker-2   Bound    pvc-a6386683-76fe-11e9-bbf3-0ae9bf7a1ac0   128Gi      RWO            gp2            84d
   ```

1. Repeat steps 2-5 for the Verify cluster, with the additional step
   of making the pre-release a full release in the `alphagov/gsp`
   GitHub Releases page.

## Context: but why is this so complicated?

This is hard because Kubernetes doesn't support resizing volumes
through StatefulSets.  There is [a proposal to support volume
expansion through
StatefulSets](https://github.com/kubernetes/enhancements/issues/661)
but it is still work in progress.  Until then, in order to change the
size of a volume claim template in a StatefulSet, you must destroy and
recreate the set.

But: the volume claim template in the StatefulSet only controls how
*new* volumes are created.  When you destroy and recreate the
StatefulSet, the existing PersistentVolumes do not get resized.  That
is why you must edit the PersistentVolumeClaims and update them
individually.  PersistentVolumeClaims *do* support on-line resizing
(at least, it does under the settings present in GSP -- see
`allowVolumeExpansion` in
[default-storage-class.yaml](../../charts/gsp-cluster/templates/00-aws-auth/default-storage-class.yaml))
but this is much more laborious.
