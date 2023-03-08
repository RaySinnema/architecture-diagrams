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


## Genetic algorithms

A quick Google search shows that several people have attempted to automatically lay out graphs using Genetic Algorithms.
Most of this work is on graphs rather than diagrams and most present results that wouldn't be satisfactory for diagrams.
But some seem to produce decent results, so this sounds like a promising approach to try.

A Genetic Algorithm requires the following:

1. A representation of solutions as a _genome_.
  This genome usually consist of multiple units, the _genes_.
  The GA starts with an initial collection of random genomes, the _pool_.
2. A way to evaluate how well the solution solves the problem.
  This _fitness function_ takes the genome as input and outputs the _fitness_, a number between 0 (bad) and 1 (good).
3. Ways to derive new solutions from existing ones.
  These _genetic operators_ take one or more genomes as input and produce one or more new genomes.
4. A way to select genomes from the pool to form a new pool in a new _iteration_.
  This _selection mechanism_ is usually solution-agnostic and takes only the fitness of genomes into account.
  Examples are roulette wheel selection, rank selection, and tournament selection.
5. A way to determine when to stop iteration.
  This _termination mechanism_ is also usually solution-agnostic, looking only at the fitness of genomes in the pool.
  The simplest termination mechanism is to perform a fixed number of iterations.
  More advanced termination mechanisms look at the (distribution of) fitness of the genomes.

Let's look at these in the context of laying out diagrams.


### The genome

#### Nodes

A genome has to encode the solution, in this case the layout of the nodes and edges in the graph.
If we assume a grid layout, then we can assign integer coordinates to each cell in the grid.
The genome then has to assign nodes to those coordinates.

How big does the grid need to be?

A node in the grid has 8 neighbors.
For highly connected nodes, some of that space can't be used, since it's needed to draw edges.
To be on the safe side, let's use only 1 out of every 3 grid points for nodes.
Then the grid must have at least `3n` points.

The optimal grid is usually not square.
By making one side of the grid longer than the other, we create more room to connect nodes.
But making the grid too short also doesn't work.
In practice, a ratio of about 2:1 seems to work best.

One of the distinct characteristics of a diagram compared to a generic graph is that the nodes in a diagram are
shapes that have a size, whereas nodes in generic graphs are drawn as small circles.
Diagrams need the larger size to place text inside the shape.
Most shapes (e.g. rectangles) are drawn wider than high, although the reverse certainly occurs (e.g. person shapes).
To be visually pleasing, it then makes sense for the diagram itself to be wider than high.

If `w` is the width of the grid and `h` is its height, then this means that we want `w=2h`.
The number of points in the grid is then `wh = 2h² = 3n`, which gives `h=⌈√(3/2·n)⌉`.
For example, if `n=16`, then `h=5` and `w=10`.

Now that we have worked out the coordinates, we need to design the genome so that it assigns those coordinates to the
nodes.
The most straightforward representation lists the coordinates for the nodes one after the other.
In other words, a single gene consists of integer x and y coordinates, where `0 ≤ x ≤ w-1` and `0 ≤ y ≤ h-1`, and a 
genome is a sequence of `n` genes, one per node.

#### Edges

TODO


### The fitness function

The fitness function has to convert a genome into a number between 0 and 1.
This single number has to cover multiple dimensions of "goodness" or "badness" of the layout:

1. The fewest number of crossings, `c`.
  If there are `E` edges, then `0 ≤ c ≤ E-1`, so `Fc = 1 - c / (E-1)` works nicely.
  In order to calculate `c`, we need to know how the edges of the graph go, to see which cross.
  So we either need to add the position of the edges to the genome, or calculate these positions inside the fitness
  function.
2. The relationships between different types of nodes.
  For instance, in a container diagram you usually want to put the personas and external systems at the edge and the
  services, databases, and queues in the middle of the diagram.
3. Symmetric placement of nodes is better than asymmetric placement.
  This holds in both directions.
4. Symmetric placement of connections.
  If nodes have 3 connectors per side, then we prefer to use the middle connector if there is only one edge attached, 
  but we want to use the outer two if there are two edges attached.
5. Non-overlapping connectors.
  If there are more edges connected to a node than it has connectors, then we have to overlap edges onto the same
  connector.
  Otherwise, we prefer not to do that.
6. A smaller size of the grid.
  If instead of a `w·h` grid we can make do with something smaller, then that's better.
  In other words, we prefer solutions where entire rows or columns are unused.
7. We like the visual expression of architectural patterns.
  For example, if the architecture is based on microservices, each of which has their own database, then we'd like to
  see the same spatial relationship between the service and its database everywhere.
  This depends on being able to recognize those patterns, of course.

TODO: Finish and combine into one formula.


### The generic operators

TODO
