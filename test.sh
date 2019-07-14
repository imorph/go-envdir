# This is envdir actual tests
# from http://cr.yp.to/daemontools/install.html
# changed slightly for go-envdir
# vorks only on linux file systems (ln required)

echo '--- go-envdir requires arguments'
./go-envdir whatever; echo $?

echo '--- go-envdir complains if it cannot read directory'
ln -s test/env1 test/env1
./go-envdir test/env1 echo yes; echo $?

echo '--- go-envdir complains if it cannot read file'
rm test/env1
mkdir test/env1
ln -s Message test/env1/Message
./go-envdir test/env1 echo yes; echo $?

echo '--- go-envdir adds variables'
rm test/env1/Message
echo This is a test. This is only a test. > test/env1/Message
./go-envdir test/env1 sh -c 'echo $Message'; echo $?

echo '--- go-envdir removes variables'
mkdir test/env2
touch test/env2/Message
./go-envdir test/env1 ./go-envdir test/env2 sh -c 'echo $Message'; echo $?

echo '--- go-envdir trims tabs and spaces'
echo "trim test                " > test/env1/Message1
./go-envdir test/env1 sh -c 'echo "${Message1}" end'; echo $?

echo "--- go-envdir changes \0 to \n"
echo "see ya\0later!\0bye!" > test/env1/Message2
./go-envdir test/env1 sh -c 'echo "${Message2}"'; echo $?

echo '--- go-envdir returns child comand exit code: 42'
./go-envdir test/env1 test/return-42.sh; echo $?

rm -rf test/env1 test/env2