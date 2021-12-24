host="$1"
shift
cmd="$@"
  
# until ! mysqladmin ping -h"${MYSQL_HOST}" --silent; do
#   >&2  echo "Mysql is unavailable - sleeping"
#   sleep 1
# done
 sleep 1
>&2 echo 'Runing migrations...'
exec $cmd

