# k8s

## import GPG keys

```

kubectl -n kube-system create secret generic pkg-distributor-keys  \
    --from-file=public.key=/path/to/public.key \
    --from-file=private.key=/path/to/private.key
```
