#show setColor/2.
1 {setColor(V,C): color(C)} 1 :- vertex(V).
:- setColor(V1,C), setColor(V2,C), edge(V1, V2), V1 != V2.