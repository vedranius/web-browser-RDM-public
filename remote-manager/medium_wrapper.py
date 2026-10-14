"""Authenticated ASGI boundary for the privately installed original deploy tool."""
import hmac
import json
import os
import sys

# Container aliases preserve the Docker CLI's host HOME/context. Translate only
# the explicitly forwarded runtime values before loading the original app.
container_runtime = False
for key, value in list(os.environ.items()):
    if key.startswith('WRM_MEDIUM_ENV_'):
        container_runtime = True
        os.umask(0o077)
        os.environ[key[len('WRM_MEDIUM_ENV_'):]] = value
        del os.environ[key]

if container_runtime:
    import pwd
    # OpenSSH expands ~ against passwd, not HOME. Make standard ~/.ssh profiles
    # resolve to this operator's private Linux volume without mounting host keys.
    account_home = pwd.getpwuid(os.geteuid()).pw_dir
    os.makedirs(account_home, mode=0o700, exist_ok=True)
    os.chmod(account_home, 0o700)
    ssh_home = os.path.join(account_home, '.ssh')
    if not os.path.lexists(ssh_home):
        os.symlink(os.path.join(os.environ['HOME'], '.ssh'), ssh_home)

# The original application's modules and dependencies remain unchanged.
sys.path.insert(0, os.environ['WRM_MEDIUM_TOOL'])
from app.main import app as original_app


active_requests = 0


async def app(scope, receive, send):
    global active_requests
    if scope['type'] != 'http':
        return await original_app(scope, receive, send)
    supplied = dict(scope.get('headers', [])).get(b'x-wrm-medium-key', b'')
    expected = os.environ['WRM_MEDIUM_KEY'].encode()
    if not hmac.compare_digest(supplied, expected):
        await send({'type': 'http.response.start', 'status': 403,
                    'headers': [(b'content-type', b'text/plain')]})
        await send({'type': 'http.response.body', 'body': b'WRM authorization required'})
        return
    if scope['path'] == '/__wrm_ready':
        await send({'type': 'http.response.start', 'status': 200, 'headers': []})
        await send({'type': 'http.response.body', 'body': b'ready'})
        return
    if scope['path'] == '/__wrm_state':
        jobs = original_app.state.deploy_jobs
        with jobs._lock:
            running = sum(1 for job in jobs._jobs.values() if job.running)
        payload = json.dumps({'busy': active_requests > 0 or running > 0}).encode()
        await send({'type': 'http.response.start', 'status': 200,
                    'headers': [(b'content-type', b'application/json')]})
        await send({'type': 'http.response.body', 'body': payload})
        return
    active_requests += 1
    try:
        return await original_app(scope, receive, send)
    finally:
        active_requests -= 1
