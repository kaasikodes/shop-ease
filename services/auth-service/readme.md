# Details

The service serves the primary purpose of keeping/tracking user records and implements RBAC. The identified roles are vendor, customer, and admin. A user can be all 3 at once a vendor is responsible for putting out products to be purchased by customers although the vendors are limited by their subscription package to

- Identified tables - token, user, role, user_role (all tables have createdAt & updatedAt)
- token - id, value, type(ACCOUNT_ACTIVATION, PASSWORD_RECOVERY, AUTH_TOKEN), expires_at
- user - name(encrypt), email, pwd(encrypt), id, phone(encrypt)
- role - id, name(enum)
- user_role - id, user_id, role_id, status(ACTIVE, INACTIVE)
- The auth token should contain all the needed details of user(or should be available in redis with the expire set to that of the token) to ensure that subsequent calls from the user don't have to make calls to auth service (if down) and can just continue with the details the user is already aware of

## Flows

- On first registration all users will be customers
- Users can also specifically register to be vendors, in wish case he will first interact with the subscription service after which interacts with payment service after which payment is made webhook is triggered to inform auth to activate the vendor role, after which they are notified and have access to the vendor service to create/update **store**, manage orders, update inventories, etc.
- Users cannot register as admins but rather have to be added to the system as admins (who can view vendor activity, store items, but not modify products, or orders that vendors are responsible for)

## TODO

This what is expected

- Health Checker[Done]
- Logger[Done]
- Error Setup [Done]
- Database setup, model/repo setup, creation of handlers, config
- Validation Middleware
- Using air for hot reload
- Makefile
- Documenting endpoints with swagger
- Docker, n GRPC
- GRPC

# Study Up

Need to sudy up on db connections and houw they are reused, why the need disconnections what does row.Close() do exactly, why does not closing lead to a memory leak, and why is considered a memory leak

# Tools

The following are the cli tools installed either by chocolatey or scoop being used, also I will recommend you install scoop and cholatey if your windows like me, scoop is more developer suited though from what I have been able to surmize.

- Make
- migrate
- direnv

# Useful DB Commands

<!-- CREATE USER 'admin_user'@'%' IDENTIFIED BY 'StrongPassword'; -->

CREATE USER 'admin1'@'%' IDENTIFIED BY 'root123';

<!-- GRANT ALL PRIVILEGES ON *.* TO 'admin_user'@'%' WITH GRANT OPTION; -->

GRANT ALL PRIVILEGES ON _._ TO 'admin1'@'%' WITH GRANT OPTION;

<!-- FLUSH PRIVILEGES; -->

FLUSH PRIVILEGES;

<!-- SHOW GRANTS FOR 'admin_user'@'%'; -->

SHOW GRANTS FOR 'admin1'@'%';
