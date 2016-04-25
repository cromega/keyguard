echo -n "OTP: "
read password < /dev/tty

keyfile=$(mktemp)
curl -f -s -u "cromega:$password" "{{ .Url }}" > "$keyfile"
ret=$?
if [ $ret -eq 0 ]; then
  trap "rm -rf "$keyfile"" EXIT
  ssh-add -t 32400 "$keyfile"
else
  echo "something went wrong."
fi
