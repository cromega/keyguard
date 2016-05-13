echo -n "OTP: "
read password < /dev/tty

keyfile=$(mktemp)
trap "rm -rf "$keyfile"" EXIT

curl -k -f -s -u "cromega:$password" "{{ .Url }}" > "$keyfile"
ret=$?
if [ $ret -eq 0 ]; then
  ssh-add -t 32400 "$keyfile"
else
  echo "something went wrong."
fi
