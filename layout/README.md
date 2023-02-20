# Architecture diagram layout engines

Maintaining architecture diagrams by hand costs a lot of time.
There are tools, like [GraphViz](https://graphviz.org/), that can automatically lay out a diagram for you.
The results of these _auto-layouts_ are never quite satisfactory, however.

Laying out diagrams is an instance of the more generic problem of
[graph drawing](https://en.wikipedia.org/wiki/Graph_drawing).
We can, therefore, use generic algorithms, like
[force-directed graph drawing](https://en.wikipedia.org/wiki/Force-directed_graph_drawing) to automatically layout out
diagrams.
This is what programs like GraphViz do.
Sometimes they even support [more than one algorithm](https://graphviz.org/docs/layouts/).
These algorithms are also known as _layout engines_.

The problem is that these generic layout engines don't give good results:

- They have more crossing lines than necessary.
- They fail to highlight relationships between services that are important for understanding the architecture.
  Examples are the [BFF](https://samnewman.io/patterns/architectural/bff/) and
  [Event-Driven Architecture](https://en.wikipedia.org/wiki/Event-driven_architecture#Event_channel) pattern.
- They look, well, not like something a human would draw.

This begs the question: **Do we have to choose between costly, beautiful diagrams and cheap, messy ones, or can we have
the best of both worlds?**
At this point, I honestly don't know, but I want to find out!

If it turns out there is no good way to automatically lay out diagrams, then we need to use a hybrid approach:

1. First automatically lay out the diagram using an algorithm that doesn't suck too bad.
2. Then manually improve the layout.

But let's not give up hope just yet.


## Architecture diagram specifics

If generic algorithms don't give good results, then our only hope seems to be something that's specific to architecture
diagrams, which benefits from its distinct properties.
So what are those properties?

- Architecture diagrams are usually moderate in size, on the order of 10-100 nodes.
  Anything bigger than that and people start grouping services to preserve understandability.
- Not all nodes in the graph are created equal.
  We have users, external systems, services, databases, and queues.
  Depending on the architecture, there are specific groupings we expect to see between these, like one database per
  service in a microservice architecture, or one service writing to and one service reading from a database in
  [CQRS](https://www.martinfowler.com/bliki/CQRS.html).
  And depending on the type of diagram, we may separate some nodes from others, e.g. users and external systems on the
  outside and services grouped together on the inside in a container diagram.
