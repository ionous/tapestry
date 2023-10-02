README.MD
---------

specifying `use_mingw=true` because i had problems using the msvc toolchain to compile the go module.

**building**


`scons use_mingw=true`


**prereqs**


* go version 1.21  ( needed for pinning string pointers )

* in the `godot=cpp` directory:
```
scons use_mingw=true custom_api_file=.\gdextension\extension_api.json
```

* installing the vcredist might be necessary; i had some loading problems initially without it.


