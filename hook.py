import json
import ssl
import http.server


class WebhookHandler(http.server.BaseHTTPRequestHandler):
    def do_POST(self):
        data = self.rfile.read(int(self.headers['Content-Length']))
        print(data)
        try:
            admissionRequest = json.loads(data)
        except json.decoder.JSONDecodeError:
            self.send_error(400, "Expected JSON")
            return

        try:
            uid = admissionRequest['request']['uid']
        except KeyError:
            self.send_error(400, "Invalid AdmissionReview object")
            return

        admissionResponse = {
            'apiVersion': 'admission.k8s.io/v1',
            'kind': 'AdmissionReview',
            'response': {
                'uid': uid,
                'allowed': True
            }
        }
        httpResponse = json.dumps(admissionResponse)
        self.send_response(200)
        self.send_header('Content-type', 'application/json')
        self.end_headers()
        self.wfile.write(bytes(httpResponse, 'utf-8'))


def run(handler_class=http.server.BaseHTTPRequestHandler):
    certFile = "cert.pem"
    keyFile = "key.pem"
    httpd = http.server.HTTPServer(("", 8080), handler_class)
    httpd.socket = ssl.wrap_socket(
        httpd.socket, certfile=certFile, keyfile=keyFile, server_side=True)
    httpd.serve_forever()


if __name__ == '__main__':
    run(handler_class=WebhookHandler)
