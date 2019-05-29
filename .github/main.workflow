workflow "Run Tests" {
  resolves = "Test"
  on = "push"
}

action "Test" {
  uses = "docker://golang"
  runs = "make test"
}