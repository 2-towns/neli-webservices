#!/bin/bash
docker exec mysql sh -c "mysql -u $USERNAME -p$PASSWORD neli_webservices -e'SELECT id, DATE(date) FROM cleaning WHERE done = 0' -s" | while read cleaning_id date; do
    echo "Looking for date : $date"
    docker exec mysql sh -c "mysql -u $USERNAME -p$PASSWORD neli_webservices -e 'SELECT id, path FROM video_content WHERE leader_id = $ZOMBIE AND UNIX_TIMESTAMP(created_at) < UNIX_TIMESTAMP(\"$date 23:59:59\")' -s" | while read id path; do
            echo "Deleting video content id $id with path $path"
            docker exec mysql sh -c "mysql -u $USERNAME -p$PASSWORD neli_webservices -e 'DELETE FROM video_content WHERE id = $id' -s"
            rm $NELI_PATH/$path/$path.mp4
            rmdir $NELI_PATH/$path
            rm $NELI_PATH/public/$id.jpg
            echo "Deletion done"
    done
    docker exec mysql sh -c "mysql -u $USERNAME -p$PASSWORD neli_webservices -e 'UPDATE cleaning SET done = 1 WHERE id = $cleaning_id' -s"
    echo "Task $cleaning_id is done"
done


