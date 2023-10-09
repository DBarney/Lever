#Lever - YAHR (yet another http router)

```
 /\
 \ \  "Give me a place to stand and with a lever I will move the whole world."
  \ \           - Archimedes
   \ \
    \ \
     \ \
      \ \
       \ \
        \/
```


lever gives developers a very good place to stand. It is a minimal http router that supports middleware.

## why? oh why? oh why?

Well honestly I like small libraries that do one thing and do it well. I feel like a lot of solutions are massive compared to the functionality they provide. 
I also kept running into the same problem and tried to find a solution to it. The problem was passing additional data from middle wares into the `http.HandlerFunc` that defined the endpoints.
