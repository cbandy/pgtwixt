from prometheus_client.parser import text_string_to_metric_families
from radish import when, then
import requests
import testgres

@when('a client connects to pgtwixt')
def connect_pgtwixt(step):
    step.context.client = testgres.NodeConnection(step.context.processes['pgtwixt'], 'postgres')

@when('that client disconnects')
def disconnect_pgtwixt(step):
    step.context.client.close()

@then('pgtwixt reports a frontend and a backend connect')
def log_connect(step):
    metrics = text_string_to_metric_families(requests.get('http://{.metrics}'.format(step.context.processes['pgtwixt'])).text)
    connects = next(m for m in metrics if m.name == 'pgtwixt_connects_total').samples
    assert next(s for s in connects if "frontend" in s[1])[2] == 1
    assert next(s for s in connects if "backend" in s[1])[2] == 1

@then('pgtwixt reports a frontend and a backend disconnect')
def log_disconnect(step):
    metrics = text_string_to_metric_families(requests.get('http://{.metrics}'.format(step.context.processes['pgtwixt'])).text)
    disconnects = next(m for m in metrics if m.name == 'pgtwixt_disconnects_total').samples
    assert next(s for s in disconnects if "frontend" in s[1])[2] == 1
    assert next(s for s in disconnects if "backend" in s[1])[2] == 1
