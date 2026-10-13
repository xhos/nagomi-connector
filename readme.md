# nagomi-connector

Syncs transactions from [Wise](https://wise.com) and [SnapTrade](https://snaptrade.com) into [nagomi](https://github.com/xhos/nagomi).

Accounts are connected in nagomi-web, and [nagomi-core](https://github.com/xhos/nagomi-core) keeps the credentials encrypted. The connector polls core for sync jobs, fetches new transactions and sends them back.

## config

See [.env.example](.env.example).

## dev

```bash
run
tst
cover
regen
bump-protos
```

To add a provider, implement `provider.Provider` under `internal/provider/` and add a case to `providerFactory` in `cmd/server/main.go`.
