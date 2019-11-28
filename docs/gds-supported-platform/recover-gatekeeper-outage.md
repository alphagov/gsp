# Recover from a gatekeeper failure

## Symptoms

Several pods in the system will be stuck in a "Terminating" state, the
gatekeeper among them. 

## Causes

The node that the gatekeeper pod was running on was terminated in a way that
meant the control plane could not reschedule the pod elsewhere. All pods created
after the gatekeeper comes online will have the gatekeeper finalizer added and
this finalizer can not run once the gatekeeper pod is down in this way. This
prevents other pods from terminating properly.

## Corrective actions

1. Elevate to admin in the affected cluster
1. Execute the following command:
   ```
   kubectl -n gsp-system patch pod gatekeeper-controller-manager-0 -p '{"metadata":{"finalizers": []}}'
   ```

That should enable the gatekeeper pod to terminate and reschedule.
