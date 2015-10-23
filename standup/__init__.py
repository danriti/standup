""" standup.py """

import os

from flask import Flask, render_template_string, request
import requests
import requests.packages.urllib3
requests.packages.urllib3.disable_warnings()


HIPCHAT_TOKEN = os.environ['HIPCHAT_TOKEN']
HIPCHAT_ROOM_ID = os.environ['HIPCHAT_ROOM_ID']

app = Flask(__name__)


@app.route("/", methods=['GET', 'POST'])
def index():
    success = False
    failure = False
    if request.method == 'POST':
        s = Standup(
            name=request.values.get('name', ''),
            yesterday=request.values.get('yesterday', ''),
            today=request.values.get('today', ''),
            blocked=request.values.get('blocked', 'Nope'),
            is_blocked=request.values.get('is_blocked', False, type=bool)
        )
        try:
            s.notify_hipchat()
            success = True
        except requests.exceptions.RequestException:
            app.logger.exception('Request to HipChat failed.')
            failure = True
    return render_template_string(TEMPLATE, success=success, failure=failure)


class Standup(object):

    def __init__(self, name, yesterday, today, blocked, is_blocked):
        self._name = name
        self._yesterday = yesterday
        self._today = today
        self._blocked = blocked if blocked else 'Nope'
        self._is_blocked = is_blocked
        self._room_id = HIPCHAT_ROOM_ID
        self._token = HIPCHAT_TOKEN
        self._format = 'html'

    def notify_hipchat(self):
        params = {'auth_token': self._token}
        data = {'color': self._color(),
                'message': self._message(),
                'message_format': self._format}

        r = requests.post(self._url(), params=params, data=data)
        app.logger.info('Response: %r, %r', r.status_code, r.text)
        r.raise_for_status()

    def _color(self):
        if self._is_blocked:
            return 'red'
        return 'green'

    def _url(self):
        return 'https://api.hipchat.com/v2/room/{room_id}/notification' \
            .format(room_id=self._room_id)

    def _message(self):
        return """{name}:
<ul>
<li><b>Yesterday</b>: {yesterday}</li>
<li><b>Today</b>: {today}</li>
<li><b>Blocked</b>: {blocked}</li>
</ul>
""".format(name=self._name,
           yesterday=self._yesterday,
           today=self._today,
           blocked=self._blocked)


TEMPLATE = """<html>
<head>
  <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.5/css/bootstrap.min.css">
</head>

<body>
  <div class="container">
    <h1>Astrocats Standup!</h1>

  {% if success %}
    <div class="alert alert-success" role="alert">
      <strong>Well done!</strong> Your standup was submitted to HipChat.
    </div>
  {% elif failure %}
    <div class="alert alert-danger" role="alert">
      <strong>Oh snap!</strong> Something went wrong, so please manually post in HipChat.
    </div>
  {% endif %}

    <form action="/" method="post">
      <div class="form-group">
        <label>Name:</label>
        <input type="text" name="name" class="form-control">
      </div>
      <div class="form-group">
        <label>Yesterday:</label>
        <textarea class="form-control" name="yesterday" rows="3"></textarea>
      </div>
      <div class="form-group">
        <label>Today:</label>
        <textarea class="form-control" name="today" rows="3"></textarea>
      </div>
      <div class="checkbox">
        <label>
          <input type="checkbox" name="is_blocked"> Is Blocked?
        </label>
      </div>
      <div class="form-group">
        <label>Blocked:</label>
        <textarea class="form-control" name="blocked" rows="3"></textarea>
      </div>
      <button type="submit" class="btn btn-default">Submit</button>
    </form>
  </div>
</body>
</html>
"""


if __name__ == "__main__":
    app.run(port=3000)
