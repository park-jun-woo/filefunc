class LogMixin:
    def log(self, msg):
        print(msg)

class AuthMixin:
    def authenticate(self, token):
        return token == "valid"
