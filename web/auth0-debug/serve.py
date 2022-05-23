import os
import argparse
import functools
import http.server
import socketserver

DEFAULT_PORT=3000
BASE_PATH=os.path.join(os.path.dirname(__file__), "htdocs")

DirectoryHTTPRequestHandler = functools.partial(http.server.SimpleHTTPRequestHandler, directory=BASE_PATH)


def main(args):
    handler = DirectoryHTTPRequestHandler
    with socketserver.TCPServer(("", args.port), handler) as httpd:
        print(f"web app started at http://localhost:{args.port}")
        httpd.serve_forever()


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="auth0 debug webapp")
    parser.add_argument("-p", "--port", default=DEFAULT_PORT, type=int)
    args = parser.parse_args()
    main(args)