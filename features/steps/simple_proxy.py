from behave import when, then
import testgres

@when('a client connects to pgtwixt')
def connect_pgtwixt(context):
    context.processes['client'] = testgres.NodeConnection(context.processes.pgtwixt(), 'postgres')

@when('that client disconnects')
def disconnect_pgtwixt(context):
    context.processes['client'].close()

# pgtwixt_connections_total[frontend]
# pgtwixt_connections_total[backend]
@then('pgtwixt reports a frontend and a backend connect')
def log_connect(context):
    pass

# pgtwixt_disconnections_total[frontend]
# pgtwixt_disconnections_total[backend]
@then('pgtwixt reports a frontend and a backend disconnect')
def log_disconnect(context):
    pass
