Let n be the size of the board. We introduce 4 additional
nodes r1, r2, b1, b2 calle win nodes. r1 is connected to 
every node in row 0, r2 is connected to every node in row n,
b1 is connected to every node in column 0, b2 is connected
to every node in column n. then, a player has won if and only 
if there is a monochromatic path between two win nodes.

--

We start by defining every node as its own parent. Formally, we have
a digraph G where each vertex is an isolated node. For a node x, we
write p(x) for the node x points to (p(x) = x if x is isolated).

An element x is a root if p(x) = x, we write root(x) = x. If p(x) != x,
the root of x is defined recursively as root(x) = root(p(x)). We say two
elements (x,y) are in the same group iff root(x) = root(y).

Suppose we are given an arbitrary game state and a node i is colored in.
For each neighbor j of i with color(i) = color(j), we call union(i, j).
This method adds an arc from root(i) to root(j).

To check if red won, we just need to check if root(r1) = root(r2), and the
same for blue with b1 and b2.

We represent this digraph with an array A. We set A[i] = p(i). The find method
will find the root of the element i.