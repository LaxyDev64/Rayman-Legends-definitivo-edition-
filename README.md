<h1 align="center">Rayman-Legends-definitive-edition</h1>

<p align="center"><b>Nextendo Network game server for Rayman Legends definitivo edition.</b></p>

<p align="center">
  <img src="https://img.shields.io/badge/license-PolyForm%20Shield%201.0.0-orange" alt="License">
  <img src="https://img.shields.io/badge/go-1.23%2B-00ADD8" alt="Go 1.23+">
</p>

---

## What is this?

The NEX game server for **Super Smash Bros. Ultimate** on [Nextendo Network](https://nextendo.network)
— authentication, matchmaking, NAT traversal, and the online-init handlers (DataStore `0x73` /
Utility `0x6E`) that bring SSBU's online mode up. Built on the
[**nextendo-nex**](https://github.com/NextendoNetwork/nextendo-nex) core.

> **DataStore / Utility init is a stub.** SSBU calls a set of DataStore and Utility getters during
> the online bring-up. This server answers them with a minimal valid response (empty list) so the
> game proceeds; a full server-side implementation is not yet part of this tree, so online play is
> not complete out of the box.

## Running

```sh
cp example.env .env   # then edit .env
go run .
```

No secrets, keys, or measured data are baked into the source — everything comes from the environment.

## What this is not

Ships **no** Nintendo code, keys, measured data, or copyrighted assets. Independent reimplementation
for a community-run service; not affiliated with Nintendo. The NEX access key is a public per-title
value, not a secret.

## License

**[PolyForm Shield License 1.0.0](LICENSE.md)** — source-available.
