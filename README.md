# BuPol examination preparation tool
Super simple website for certain examination preparation. Started as a quick and simple tool for a friend to prepare
for the examination.

There are two exercises with the goal of memorizing the content of the columns/rowns. The actual content (words) are tailored 
for the examination and doesn't make much sense without knowing the context (though it also doesn't really matter in the end).

The latest version uses jsGrid which allows better user interaction instead of just showing static tables. The user can
decide the row count (the actual amount of row to memorize) and can set up a timer interval for how long an exercise/table is 
visible (time to memorize). The results can be checked row by row entering the actual value for the second column of a row. 
There is always the possibility to re-generate the content and to skip the querying of the results. After the results page
for the second exercise you get re-directed to start again with the first exercise.

The Latest version can be found here: http://3.67.79.130:8889/bupol/exercises/first/

Will use this project to firm known tech stacks but also to learn and explore new things. Esp. when it comes to web based
frontends :) Some ideas which came to my mind:
* Getting proficient with Golang
* HTTP session management (https://github.com/gorilla/sessions) and/or JWT based auth
* Some kind of CRUD "microservice" to handle progress of the user (SQL based DB, message broker like Rabbit or Kafka)
* Terraform for AWS resource management
