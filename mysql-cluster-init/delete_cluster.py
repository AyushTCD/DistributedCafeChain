from argparse import ArgumentParser
import subprocess

def execute_command(command):
    """Execute a shell command."""
    process = subprocess.Popen(command, shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    stdout, stderr = process.communicate()
    if process.returncode == 0:
        print(f"Success: {stdout.decode('utf-8')}")
    else:
        raise RuntimeError(stderr.decode('utf-8'))

def delete_docker_containers(container_names):
    """Stop and remove specified Docker containers."""
    for name in container_names:
        # Stop the container
        execute_command(f"docker stop {name}")
        # Remove the container
        execute_command(f"docker rm {name}")

def delete_docker_network(network_name):
    """Remove Docker network."""
    execute_command(f"docker network rm {network_name}")

def cleanup_mysql_cluster(data_nodes, sql_nodes):
    network_name = "mysql-cluster"
    management_node_name = "management-1"
    data_node_names = ["ndb-" + str(i) for i in range(1, data_nodes + 1)]
    sql_node_names = ["mysqld-" + str(i) for i in range(1, sql_nodes + 1)]

    all_container_names = [management_node_name] + data_node_names + sql_node_names

    # Delete Docker containers
    delete_docker_containers(all_container_names)

    # Delete Docker network
    delete_docker_network(network_name)

if __name__ == "__main__":
    parser = ArgumentParser()
    parser.add_argument("--data", type=int, default=2)
    parser.add_argument("--sql", type=int, default=2)

    options = parser.parse_args()
    cleanup_mysql_cluster(options.data, options.sql)
