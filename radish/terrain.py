import contextlib

from radish import before, after

@before.each_scenario
def init_contextlib(scenario):
    scenario.context.cleanups = contextlib.ExitStack()

@after.each_scenario
def close_contextlib(scenario):
    scenario.context.cleanups.close()

