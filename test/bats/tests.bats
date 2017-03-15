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
  run $treasury write test/treasury/key1 secret1
  [ $status -eq 0 ]
  [[ ${lines[0]} =~ "Success!" ]]
}

@test "write second" {
  run $treasury write test/treasury/key2 secret2
  [ $status -eq 0 ]
  [[ ${lines[0]} =~ "Success!" ]]
}

@test "write second not forced" {
  run $treasury write test/treasury/key2 secret2
  [ $status -eq 0 ]
  [[ ${lines[0]} =~ "Success!" ]]
}

@test "write second forced" {
  run $treasury write test/treasury/key2 secret2 --force
  [ $status -eq 0 ]
  [[ ${lines[0]} =~ "Success!" ]]
}

@test "write-wrong-data" {
  run $treasury write test secret1
  [ $status -eq 255 ]
  [[ ${lines[0]} =~ "Error" ]]
}

@test "read" {
  run $treasury read test/treasury/key1
  [ $status -eq 0 ]
  [[ ${lines[0]} =~ "secret" ]]
}

@test "read-wrong-data" {
  run $treasury read test
  [ $status -eq 255 ]
  [[ ${lines[0]} =~ "Error" ]]
}

@test "export single" {
  run $treasury export test/treasury/key1
  [ $status -eq 0 ]
  [[ ${lines[0]} == "export key1='secret1'" ]]
}

@test "export all" {
  run $treasury export test/treasury/
  [ $status -eq 0 ]
  echo ${lines[0]}
  [[ ${lines[0]} == "export key1='secret1'" ]]
  [[ ${lines[1]} == "export key2='secret2'" ]]
}

@test "import forced" {
  run $treasury import test/treasury/ test/bats/bats.env.test --force
  [ $status -eq 0 ]
  [[ ${lines[0]} =~ "Success!" ]]
  [[ ${lines[1]} =~ "Success!" ]]
  [[ ${lines[2]} == "Import successful" ]]
}

@test "read imported key3" {
  run $treasury read test/treasury/key3
  [ $status -eq 0 ]
  [[ ${lines[0]} =~ "secret3" ]]
}

@test "read imported key4" {
  run $treasury read test/treasury/key4
  [ $status -eq 0 ]
  [[ ${lines[0]} =~ "secret4" ]]
}
