#!/usr/bin/env bats

treasury=$PWD/treasury
randomKey=$(cat /dev/urandom | env LC_CTYPE=C tr -dc a-zA-Z0-9 | head -c 16)
valid_aws_region=eu-west-1
invalid_aws_region=us-west-1

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

@test "read-with-valid-region" {
  run $treasury read development/treasury/key1 -r $valid_aws_region
  [ $status -eq 0 ]
  [[ ${lines[0]} =~ "secret" ]]
}

@test "read-with-wrong-region" {
  run $treasury read development/treasury/key1 -r $invalid_aws_region
  [ $status -eq 255 ]
  [[ ${lines[0]} =~ "Error:" ]]
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
  run $treasury template --src test/resources/bats-source.secret.tpl --dst test/output/bats-output.secret
  [ $status -eq 0 ]
  [[ ${lines[0]} == "File with secrets successfully generated" ]]
}

@test "template-and-var-append" {
  run $treasury template --src test/resources/bats-source.secret.tpl --dst test/output/bats-output.secret --append 'key1:treasury'
  [ $status -eq 0 ]
  run grep "key1=secret1treasury" test/output/bats-output.secret
  [ $status -eq 0 ]
}   

@test "template-and-var-append-multiple-variables" {
  run $treasury template --src test/resources/bats-source.secret.tpl --dst test/output/bats-output.secret --append 'key1:treasury' --append 'key2:?pool=20'
  [ $status -eq 0 ]
  run grep "key1=secret1treasury" test/output/bats-output.secret
  [ $status -eq 0 ]
  run grep "key2=secret2?pool=20" test/output/bats-output.secret
  [ $status -eq 0 ]
}

@test "template-and-var-append-bad-input" {
  run $treasury template --src test/resources/bats-source.secret.tpl --dst test/output/bats-output.secret --append 'key1::treasury'
  [ $status -eq 0 ]
  run grep "key1=secret1:treasury" test/output/bats-output.secret
  [ $status -eq 0 ]
}

@test "template wrong key" {
  run $treasury template --src test/resources/bats-wrong-source.secret.tpl --dst test/output/bats-output.secret
  [ $status -eq 255 ]
  [[ ${lines[0]} =~ "Error" ]]
}

@test "write file content to treasury key" {
  run $treasury write development/treasury/key5 test/resources/test_file --file
  [ $status -eq 0 ]
  run $treasury read development/treasury/key5 
  [[ ${lines[0]} =~ "H4sIAAAAAAAA/yopSk0sLi2qBAQAAP//MDbE1QgAAAA=" ]]
}

@test "write too large file content to treasury key" {
  run $treasury write development/treasury/key5 test/resources/test_large_file --file
  [ $status -eq 255 ]
}

@test "check version" {
  run $treasury version
  [ $status -eq 0 ]
}
