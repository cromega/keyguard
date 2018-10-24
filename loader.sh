echo -n "OTP: "
read password < /dev/tty

mktemp="mktemp"
if command -v gmktemp; then
  mktemp="gmktemp"
fi

keyfile=$($mktemp)
trap "rm -rf $keyfile" EXIT

curl -k -f -s -u "keyguard:$password" "{{ .URL }}" > "$keyfile"
if [ $? ]; then
  ssh-add -t {{ .Expiry }} "$keyfile"
else
  echo "something went wrong."
fi
