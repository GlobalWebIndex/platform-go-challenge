
#
# source this file with . tools.sh
#

fserver="http://localhost:8080"
aserver="http://localhost:8080"
token="./token.jwt"

gwifav-login() {
	curl -s -XPOST -d"email=$1&passwd=$2" ${fserver}/login > ${token}
}

gwifav-add() {
	curl <<EOF --no-progress-meter -w '%header{Location}\n' -H "Authorization: Bearer `cat ${token}`" -d@- -XPOST ${fserver}/favorites
{
  "assetKey": "$1",
  "description": "${2:-description of ${1}}"
}
EOF
}

gwifav-get() {
	curl --no-progress-meter -H "Authorization: Bearer `cat ${token}`" ${fserver}/favorites/$1
}

gwifav-getall() {
	curl --no-progress-meter -H "Authorization: Bearer `cat ${token}`" ${fserver}/favorites
}

gwifav-del() {
	curl --no-progress-meter -H "Authorization: Bearer `cat ${token}`" -XDELETE ${fserver}/favorites/$1
}

gwifav-update() {
	curl <<EOF --no-progress-meter -H "Authorization: Bearer `cat ${token}`" -d@- -XPUT ${fserver}/favorites/$1
{
  "description": "${2:-updated description}"
}
EOF
}

gwifav-blast() {
	assets=(charts insights audiences)
	declare -i i=0 n=${1:-20}
	for (( i=0; i<n; i++)) ;do
		key="/${assets[$(($RANDOM % ${#assets[@]}))]}/$(($RANDOM % 1000))"
		descr="initial description of ${key}"
		gwifav-add "$key" "$descr"
	done
}

gwifav-asset() {
	curl --no-progress-meter -H "Authorization: Bearer `cat ${token}`" ${aserver}/assets/$1
}
