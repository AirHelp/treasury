#!/usr/bin/env bats

treasury=$PWD/treasury

@test "Check that the treasury binary is available" {
    command $treasury
}

@test "usage" {
  run $treasury
  [ $status -eq 0 ]
}

@test "help" {
  run $treasury --help
  [ $status -eq 0 ]
}

# @test "version" {
#     run $treasury --version
#     [ $status -eq 0 ]
#     [[ ${lines[0]} =~ "treasury version" ]]
# }

@test "write" {
  run $treasury write test/treasury/test-key secret
  [ $status -eq 0 ]
  [[ ${lines[0]} =~ "Success!" ]]
}

@test "write second" {
  run $treasury write test/treasury/test-key2 secret2
  [ $status -eq 0 ]
  [[ ${lines[0]} =~ "Success!" ]]
}

@test "write-wrong-data" {
  run $treasury write test secret
  [ $status -eq 255 ]
  [[ ${lines[0]} =~ "Error" ]]
}

@test "read" {
  run $treasury read test/treasury/test-key
  [ $status -eq 0 ]
  [[ ${lines[0]} =~ "secret" ]]
}

@test "read-wrong-data" {
  run $treasury read test
  [ $status -eq 255 ]
  [[ ${lines[0]} =~ "Error" ]]
}

@test "export single" {
  run $treasury export test/treasury/test-key
  [ $status -eq 0 ]
  [[ ${lines[0]} == "export test-key='secret'" ]]
}

@test "export all" {
  run $treasury export test/treasury/
  [ $status -eq 0 ]
  [[ ${lines[0]} == "export test-key='secret'" && ${lines[1]} == "export test-key2='secret2'" ]]
}
