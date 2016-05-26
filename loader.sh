echo -n "OTP: "
read password < /dev/tty

mktemp="mktemp"
if command -v gmktemp; then
  mktemp="gmktemp"
fi

keyfile=$($mktemp)
trap "rm -rf "$keyfile"" EXIT

curl -k -f -s -u "cromega:$password" "{{ .Url }}" > "$keyfile"
ret=$?
if [ $ret -eq 0 ]; then
  ssh-add -t 32400 "$keyfile"
else
  echo "something went wrong."
fi
