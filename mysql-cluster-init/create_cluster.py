from argparse import ArgumentParser
import subprocess
from sys import platform

CONFIG_PATH = "C:/DOCKER_CLUSTER/" if platform == "win32" else "./DOCKER_CLUSTER"

def execute_command(command):
    """Execute a shell command."""
    process = subprocess.Popen(command, shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    stdout, stderr = process.communicate()

    if process.returncode == 0:
        print(f"Success: {stdout.decode('utf-8')}")
    else:
        raise RuntimeError(stderr.decode("utf-8"))

def create_network(network_name, ip):
    """Create a Docker network."""
    command = f"docker network create {network_name} --subnet={ip}0.0/16"
    execute_command(command)

def run_mysql_container(container_name, network_name, node_type, node_id, connect_string, ip_addr, port):
    """Run a MySQL NDB Cluster container."""
    command = f"docker run -d --net={network_name} --name={container_name} --ip={ip_addr} "

    mysql_cnf = CONFIG_PATH + "/mysql/my.cnf:/etc/my.cnf"
    mysql_cluster_cnf = CONFIG_PATH + "/mysql/mysql-cluster.cnf:/etc/mysql-cluster.cnf" 

    if node_type == "ndb_mgmd":
        command += f"--hostname {container_name} -v {mysql_cnf} -v {mysql_cluster_cnf} mysql/mysql-cluster ndb_mgmd --ndb-nodeid=1 --reload --initial"
    elif node_type == "data":
        command += f"-v {mysql_cluster_cnf} mysql/mysql-cluster ndbd --ndb-nodeid={node_id} --connect-string {connect_string}"
    elif node_type == "sql":
        command += f"-p {port}:3306 -v {mysql_cluster_cnf} -e MYSQL_RANDOM_ROOT_PASSWORD=true mysql/mysql-cluster mysqld --ndb-nodeid={node_id} --ndb-connectstring {connect_string}"
    else:
        print("Invalid node type specified.")
        return

    execute_command(command)

def setup_mysql_cluster(data_nodes, sql_nodes):
    network_name = "mysql-cluster"
    management_node_name = "management-1"
    ip_base = "10.100."
    management_ip = ip_base + "0.2"

    # Generate mysql-cluster.cnf
    cluster_cnf = [
        "[ndbd default]",
        "NoOfReplicas=2",
        "DataMemory=80M",
        "IndexMemory=18M",
        "[ndb_mgmd]",
        "NodeId=1",
        "hostname=" + management_ip,
        "datadir=/var/lib/mysql",
    ]

    for i in range(1, data_nodes + 1):
        cluster_cnf += [
            "[ndbd]",
            "NodeId=" + str(1 + i),
            "hostname=" + ip_base + "1." + str(i + 1),
            "datadir=/var/lib/mysql",
        ]

    for i in range(1, sql_nodes + 1):
        cluster_cnf += [
            "[mysqld]",
            "NodeId=" + str(1 + data_nodes + i),
            "hostname=" + ip_base + "2." + str(i + 1),
        ]

    with open(CONFIG_PATH + "/mysql/mysql-cluster.cnf", mode="w") as cnf:
        cnf.write("\n".join(cluster_cnf))

    # Create Docker network
    create_network(network_name, ip_base)

    # Start management node
    run_mysql_container(management_node_name, network_name, "ndb_mgmd", None, None, management_ip, None)

    # Start data nodes
    for i in range(1, data_nodes + 1):
        run_mysql_container("ndb-" + str(i), network_name, "data", i + 1, management_ip, ip_base + "1." + str(i + 1), None)

    # Start MySQL server nodes
    for i in range(1, sql_nodes + 1):
        run_mysql_container("mysqld-" + str(i), network_name, "sql", i + data_nodes + 1, management_ip, ip_base + "2." + str(i + 1), 3306 + i)

if __name__ == "__main__":
    parser = ArgumentParser()
    parser.add_argument("--data", type=int, default=2)
    parser.add_argument("--sql", type=int, default=2)

    options = parser.parse_args()
    setup_mysql_cluster(options.data, options.sql)
