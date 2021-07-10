# OpenShift vs. Rancher

- OpenShift CLI introduces `oc`, which is compatible with `kubectl`, but adds a
  few extra special helpers that help to get job done.
- OpenShift is a Certified Kubernetes distribution by the CNCF.
- OpenShift adds value with extra features over kubectl, like the ability to
  create new namespaces, easily switch context, and commands for developers such
  as the ability to build container images and deploy applications directly from
  source code or binaries, also known as the Source-to-image process (s2i).
- That is, `oc project <project>` and `oc get projects`.
- OpenShift ODO is a command-line allowing to quickly deploy local code to a
  remote OpenShift cluster. Making it possible to have an efficient inner loop
  where all changes can instantly be synced with the running container, thus no
  need to rebuild and push the image.
- OpenShift pays a lot of attention for security and quality assurance. "We
  believe in delivering a stable and long-supported platform."
- OpenShift not only delivers Kubernetes, but also all the essential open source
  tools that make it an enterprise-ready solution to confidently run production.
- For developers who still want to run everything on their laptop, they can rely
  on Codeready Containers, which is an all-in-one OpenShift 4 running on the
  laptop.
- Local development using [CodeReady containers][crc] seems to be a great idea.
- Podman is Docker alternative which allows to run several Docker images inside
  a single pod (compare to Kubernetes pod) - local development made easy?
- So far, Red Hat stands out as the market leader with 44 percent market share.
  Rancher is a small player with 3 % share compared to Red Hat.

[crc]: https://developers.redhat.com/products/codeready-containers/overview
