***Create MySQL NDB cluster and initialise database + tables***

1. Install MySQL and Docker on your machine (add Docker and MySQL Shell to path).
2. Run `create_cluster.py --data=4 --sql=3` to create the cluster *(change the number of data and sql nodes as desired)*
3. Run `create_user.py` to access the mysql server node and create a user (update username as required)
4. Run `create_database_tables.py` to create the database and tables with schema
5. Run `delete_cluster.py --data=4 --sql=3` to stop and delete all cluster container instances *(ensure the numbers match `create_cluster.py`)*

*MySQL nodes will be available on ports 3307 and up, e.g. mysqld-1 on 3307, mysqld-2 on 3308, mysqld-3 on 3309 and so on*
