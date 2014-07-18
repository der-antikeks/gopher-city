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

*/
