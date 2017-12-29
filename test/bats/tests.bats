#!/usr/bin/env bats

treasury=$PWD/treasury
randomKey=$(cat /dev/urandom | env LC_CTYPE=C tr -dc a-zA-Z0-9 | head -c 16)

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
  run $treasury write development/treasury/key1 secret1
  [ $status -eq 0 ]
  [[ ${lines[0]} =~ "Success!" ]]
}

@test "write second" {
  run $treasury write development/treasury/key2 secret2
  [ $status -eq 0 ]
  [[ ${lines[0]} =~ "Success!" ]]
}

@test "write second not forced" {
  run $treasury write development/treasury/key2 secret2
  [ $status -eq 0 ]
  [[ ${lines[0]} =~ "Success!" ]]
}

@test "write second forced" {
  run $treasury write development/treasury/key2 secret2 --force
  [ $status -eq 0 ]
  [[ ${lines[0]} =~ "Success!" ]]
}

@test "write-wrong-data" {
  run $treasury write test secret1
  [ $status -eq 255 ]
  [[ ${lines[0]} =~ "Error" ]]
}

@test "write random key" {
  run $treasury write development/application/"${randomKey}" secret
  [ $status -eq 0 ]
  [[ ${lines[0]} =~ "Success!" ]]
}

@test "read" {
  run $treasury read development/treasury/key1
  [ $status -eq 0 ]
  [[ ${lines[0]} =~ "secret" ]]
}

@test "read-wrong-data" {
  run $treasury read test
  [ $status -eq 255 ]
  [[ ${lines[0]} =~ "Error" ]]
}

@test "export single" {
  run $treasury export development/treasury/key1
  [ $status -eq 0 ]
  [[ ${lines[0]} == "export key1='secret1'" ]]
}

@test "export all" {
  run $treasury export development/treasury/
  [ $status -eq 0 ]
  echo ${lines[0]}
  [[ ${lines[0]} == "export key1='secret1'" ]]
  [[ ${lines[1]} == "export key2='secret2'" ]]
}

@test "import forced" {
  run $treasury import development/treasury/ test/bats/bats.env.test --force
  [ $status -eq 0 ]
  [[ ${lines[0]} == "Import successful" ]]
}

@test "read imported key3" {
  run $treasury read development/treasury/key3
  [ $status -eq 0 ]
  [[ ${lines[0]} =~ "secret3" ]]
}

@test "read imported key4" {
  run $treasury read development/treasury/key4
  [ $status -eq 0 ]
  [[ ${lines[0]} =~ "secret4" ]]
}

@test "template" {
  run $treasury template --src test/resources/source.secret.tpl --dst test/output/bats-output.secret
  [ $status -eq 0 ]
  [[ ${lines[0]} == "File with secrets successfully generated" ]]
}

@test "check version" {
  run $treasury version
  [ $status -eq 0 ]
}
