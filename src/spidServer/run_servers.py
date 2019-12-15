import os
import subprocess


def main(mapper_address):
    base_dir = os.path.abspath('_server_instances')
    for i, _ in os.listdir(base_dir):
        subprocess.Popen([f"{base_dir}\\s{i}\\run.bat", mapper_address],
                         creationflags=subprocess.CREATE_NEW_CONSOLE)


if __name__ == '__main__':
    import argparse
    parser = argparse.ArgumentParser(description='Setup SPID server instances')
    parser.add_argument('-m', '--mapper-address', type=str, required=True, help='Server mapper address')

    args = parser.parse_args()

    main(args.mapper_address)
