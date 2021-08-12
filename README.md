## 단축 URL 제공 서비스

1. 기본 세팅
    - main() 실행 후 http://localhost:8080 접속 확인
    - DB는 Golang 에서 제공하는 자료구조를 이용 ( array, slice, map ...)
    - Web Framework 사용에는 제한이 없음 (Beego, Gin, Gorilla ...)
2. 요구 사항
    1. URL(이하 oriUrl) 을 요청 받아 짧은 URL(이하 shortUrl) 을 응답한다.
    2. shortUrl 은 8 Character 이내로 생성한다.
    3. shortUrl 은 중복되지 않는다.
    4. 동일한 oriUrl 요청은 동일한 shortUrl 을 응답한다.
    5. shortUrl 을 요청 받아 원래 oriUrl 을 응답한다.
    6. shortUrl 요청 수를 기록한다.
    7. 등록된 oriUrl, shortUrl, 요청 수를 조회 할 수 있다.