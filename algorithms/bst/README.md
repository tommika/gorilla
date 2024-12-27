BST
===

Generic binary search tree data structure and algorithms

The bst package implements a binary search tree, with the following
associative-array (map) operations:
* Get
* Put
* Delete
* Size

and the following ordered-collection operations:
* VisitInOrder
* Min
* Max

Trees can be either unbalanced or balanced. Balanced trees are implemented
using a Red-Black tree algorithm.

This implementation does not maintain parent pointers within the nodes of the
tree.  As such, many operations (in particular re-balancing after insertion and
deletion) must remember the path from the root to an impacted node for the
duration of the operation.