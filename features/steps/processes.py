import time

from behave import given
import testgres

@given('it is running')
def start_something(context):
    if isinstance(context.subject, testgres.PostgresNode):
        return start_postgres(context, context.subject.name)

    if context.subject is context.processes.get('pgtwixt'):
        return start_pgtwixt(context)

    raise TypeError("Unknown context.subject: {!r}".format(context.subject))

@given('postgres "{}" is configured to trust connections')
def configure_postgres(context, name):
    context.subject = context.processes.postgres(name)

@given('postgres "{}" is running')
def start_postgres(context, name):
    context.subject = context.processes.postgres(name).start()

@given('pgtwixt is configured to connect to "{}"')
def configure_pgtwixt_backend(context, postgres_name):
    postgres = context.processes['postgres'][postgres_name]
    context.subject = context.processes.pgtwixt()
    context.subject.location = 'localhost:{}'.format(postgres.port)

@given('pgtwixt is running')
def start_pgtwixt(context):
    context.processes.pgtwixt().start()
    time.sleep(3) # FIXME some kind of wait
