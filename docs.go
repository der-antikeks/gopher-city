package main

/*
	Residential
		provide Worker, Consumer
		require Job, Market

	Commercial
		provide Job, Market
		require Worker, Consumer, Freight

	Industrial
		provide Job, Freight
		require Worker


                +-----------+
        +------->Residential|
        |       +----+------+
        |            |
     +--+--+         |           +---------+
     |Goods|         |           |Resources|
     +-^---+         |           +----+----+
       |             |                |
       |          +--v---+            |
       |          |Gopher|            |
       |          +-+--+-+            |
       |            |  |              |
+------+---+        |  |         +----v-----+
|Commercial<--------+  +--------->Industrial|
+----^-----+                     +----+-----+
     |                                |
     |           +--------+           |
     +-----------+Products<-----------+
                 +--------+

Low Building > 1x1   > 4 Gopher
Mid Building > 2x2x2 > 32 Gopher
Hig Building > 3x3x3 > 108 Gopher

1 Gopher > -0,4 Goods
1 Good > -1,5 Gopher, -0,5 Products
1 Product > -2 Gopher

5 Low Residential > 20 Gopher, -8 Goods
3 Low Commercial > -12 Gopher, +8 Goods, -4 Products
2 Low Industrial > -8  Gopher,           +4 Products

*/
