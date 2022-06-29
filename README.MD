## Platform Engineer Exam

### Setup Requirements
1. Docker
2. Port **8081** is free
3. TCP Port **3337** is free

### Boot Instructions
1. sh start.sh
2. Basic auth creds:
   1. Username: "admin@gmail.com"
   2. Password: "123456"

### Some Notes
1. I did not use for loops for my test code deviating from how it's usually done, because the mock libraries I am using to test won't be easily readable if I go that route
2. I also served a UI from: http://localhost:8081
3. Docs is served from: http://localhost:8081/docs
4. Architecture followed is dependency injection
5. Notable libraries used:
   1. [Fiber](https://github.com/gofiber/fiber) - Multiplexer
   2. [Dingo](https://github.com/sarulabs/dingo) - Dependency Injection Library
   3. [Counterfeiter](https://github.com/maxbrunsfeld/counterfeiter) - Interface Mock
   4. [SQLBoiler](https://github.com/volatiletech/sqlboiler) - ORM
   5. [Go-SQLMock](https://github.com/DATA-DOG/go-sqlmock) - SQL Mock 
   6. [Swaggo](https://github.com/swaggo/swag) - Doc Generator
   7. [Logrus](https://github.com/sirupsen/logrus) - Logger
   8. [Viper](https://github.com/spf13/viper) - Env Parser
   9. [Validator](https://github.com/go-playground/validator) - Validation
   10. [Testify](https://github.com/stretchr/testify) - Testing