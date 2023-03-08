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
  [Event-Driven Architecture](https://en.wikipedia.org/wiki/Event-driven_architecture#Event_channel) patterns.
- They look, well, not like something an architect would draw by hand.

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
- Nodes in a diagram have shapes with non-negligible sizes.
  Text is usually placed inside the node's shape rather than next to the node.
- Diagrams usually lay out nodes in a grid and connect them with orthogonal edges, not straight lines.
- Not all nodes in the graph are created equal.
  We have users, external systems, services, databases, and queues.
  Depending on the type of diagram, we may separate nodes of different types, e.g. a container diagram would show
  users and external systems on the outside and services grouped together on the inside.
- Depending on the architecture, there are specific groupings we expect to see between the node types, like one
  database per service in a microservice architecture, or one service writing to and one service reading from a database
  in [CQRS](https://www.martinfowler.com/bliki/CQRS.html).

Let's start with the moderate size.
Even the lower bound of 10 is problematic.
Let's assume for a moment that all we have to do is line up the nodes in one straight line.
This is obviously an over-simplification, the real problem is much harder.
Even for this simple layout, there are 10! = 3.6M possibilities.
If we process one possibility per ms, then we'd need about an hour.
For larger diagrams, this number increases exponentially.
So obviously this isn't going to work.

If the search space is too big to search exhaustively, then a couple of options remain:

1. Prune the search space to make it smaller by quickly rejecting significant parts of it.
  Then comprehensively search the remainder of the space.
2. Use heuristics.
3. Use a local search algorithm, like Simulated Annealing (SA) or Genetic Algorithms (GA).

I'm not aware of any techniques that work well using the first or second approaches.
The first one doesn't even sound very promising.
The second one would require very deep insight into the problem space and even then isn't guaranteed to work.

That leaves the third option.
SA and GA approaches are similar in spirit, but technically different.
I have prior experience with GAs, so let me start there.

A quick Google search shows that many people have attempted to automatically lay out graphs using Genetic Algorithms.
Most of this work is on graphs rather than diagrams and most present results that wouldn't be satisfactory for diagrams.
But some seem to produce decent results, so this sounds like a promising approach to try.
