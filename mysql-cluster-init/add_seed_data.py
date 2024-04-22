import subprocess
from os import path

def execute_command(command, capture_output=False, stdin=None):
    """Execute a shell command."""
    process = subprocess.Popen(command, shell=True, stdin=stdin, stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True)
    stdout, stderr = process.communicate()
    if process.returncode == 0 and capture_output:
        return stdout
    elif process.returncode == 0:
        print(f"Success: {stdout}")
    else:
        print(f"Error: {stderr}")
    return None

if __name__ == "__main__":
    user = "isabelle"
    password = "password"
    host = "localhost"
    port = "3304"

    with open("./seed_data.sql") as schema:
        command = f"mysql -u {user} -h {host} -P {port} -p{password}"
        execute_command(command, stdin=schema)
