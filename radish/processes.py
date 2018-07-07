import contextlib
import subprocess
import time

from radish import after, before, given
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

    @property
    def output(self):
        if self._process:
            return self._process.stdout

def pgtwixt_instance(context):
    if 'pgtwixt' not in context.processes:
        context.processes['pgtwixt'] = PgTwixt()

    return context.processes['pgtwixt']

def postgres_instance(context, name):
    if name not in context.processes['postgres']:
        node = testgres.PostgresNode(name).init()
        context.processes['postgres'][name] = node

    return context.processes['postgres'][name]

@before.each_scenario
def processes_before(scenario):
    scenario.context.processes = {'postgres': {}}

@after.each_scenario
def processes_after(scenario):
    processes = scenario.context.processes

    if 'pgtwixt' in processes:
        process = processes['pgtwixt']
        process.cleanup()

        if scenario.state == 'failed':
            print("PgTwixt output:\n{}\n".format(process.output.read().decode('utf8')))

    with contextlib.ExitStack() as stack:
        for node in processes['postgres'].values():
            stack.push(node)

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
    step.context.subject.location = 'host=127.0.0.1 port={} sslmode=prefer'.format(postgres.port)

@given('pgtwixt is running')
def start_pgtwixt(step):
    step.context.subject = pgtwixt_instance(step.context).start()
    time.sleep(3) # FIXME some kind of wait
