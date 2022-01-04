# README

⚠️ WIP - Work in progress

This is an extension to [goldmark markdown](https://github.com/yuin/goldmark) that enables custom parsing logic for threats/controls, to be used in a documentation system. Heavily based on [goldmark-meta](https://github.com/yuin/goldmark-meta).

It overrides the normal code-block parsing when detecting `yml control` or `yml threat` and
decodes the yaml code within.

````markdown
```yml control
slug: sslmode-full-verify
title: SSLMODE Full Verify
description:
Per default, some postgres clients don't properly verify the TLS/SSL connection to postgres.
Ensure that clients (both services and developers) have their SSLMODE set to full-verify.
```
````

````markdown
```yml threat
slug: shared-password
title: Shared Password
description: Avoid using the same password for multiple services and/or users. In case of an incident, it will be difficult to track down where the credentials were compromised and shared credentials are more likely to leak.
```
````
