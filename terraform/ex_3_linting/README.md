Using [https://github.com/fugue/regula](regula) for linting.

$ ./regula/bin/regula run fatal
-> fatal error: policy with root access
$ ./regula/bin/regula run medium
-> error: overly permissive policy
$ ./regula/bin/regula run ok
-> no rule should fail
