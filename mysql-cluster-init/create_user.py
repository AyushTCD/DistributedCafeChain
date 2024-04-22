import subprocess
import re

def execute_command(command, capture_output=False):
    """Execute a shell command."""
    process = subprocess.Popen(command, shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True)
    stdout, stderr = process.communicate()
    if process.returncode == 0 and capture_output:
        return stdout
    elif process.returncode == 0:
        print(f"Success: {stdout}")
    else:
        print(f"Error: {stderr}")
    return None

def get_mysql_root_password(container_name):
    """Retrieve the auto-generated MySQL root password from container logs."""
    logs = execute_command(f"docker logs {container_name}", capture_output=True)
    password_match = re.search(r"GENERATED ROOT PASSWORD: (\S+)", logs)
    if password_match:
        return password_match.group(1)
    else:
        print("Root password not found in logs.")
        return None


def create_user(container_name, root_password, new_user, new_password):
    """Create a new MySQL user using the root password."""
    # MySQL command to create a new user and grant privileges
    sql_command = (

        f"CREATE USER IF NOT EXISTS '{new_user}'@'%' IDENTIFIED BY '{new_password}'; "
        f"GRANT ALL PRIVILEGES ON *.* TO '{new_user}'@'%' WITH GRANT OPTION; FLUSH PRIVILEGES;"
    )
    command = f"docker exec {container_name} mysql -uroot -p\"{root_password}\" -e \"{sql_command}\""
    execute_command(command)


if __name__ == "__main__":
    username = "isabelle"
    user_password = "password"

    # Get ROOT password
    root_password = get_mysql_root_password("mysqld-1")

    create_user("mysqld-1", root_password, username, user_password)
