# Helpful Info

- When using the the zap logger this format can used for derived fields in grafana to link logs to tempo : "trace_id":"([a-f0-9]{32})"
- When using the the default logger (with the DefaultLogFormat) this format can used for derived fields in grafana to link logs to tempo : trace_id ([a-f0-9]+)
- You could use this regex to account for both of the above cases : "trace_id":"([a-f0-9]{32})"|trace_id\s+([a-f0-9]+)
