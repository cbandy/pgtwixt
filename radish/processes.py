import subprocess
import time

from radish import given
import testgres

class PgTwixt(object):
    def __init__(self):
        self.host = '127.0.0.1'
        self.location = None
        self.metrics = '127.0.0.1:{}'.format(testgres.reserve_port())
        self.port = testgres.reserve_port()
        self._process = None

    def start(self):
        self._process = subprocess.Popen(
                ['radish/pgtwixt', '{0.host}:{0.port}'.format(self), self.metrics, self.location],
                stdout=subprocess.PIPE,
                stderr=subprocess.STDOUT)
        return self

    def cleanup(self):
        if self._process:
            self._process.terminate()

def pgtwixt_instance(context):
    if not hasattr(context, 'processes'): context.processes = {}

    if 'pgtwixt' not in context.processes:
        process = PgTwixt()
        context.cleanups.callback(process.cleanup)
        context.processes['pgtwixt'] = process

    return context.processes['pgtwixt']

def postgres_instance(context, name):
    if not hasattr(context, 'processes'): context.processes = {}
    context.processes.setdefault('postgres', {})

    if name not in context.processes['postgres']:
        node = testgres.PostgresNode(name).init()
        context.cleanups.push(node)
        context.processes['postgres'][name] = node

    return context.processes['postgres'][name]


@given('it is running')
def start_something(step):
    if isinstance(step.context.subject, testgres.PostgresNode):
        return start_postgres(step, step.context.subject.name)

    if step.context.subject is step.context.processes.get('pgtwixt'):
        return start_pgtwixt(step)

    raise TypeError("Unknown context.subject: {!r}".format(step.context.subject))

@given('postgres "{}" is configured to trust connections')
def configure_postgres(step, name):
    step.context.subject = postgres_instance(step.context, name)

@given('postgres "{}" is running')
def start_postgres(step, name):
    step.context.subject = postgres_instance(step.context, name).start()

@given('pgtwixt is configured to connect to "{}"')
def configure_pgtwixt_backend(step, postgres_name):
    postgres = postgres_instance(step.context, postgres_name)
    step.context.subject = pgtwixt_instance(step.context)
    step.context.subject.location = '127.0.0.1:{}'.format(postgres.port)

@given('pgtwixt is running')
def start_pgtwixt(step):
    step.context.subject = pgtwixt_instance(step.context).start()
    time.sleep(3) # FIXME some kind of wait
