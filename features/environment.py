import subprocess

import testgres

class PgTwixt(object):
    def __init__(self):
        self.location = None
        self.port = testgres.reserve_port()
        self._process = None

    def start(self):
        self._process = subprocess.Popen(
                ['go', 'run', 'cmd/pgtwixt/main.go', '127.0.0.1:{}'.format(self.port), self.location],
                stdout=subprocess.PIPE,
                stderr=subprocess.STDOUT)

    def cleanup(self):
        if self._process:
            self._process.terminate()

class Processes(dict):
    def pgtwixt(self):
        if 'pgtwixt' not in self:
            self['pgtwixt'] = PgTwixt()
        return self['pgtwixt']

    def postgres(self, name):
        if name not in self['postgres']:
            self['postgres'][name] = testgres.PostgresNode(name).init()
        return self['postgres'][name]


def before_all(context):
    context.processes = Processes({'postgres': {}})

def after_scenario(context, scenario):
    if 'pgtwixt' in context.processes:
        context.processes['pgtwixt'].cleanup()

    for node in context.processes['postgres'].values():
        node.cleanup()

