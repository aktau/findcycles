findcycles(1)
=============

findcycles(1) filters a DOT format graph, keeping only nodes that are in in
cycles. Input is taken on stdin, output is provied on stdout.

It proved useful to me while debugging cycles in systemd units. For systemd
graphs, only the After edges are relevant for ordering cycles. You can
filter them out with grep(1) or some other tool to get a more concise output
graph.

To generate a suitable graph from systemd, you can issue a shell command as
described in the following link: http://unix.stackexchange.com/a/227963.

Afterwards, I did the following to study the cycles:

```bash
$ <remount-cycle.dot grep -P 'digraph|\}|green' | findcycles |
    dot -Tsvg > cycles-order.svg
$ <remount-cycle.dot findcycles |
    dot -Tsvg > cycles-all.svg
```

NOTE: findcycles(1) may very well contain bugs when used outside the realm
of systemd graphs, you have been warned
