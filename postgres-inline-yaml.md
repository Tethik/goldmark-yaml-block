---
inherits: ["sql"]
related:
  - RDS
---

# Postgres

Postgres is a very commonly used database. You often find it hosted using [RDS](./RDS.md).

```yml control
slug: sslmode-full-verify
title: SSLMODE Full Verify
description:
  Per default, some postgres clients don't properly verify the TLS/SSL connection to postgres.
  Ensure that clients (both services and developers) have their SSLMODE set to full-verify.
```

```yml threat
slug: shared-password
title: Shared Password
description: Avoid using the same password for multiple services and/or users. In case of an incident, it will be difficult to track down where the credentials were compromised and shared credentials are more likely to leak.
```

```yml control
title: Role-based authorization
description:
  Use Postgres' built-in role-based authorization system to tightly control which tables a database user can
  access. E.g. if you have a service that only does read-only work, then you can restrict it with a read-only role.
```

```yml control
title: Unique Users
description: Developers and services should use a unique user each, in order to ensure that audit logs can track what user performed what query, to avoid credential reuse, and to better be able to lock down authorization per user/service.
mitigates:
  - postgres/shared-password
```
