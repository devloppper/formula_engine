package formula_engine

/*
	SUM(1 + SUM(2,3), SUM(2,3)) + (3 + 4) * 2
-->运算树
	-ROOT
   --SUM;			+; (;3;+;4;);+;2
  ---1;+;SUM;,;SUM;
 ----2; 3 | 2 ,3
*/

type treeNode struct {
	*unit
	parent *unit
}

func newUnitTree() {

}
