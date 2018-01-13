from radish import when, then
import testgres

@when('a client connects to pgtwixt')
def connect_pgtwixt(step):
    step.context.client = testgres.NodeConnection(step.context.processes['pgtwixt'], 'postgres')

@when('that client disconnects')
def disconnect_pgtwixt(step):
    step.context.client.close()

# pgtwixt_connects_total[frontend,bind]
# pgtwixt_connects_total[backend,host]
@then('pgtwixt reports a frontend and a backend connect')
def log_connect(step):
    pass

# pgtwixt_disconnects_total[frontend,bind]
# pgtwixt_disconnects_total[backend,host]
@then('pgtwixt reports a frontend and a backend disconnect')
def log_disconnect(step):
    pass
